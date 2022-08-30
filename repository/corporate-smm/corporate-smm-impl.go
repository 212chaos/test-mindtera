package repo_corporatesmm

import (
	"context"

	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	"gorm.io/gorm/clause"
)

type CorporateSMMRepositoryImpl struct {
	sugar    logger.CustomLogger
	pgConfig gormpg.PostgresConfig
}

func NewCorporateSMMRepository(sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig) CorporateSMMRepository {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate SMM table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateSMM{})
	}

	return &CorporateSMMRepositoryImpl{
		sugar:    sugar,
		pgConfig: pgConfig,
	}
}

func (c *CorporateSMMRepositoryImpl) GetCorporateSMM(ctx context.Context, corporateSMM *[]entity.CorporateSMM) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateSMM{}).
		Select("id, code, name, english_name,sequence_number, description, icon, description_english").
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Find(corporateSMM)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to get:%v", err.Error())
		return err
	}
	return err
}

func (c *CorporateSMMRepositoryImpl) UpsertCorporateSMM(ctx context.Context, corporateSMM *entity.CorporateSMM) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateSMM{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "code"},
			},
			UpdateAll: true,
		}).
		Create(corporateSMM)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to upsert:%v", err.Error())
		return err
	}
	return err
}
