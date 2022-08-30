package repo_corporateemployee

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"gorm.io/gorm"

	publicuser "github.com/mindtera/consumer-service/model"
	cm "github.com/mindtera/go-common-module/common/v2/model"
)

type CorporateEmployeeRepository interface {
	GetCorporateEmployeeRepository(ctx context.Context, corporateID uuid.UUID, corporateEmployees *[]entity.CorporateEmployeeEntity, pagingModel *cm.PaginationModel) (err error)
	GetCorporateEmployeeByID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error)
	GetCorporateEmployeeByPublicID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error)
	GetCorporateEmployeeByName(ctx context.Context, filter map[string]string, corporateEmployees *[]entity.CorporateEmployeeEntity, pagingModel *cm.PaginationModel) (err error)
	GetCorporateEmployeeIDRepository(ctx context.Context, corporateID uuid.UUID, corporateEmployees *[]entity.CorporateEmployeeEntity, pagingModel *cm.PaginationModel) (err error)
	UpsertCorporateEmployeeRepository(ctx context.Context, corporateEmployees *[]entity.CorporateEmployeeEntity) (err error)
	GetEmployeeSum(ctx context.Context, corporateID uuid.UUID, count *int64) (err error)
	GetEmployeeDeleted(ctx context.Context, corporateID uuid.UUID, count *int64) (err error)
	DeleteEmployee(ctx context.Context, employeeEntity *entity.CorporateEmployeeEntity) (err error)
	// with transaction
	GetEmployeeSumWithTx(ctx context.Context, tx *gorm.DB, corporateID uuid.UUID, count *int64) (err error)
	// get employee public user data
	GetUserByUserEmail(ctx context.Context, emails []string, publicUser *[]publicuser.User) (err error)
	// get corporate employee by email
	GetCorporateEmployeeWithQuery(ctx context.Context, corporateEmployee *entity.CorporateEmployeeEntity) (err error)
	// get employee by public id and subscription and validate whether subscription is valid or not
	GetCorporateEmplByPublicAndSubsID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error)
	// get all active public id by department
	GetUsersByDepartment(ctx context.Context, corporate_id uuid.UUID, department string, publicId *[]string) (err error)
	// get all active public id by email
	GetUsersByEmail(ctx context.Context, email []string, employees *[]entity.CorporateEmployeeEntity) (err error)
	// get amount filled employee
	GetEmployeeAssessStatus(ctx context.Context, subsId uuid.UUID, assessStatus *model.AssessmentStatus) (err error)
	// get amount filled employee
	GetEmployeeWellBeingAssessStatus(ctx context.Context, subsId uuid.UUID, assessStatus *model.WellBeingStatus) (err error)
	GetEmployeeByCorpAndSubsId(ctx context.Context, corpId, subsId uuid.UUID, corpEmpls *[]entity.CorporateEmployeeEntity) (err error)
}
