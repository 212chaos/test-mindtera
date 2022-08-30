package repository_corporateprogram

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"

	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateProgramImpl struct {
	sugar     logger.CustomLogger
	pgConfig  gormpg.PostgresConfig
	commonSvc commonsvc.CommonService
}

//construct
func NewCorporateProgram(sugar logger.CustomLogger, pgConfig gormpg.PostgresConfig, commonSvc commonsvc.CommonService) CorporateProgram {
	return &CorporateProgramImpl{
		sugar:     sugar,
		pgConfig:  pgConfig,
		commonSvc: commonSvc}
}

// get calculated program employee
func (c *CorporateProgramImpl) GetCorporateEmployeeProgram(ctx context.Context, employee entity.EmployeeTaskEntity) (result []model.EmployeeProgram, err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	// get skipped program
	var skippedProgram []model.EmployeeProgram
	var skip []map[string]any
	tx.Raw(c.getSkippedProgramAmount(),
		employee.UserID,
		employee.SubscriptionID).
		Scan(&skip)
	c.commonSvc.ObjectMapper(&skip, &skippedProgram)

	skippedChan := make(chan map[string]int32)
	go func(val []model.EmployeeProgram) {
		hash := make(map[string]int32)
		for _, v := range val {
			hash[v.ProgramId.String()] = v.Amount
		}
		skippedChan <- hash
	}(skippedProgram)

	tx.Raw(c.getAvailableProgramAmount(),
		employee.CorporateId,
		employee.UserID,
		employee.SubscriptionID,
		employee.UserID).
		Scan(&result)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to get:%v", err.Error())
		return result, err
	}

	if skippedHash := <-skippedChan; skippedHash != nil {
		for i := range result {
			result[i].Amount -= skippedHash[result[i].ProgramId.String()]
			delete(skippedHash, result[i].ProgramId.String())
		}
		for key, val := range skippedHash {
			result = append(result, model.EmployeeProgram{
				ProgramId: uuid.MustParse(key),
				Amount:    -val,
			})
		}
	}
	return result, err
}
