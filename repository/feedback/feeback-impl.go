package repo_feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"
	"gorm.io/gorm/clause"

	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonservice "github.com/mindtera/go-common-module/common/v2/service/common"
)

var (
	ValidShownCategory = map[string]string{
		"DEPARTMENT":          "id, department, record_flag",
		"NAME_AND_DEPARTMENT": "id, department, employee_name, record_flag"}
)

type FeedbackRepositoryImpl struct {
	sugar     logger.CustomLogger
	pgConfig  gormpg.PostgresConfig
	assert    assert.Assert
	commonSvc commonservice.CommonService
}

func NewFeedbackRepository(sugar logger.CustomLogger, pgConfig gormpg.PostgresConfig,
	assert assert.Assert, commonSvc commonservice.CommonService) FeedbackRepository {
	return &FeedbackRepositoryImpl{
		sugar:     sugar,
		assert:    assert,
		pgConfig:  pgConfig,
		commonSvc: commonSvc,
	}
}

func (f *FeedbackRepositoryImpl) GetFeedbackByFilter(ctx context.Context, feedbackFilter model.FeedbackFilter,
	pagingMeta *mdl.PaginationModel, feedback *[]entity.Feedback) (err error) {
	filter := model.FeedbackQueryFilter{
		ShownCategoryCode: feedbackFilter.ShownCategory,
		CategoryCode:      feedbackFilter.Category,
		CorporateId:       feedbackFilter.CorporateId,
	}

	// performing concurrency
	ctxBackground := f.commonSvc.ContextBackground(ctx)

	totalChan := make(chan int64)
	dataChan := make(chan []entity.Feedback)
	resultErr := make(chan error)

	go f.feedTotalPage(ctxBackground, totalChan, resultErr, filter, feedbackFilter)
	go f.feedTotalData(ctxBackground, dataChan, resultErr, filter, *pagingMeta, feedbackFilter)

	for i := 0; i < 2; i++ {
		if err = <-resultErr; err != nil {
			f.sugar.WithContext(ctx).Errorf("error in channel:%v", err)
			return err
		}
	}
	// getting all data from channel
	f.sugar.WithContext(ctx).Infof("getting value from in channel:%v", err)
	data := <-dataChan
	(*pagingMeta).TotalData = int(<-totalChan)
	(*pagingMeta).DataPerPage = len(data)
	(*feedback) = data

	return err
}

func (f *FeedbackRepositoryImpl) UpsertFeedback(ctx context.Context, feedback *entity.Feedback) (err error) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)
	// query feedback
	f.sugar.WithContext(ctx).Infof("upserting feedback to database with")
	tx.Model(&entity.Feedback{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			UpdateAll: true,
		}).
		Create(feedback)

	// checking error
	if err = tx.Error; err != nil {
		f.sugar.WithContext(ctx).Errorf("error when creating data in database:%v", err)
	}
	return err
}

func (f *FeedbackRepositoryImpl) MarkFeedbackAsRead(ctx context.Context, feedbackId, companyId uuid.UUID) (err error) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)

	// query feedback
	f.sugar.WithContext(ctx).Infof("marked feedback to database with id:%v corp id:%v", feedbackId, companyId)
	tx.Model(&entity.Feedback{}).
		Where("id = ? AND corporate_id = ?", feedbackId, companyId).
		Updates(entity.Feedback{ReadFlag: true})

	// checking error
	if err = tx.Error; err != nil {
		f.sugar.WithContext(ctx).Errorf("error when marked data in database:%v", err)
	}
	return err
}

func (f *FeedbackRepositoryImpl) GetFeedbackCategory(ctx context.Context, feedbackCategory *[]entity.FeedbackCategory) (err error) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)

	// query feedback
	f.sugar.WithContext(ctx).Infof("fetching feedback category from database ")
	tx.Model(&entity.FeedbackCategory{}).
		Select(`
		id, 
		code, 
		corporate_feedback_categories.name as localization_name_Indonesia,
		corporate_feedback_categories.english_name as localization_name_English`).
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Order("created_at ASC").
		Find(feedbackCategory)

	// checking error
	if err = tx.Error; err != nil {
		f.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)
	}
	return err
}

func (f *FeedbackRepositoryImpl) GetFeedbackShownCategory(ctx context.Context, feedbackShownCategory *[]entity.FeedbackShownCategory) (err error) {
	// generate transaction
	tx := f.pgConfig.GenerateTransaction(ctx)

	// query feedback
	f.sugar.WithContext(ctx).Infof("fetching feedback shown category from database")
	tx.Model(&entity.FeedbackShownCategory{}).
		Select(`
		id, 
		code, 
		corporate_feedback_shown_categories.description as localization_description_Indonesia,
		corporate_feedback_shown_categories.english_description as localization_description_English,
		corporate_feedback_shown_categories.name as localization_name_Indonesia,
		corporate_feedback_shown_categories.english_name as localization_name_English`).
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Order("created_at ASC").
		Find(feedbackShownCategory)

	// checking error
	if err = tx.Error; err != nil {
		f.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)
	}
	return err
}
