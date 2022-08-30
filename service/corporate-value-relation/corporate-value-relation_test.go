package service_corporatevaluerelation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/model"
	repo_corporateassessmentrelation "github.com/mindtera/corporate-service/repository/corporate-assessment-relation"
	corpemplrepo "github.com/mindtera/corporate-service/repository/corporate-employee"
	repo_corporatesmm "github.com/mindtera/corporate-service/repository/corporate-smm"
	corporatevaluerepository "github.com/mindtera/corporate-service/repository/corporate-value"
	corporatevaluerelationrepository "github.com/mindtera/corporate-service/repository/corporate-value-relation"
	usersubsrepo "github.com/mindtera/corporate-service/repository/user-subscription"
	service_corporatevalue "github.com/mindtera/corporate-service/service/corporate-value"
	corpvalrel "github.com/mindtera/corporate-service/service/corporate-value-relation"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	config_googlecloud "github.com/mindtera/go-common-module/common/v2/configuration/google-cloud"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	redisconfig "github.com/mindtera/go-common-module/common/v2/configuration/redis"
	cusassert "github.com/mindtera/go-common-module/common/v2/service/assert"
	service_common "github.com/mindtera/go-common-module/common/v2/service/common"
	service_redis "github.com/mindtera/go-common-module/common/v2/service/redis"
)

var (
	sugar               = logger.NewCustomLogger()
	sugarWithoutContext = logger.NewWrappedZapLogger()
	pgConf              = gormpg.NewPostgresConfig(sugarWithoutContext)
	userSubsRepo        = usersubsrepo.NewUserSubscriptionRepository(sugar, pgConf)
	corpEmplRepo        = corpemplrepo.NewCorporateEmployeeRepository(sugar, pgConf, userSubsRepo)
	corpValRelRepo      = corporatevaluerelationrepository.NewCorporateValueRelationRepo(sugar, pgConf)
	cusAssert           = cusassert.NewAssert()
	common              = service_common.NewCommonService()
	redisConf           = redisconfig.NewRedisConfig()
	redisSvc            = service_redis.NewRedisSvc(sugarWithoutContext, redisConf)
	gcloud              = config_googlecloud.NewGoogleStorageConfig()
	corpValRepo         = corporatevaluerepository.NewCorporateValueRepository(sugar, pgConf)
	corpValSvc          = service_corporatevalue.NewCorporateValueService(sugar, gcloud, corpValRepo,
		repo_corporateassessmentrelation.NewCorporateAssessmentRelationRepo(sugar, pgConf),
		repo_corporatesmm.NewCorporateSMMRepository(sugar, pgConf),
		redisSvc, common)
	corpValRelSvc = corpvalrel.NewCorporateValueRelationService(sugar,
		&config_googlecloud.GoogleStorageConfigImpl{},
		corpValRelRepo,
		&repo_corporateassessmentrelation.CorporateAssessmentRelationRepositoryImpl{},
		corpEmplRepo,
		redisSvc,
		cusAssert,
		common,
		corpValSvc)
)

type TestCase struct {
	Type     bool
	Name     string
	CorpInfo *pb.CorpInfo
}

func TestGetCorporateValueIds(t *testing.T) {
	//TODO: change use case mock data
	testCase := []TestCase{
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.MustParse("382d1995-db6c-43e2-afe6-f4a929816982").String(),
				SubsId:     uuid.MustParse("522e00f5-5049-4090-86ea-f7ca8f48a40d").String(),
				EmployeeId: uuid.MustParse("eeac3bb5-d856-4a98-a139-220af8cf713c").String(),
			},
			Name: "Positive Test Case",
			Type: true,
		},
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.NewString(),
				SubsId:     uuid.NewString(),
				EmployeeId: uuid.NewString(),
			},
			Name: "Negative Test Case",
			Type: false,
		},
	}

	for _, v := range testCase {
		t.Run(v.Name, func(tr *testing.T) {
			var val []model.CorpValueId
			err := corpValRelSvc.GetCorporateValueIds(context.Background(), v.CorpInfo, &val)
			if !v.Type {
				if err == nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			} else {
				if err != nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			}
		})
	}
}
func TestGetCorporateValueDetail(t *testing.T) {
	//TODO: change use case mock data
	testCase := []TestCase{
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.MustParse("382d1995-db6c-43e2-afe6-f4a929816982").String(),
				SubsId:     uuid.MustParse("522e00f5-5049-4090-86ea-f7ca8f48a40d").String(),
				EmployeeId: uuid.MustParse("eeac3bb5-d856-4a98-a139-220af8cf713c").String(),
			},
			Name: "Positive Test Case",
			Type: true,
		},
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.NewString(),
				SubsId:     uuid.NewString(),
				EmployeeId: uuid.NewString(),
			},
			Name: "Negative Test Case",
			Type: false,
		},
	}

	for _, v := range testCase {
		t.Run(v.Name, func(tr *testing.T) {
			var val []model.CorpValue
			err := corpValRelSvc.GetCorporateValueDetail(context.Background(), v.CorpInfo, &val)
			if !v.Type {
				if err == nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			} else {
				if err != nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			}
		})
	}
}

func TestGetCorporateValueName(t *testing.T) {
	//TODO: change use case mock data
	testCase := []TestCase{
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.MustParse("382d1995-db6c-43e2-afe6-f4a929816982").String(),
				SubsId:     uuid.MustParse("522e00f5-5049-4090-86ea-f7ca8f48a40d").String(),
				EmployeeId: uuid.MustParse("eeac3bb5-d856-4a98-a139-220af8cf713c").String(),
			},
			Name: "Positive Test Case",
			Type: true,
		},
		{
			CorpInfo: &pb.CorpInfo{
				Id:         uuid.NewString(),
				SubsId:     uuid.NewString(),
				EmployeeId: uuid.NewString(),
			},
			Name: "Negative Test Case",
			Type: false,
		},
	}

	for _, v := range testCase {
		t.Run(v.Name, func(tr *testing.T) {
			var val []model.CorpValueName
			err := corpValRelSvc.GetCorporateValueName(context.Background(), v.CorpInfo, &val)
			if !v.Type {
				if err == nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			} else {
				if err != nil {
					tr.Errorf("ERROR for payload %v", v.CorpInfo)
				}
			}
		})
	}
}
