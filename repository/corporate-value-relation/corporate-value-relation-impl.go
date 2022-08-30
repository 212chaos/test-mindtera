package corporatevaluerelationrepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CorporateValueRelationRepositoryImpl struct {
	sugar    logger.CustomLogger
	pgConfig gormpg.PostgresConfig
}

func NewCorporateValueRelationRepo(
	sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig) CorporateValueRelationRepo {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate value relation table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateValueRelation{})
	}

	return &CorporateValueRelationRepositoryImpl{
		sugar:    sugar,
		pgConfig: pgConfig}
}

func (c *CorporateValueRelationRepositoryImpl) GetCorporateValueRelationByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValRel *[]entity.CorporateValueRelation) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueRelation{}).
		Select("id, corporate_id, corporate_value_id, record_flag").
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Where("corporate_id=?", corporateID).
		Preload("CorporateValueDetail", common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, english_name, code, record_flag, program_tagging_id")
			}).
		Find(corpValRel)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) GetCorporateValueRelationWithAssessmentByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValRel *[]model.CorporateRelationAssessmentModel) (err error) {
	c.sugar.WithContext(ctx).Info(corporateID)
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateValueRelation{}).
		Select(`
			corporate_value_relations.id, 
			corporate_value_relations.corporate_id, 
			corporate_value_relations.corporate_value_id, 
			corporate_value_relations.record_flag,
			corporate_assessment_relations.parent_smm,
			corporate_value_entities.name as name, 
			corporate_value_entities.english_name as english_name,
			corporate_value_entities.code as code,
			corporate_assessment_relations.assessment_id`).
		Joins(`left join corporate_assessment_relations on 
			corporate_assessment_relations.corporate_value_id = corporate_value_relations.corporate_value_id`).
		Joins(`left join corporate_value_entities on 
			corporate_value_entities.id = corporate_value_relations.corporate_value_id`).
		Where("corporate_value_relations.corporate_id=?", corporateID).
		Where(`corporate_value_relations.deleted_at IS NULL AND 
			corporate_value_relations.record_flag = 'ACTIVE'`).
		Where(`corporate_value_entities.deleted_at IS NULL AND 
			corporate_value_entities.record_flag = 'ACTIVE' AND corporate_value_entities.hide = 'false'`).
		Where(`corporate_assessment_relations.deleted_at IS NULL AND 
			corporate_assessment_relations.record_flag = 'ACTIVE'`).
		Find(corpValRel)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) GetCorpValRelWithAssessByCorporateID(ctx context.Context, corpId uuid.UUID, corpValRel *[]model.CorpValueId) (err error) {
	c.sugar.WithContext(ctx).Infof("getting relation for corp:%v", corpId)
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateValueRelation{}).
		Select(`corporate_value_relations.corporate_value_id as "corp_value_id", 
			corporate_assessment_relations.parent_smm as "corp_smm",
			corporate_assessment_relations.assessment_id as "task_assess_id"`).
		Joins(`left join corporate_assessment_relations on 
			corporate_assessment_relations.corporate_value_id = corporate_value_relations.corporate_value_id`).
		Joins(`left join corporate_value_entities on 
			corporate_value_entities.id = corporate_value_relations.corporate_value_id`).
		Where("corporate_value_relations.corporate_id=?", corpId).
		Where(`corporate_value_relations.deleted_at IS NULL AND 
			corporate_value_relations.record_flag = 'ACTIVE'`).
		Where(`corporate_value_entities.deleted_at IS NULL AND 
			corporate_value_entities.record_flag = 'ACTIVE' AND corporate_value_entities.hide = 'false'`).
		Where(`corporate_assessment_relations.deleted_at IS NULL AND 
			corporate_assessment_relations.record_flag = 'ACTIVE'`).
		Find(corpValRel)
	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) GetCorpValRelNameByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValName *[]model.CorpValueName) (err error) {
	c.sugar.WithContext(ctx).Infof("getting relation name for corp:%v", corporateID)
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateValueRelation{}).
		Select(`
		corporate_value_relations.corporate_value_id, 
		corporate_value_entities.name as localization_name_Indonesia, 
		corporate_value_entities.english_name as localization_name_English`).
		Joins(`left join corporate_value_entities on 
		corporate_value_entities.id = corporate_value_relations.corporate_value_id`).
		Where("corporate_value_relations.corporate_id=?", corporateID).
		Where(`corporate_value_relations.deleted_at IS NULL AND 
		corporate_value_relations.record_flag = 'ACTIVE'`).
		Where(`corporate_value_entities.deleted_at IS NULL AND 
		corporate_value_entities.record_flag = 'ACTIVE' AND corporate_value_entities.hide = 'false'`).
		Find(corpValName)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) GetCorpValRelDetWithAssessByCorporateID(ctx context.Context, corpId uuid.UUID, corpValRel *[]model.CorpValue) (err error) {
	c.sugar.WithContext(ctx).Infof("getting relation for corp:%v", corpId)
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateValueRelation{}).
		Select(`corporate_value_relations.corporate_value_id as "corp_value_id", 
			corporate_assessment_relations.parent_smm as "corp_smm",
			corporate_assessment_relations.assessment_id as "task_assess_id",
			corporate_value_entities.*,
			corporate_value_entities.sequence_number as "detail_sequence_number",
			corporate_value_entities.code as "detail_code",
			corporate_value_entities.id as "detail_id",
			corporate_value_entities.name as "detail_name",
			corporate_value_entities.english_name as "detail_english_name"`).
		Joins(`left join corporate_assessment_relations on 
			corporate_assessment_relations.corporate_value_id = corporate_value_relations.corporate_value_id`).
		Joins(`left join corporate_value_entities on 
			corporate_value_entities.id = corporate_value_relations.corporate_value_id`).
		Where("corporate_value_relations.corporate_id=?", corpId).
		Where(`corporate_value_relations.deleted_at IS NULL AND 
			corporate_value_relations.record_flag = 'ACTIVE'`).
		Where(`corporate_value_entities.deleted_at IS NULL AND 
			corporate_value_entities.record_flag = 'ACTIVE' AND corporate_value_entities.hide = 'false'`).
		Where(`corporate_assessment_relations.deleted_at IS NULL AND 
			corporate_assessment_relations.record_flag = 'ACTIVE'`).
		Find(corpValRel)
	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) UpsertCorporateValueRelationByCorporateID(ctx context.Context, corpValRel *[]entity.CorporateValueRelation) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueRelation{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "corporate_id"},
				{Name: "corporate_value_id"}},
			UpdateAll: true}).
		Create(corpValRel)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to upsert:%v", err.Error())
	}
	return err
}

func (c *CorporateValueRelationRepositoryImpl) DeleteCorporateValueRelationByCorporateID(ctx context.Context, corporateID uuid.UUID) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueRelation{}).
		Unscoped().
		Where("corporate_id = ?", corporateID).
		Delete(entity.CorporateValueRelation{})

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to upsert:%v", err.Error())
	}
	return err
}
