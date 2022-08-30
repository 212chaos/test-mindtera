package service_corporateemployee

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	cm "github.com/mindtera/go-common-module/common/v2/model"
	qm "github.com/mindtera/quiz-assessment-service/model"
)

type CorporateEmployeeService interface {
	GetCorporateEmployeeService(ctx context.Context, corporatePaging *cm.PaginationResponseModel) (err error)
	GetCorporateEmployeeStatusService(ctx context.Context, corporateID uuid.UUID, corporateStatus *model.CorporateStatus) (err error)
	GetCorporateEmployeeByNameService(ctx context.Context, employeeName string, corporatePaging *cm.PaginationResponseModel) (err error)
	GetCorporateEmployeeByPublicIDService(ctx context.Context, corporateEmployeeInfo *model.CorporateInformation) (err error)
	// getting corporate employee entity from email
	GetCorporateEmployeeByEmailService(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error)
	UpsertCorporateEmployeeService(ctx context.Context, corporateEmployees *[]entity.CorporateEmployeeEntity) (errorRecord []model.CorporateEmployeeError, err error)
	DeleteCorporateEmployeeService(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error)
	// getting dashboard status
	GetCorporateEmployeeDashboardService(ctx context.Context, query *model.CorpQuery, assessment *qm.CorporateAssessmentResponseModel) (errMsg cm.ErrorMessage)
	// getting sum of employee already filled the assessment
	GetCorpEmployeeAssessmentStatus(ctx context.Context, assessStatus *model.AssessmentStatus) (err error)
	GetWellBeingAssessmentStatus(ctx context.Context, wellBeingStatus *model.WellBeingStatus) (errMsg cm.ErrorMessage)

	// get active employee by corp and subs id
	UpdateCorpEmployeeByCorpAndSubsId(ctx context.Context, subsEntity entity.CorporateSubscriptionEntity, status string) (err error)
	GenerateEmployeeChanByCorpAndSubsId(ctx context.Context, employees []entity.CorporateEmployeeEntity) <-chan entity.CorporateEmployeeEntity
}
