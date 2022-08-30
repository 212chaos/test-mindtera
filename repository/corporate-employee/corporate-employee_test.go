package repo_corporateemployee_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	corpemplrepo "github.com/mindtera/corporate-service/repository/corporate-employee"
	usersubsrepo "github.com/mindtera/corporate-service/repository/user-subscription"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	cusassert "github.com/mindtera/go-common-module/common/v2/service/assert"
)

var (
	sugar          = logger.NewCustomLogger()
	sugarNoContext = logger.NewWrappedZapLogger()
	pgConf         = gormpg.NewPostgresConfig(sugarNoContext)
	userSubsRepo   = usersubsrepo.NewUserSubscriptionRepository(sugar, pgConf)
	corpEmplRepo   = corpemplrepo.NewCorporateEmployeeRepository(sugar, pgConf, userSubsRepo)
	cusAssert      = cusassert.NewAssert()
)

type TestCase struct {
	Entity entity.CorporateEmployeeEntity
	Name   string
	Type   bool
}

func TestGetCorporateEmplByPublicAndSubsID(t *testing.T) {
	//TODO: change use case mock data
	testCase := []TestCase{
		{
			Entity: entity.CorporateEmployeeEntity{
				CorporateSubscriptionID: uuid.MustParse("522e00f5-5049-4090-86ea-f7ca8f48a40d"),
				PublicUserID:            uuid.MustParse("eeac3bb5-d856-4a98-a139-220af8cf713c"),
			},
			Name: "Positive Test Case",
			Type: true,
		},
		{
			Entity: entity.CorporateEmployeeEntity{
				CorporateSubscriptionID: uuid.New(),
				PublicUserID:            uuid.New(),
			},
			Name: "Negative Test Case",
			Type: false,
		},
	}

	for _, v := range testCase {
		t.Run(v.Name, func(tr *testing.T) {
			err := corpEmplRepo.GetCorporateEmplByPublicAndSubsID(context.Background(), &v.Entity)
			if err != nil {
				tr.Errorf("ERROR for payload %v", v.Entity)
			}
			if !(cusAssert.IsUUIDEmpty(v.Entity.ID.String())) == v.Type {
				tr.Errorf("ERROR: got:%v", v.Entity.Email)
			}
			if !(cusAssert.IsEmpty(v.Entity.Email)) == v.Type {
				tr.Errorf("ERROR: got:%v", v.Entity.Email)
			}
		})
	}

}
