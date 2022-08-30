package corporatevaluerepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
)

type CorporateValueRepoImpl struct {
	sugar    logger.CustomLogger
	pgConfig gormpg.PostgresConfig
}

func NewCorporateValueRepository(sugar logger.CustomLogger, pgConfig gormpg.PostgresConfig) CorporateValueRepository {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate value table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateValueEntity{})
	}

	return &CorporateValueRepoImpl{
		sugar:    sugar,
		pgConfig: pgConfig,
	}
}

func (crv *CorporateValueRepoImpl) GetCorporateValue(ctx context.Context, corpVals *[]entity.CorporateValueEntity) (err error) {
	tx := crv.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueEntity{}).
		Select("id, name, english_name, code, record_flag, program_tagging_id").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE).
		Order("sequence_number asc").
		Find(corpVals)

	if err = tx.Error; err != nil {
		tx.Rollback()
		crv.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}

	return err
}

func (crv *CorporateValueRepoImpl) GetCorporateValueByID(ctx context.Context, cids []uuid.UUID, corpVals *[]entity.CorporateValueEntity) (err error) {
	tx := crv.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueEntity{}).
		Select("id, name, english_name, code, record_flag, program_tagging_id").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE).
		Where("id IN ?", cids).
		Order("sequence_number asc").
		Find(corpVals)

	if err = tx.Error; err != nil {
		tx.Rollback()
		crv.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}

	return err
}

func (crv *CorporateValueRepoImpl) GetCorporateValueByCode(ctx context.Context, corpVal *entity.CorporateValueEntity) (err error) {
	tx := crv.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateValueEntity{}).
		Select("id, name, english_name, code, record_flag, program_tagging_id").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_WITH_HIDE).
		Where("code", corpVal.Code).
		Find(corpVal)

	if err = tx.Error; err != nil {
		tx.Rollback()
		crv.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}

	return err
}
