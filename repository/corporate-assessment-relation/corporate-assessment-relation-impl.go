package repo_corporateassessmentrelation

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	"gorm.io/gorm"
)

type CorporateAssessmentRelationRepositoryImpl struct {
	sugar    logger.CustomLogger
	pgConfig gormpg.PostgresConfig
}

func NewCorporateAssessmentRelationRepo(sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig) CorporateAssessmentRelationRepository {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate assessment relation table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateAssessmentRelation{})
	}
	return &CorporateAssessmentRelationRepositoryImpl{
		sugar:    sugar,
		pgConfig: pgConfig}
}

func (c *CorporateAssessmentRelationRepositoryImpl) GetCorporateAssessmentRelation(ctx context.Context, corporateAssessment *[]entity.CorporateAssessmentRelation) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	tx.Model(&entity.CorporateAssessmentRelation{}).
		Select("id, corporate_value_id, assessment_id, item_code, parent_smm").
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Preload("CorporateValueDetail", common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, english_name, code, record_flag, program_tagging_id")
			}).
		Find(corporateAssessment)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}

	return err
}

func (c *CorporateAssessmentRelationRepositoryImpl) GetCorporateAssessmentRelationByValueID(ctx context.Context, corporateValueID []uuid.UUID, corporateAssessment *[]entity.CorporateAssessmentRelation) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateAssessmentRelation{}).
		Select("id, corporate_value_id, assessment_id, item_code, parent_smm, record_flag").
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Where("corporate_value_id IN ?", corporateValueID).
		Preload("CorporateValueDetail", common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, english_name, code, record_flag, program_tagging_id")
			}).
		Find(corporateAssessment)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}

	return err
}
