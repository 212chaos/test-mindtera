package repo_feedback

import (
	"context"

	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
	"gorm.io/gorm"
)

func (f *FeedbackRepositoryImpl) feedTotalPage(ctx context.Context,
	totalChan chan int64, err chan error, categoryFilter model.FeedbackQueryFilter, feedbackFilter model.FeedbackFilter) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)
	f.sugar.WithContext(ctx).Info("performing feedback counting for :%v", feedbackFilter)
	// get total amount
	var total int64
	tx.Model(&entity.Feedback{}).
		Select("count(1)").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("(created_at >= ? AND created_at <= ?)",
			feedbackFilter.StartDate,
			feedbackFilter.EndDate).
		Where(categoryFilter).
		Order("created_at DESC").
		Count(&total)
	// checking error
	if tx.Error != nil {
		f.sugar.WithContext(ctx).Errorf("error when fetching total data from database:%v", err)
	}
	// updating channel
	err <- tx.Error
	totalChan <- total
	f.sugar.WithContext(ctx).Info("got total data:%v", total)
}

func (f *FeedbackRepositoryImpl) feedTotalData(ctx context.Context,
	dataChan chan []entity.Feedback, err chan error,
	categoryFilter model.FeedbackQueryFilter,
	pagingMeta mdl.PaginationModel, feedbackFilter model.FeedbackFilter) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)
	f.sugar.WithContext(ctx).Infof("fetching feedback from database with filter :%v", feedbackFilter)
	// query feedback
	var feedback []entity.Feedback
	tx.Model(&entity.Feedback{}).
		Select(`
				id, 
				user_id, 
				corporate_id,
				body, 
				sender, 
				created_at, 
				category_code,
				shown_category_code,
				read_flag`).
		Scopes(f.pgConfig.PaginateQuery(&pagingMeta)).
		Preload("EmployeeDetail", common.ACTIVE_RECORD_FLAG_QUERY,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("public_user_id, department, corporate_id, employee_name, record_flag").
					Preload("DepartmentDetail", common.ACTIVE_RECORD_FLAG_QUERY,
						func(db *gorm.DB) *gorm.DB {
							return db.Select("code, name, corporate_id").
								Where("corporate_id = ?", categoryFilter.CorporateId)
						})
			}).
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Where("(created_at >= ? AND created_at <= ?)",
			feedbackFilter.StartDate,
			feedbackFilter.EndDate).
		Where(categoryFilter).
		Order("created_at DESC").
		Find(&feedback)
	// checking error
	if tx.Error != nil {
		f.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)
	}
	// updating channel
	err <- tx.Error
	dataChan <- feedback
}
