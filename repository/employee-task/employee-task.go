package repo_employeetask

import (
	"context"
	"time"

	consumermodel "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	dashentity "github.com/mindtera/dashboard-auth/entity"
	cm "github.com/mindtera/go-common-module/common/v2/model"
	fcommonmodel "github.com/mindtera/go-common-module/fajar/model"
)

type EmployeeTaskRepository interface {
	// upsert batch employee task repository
	UpsertBatchEmployeeTaskRepository(ctx context.Context, employeeTasks *[]entity.EmployeeTaskEntity) (err error)
	// get employee task based on user id, corporate subscription, and task type
	GetEmployeeTaskByUserIDSubscription(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity, pagingModel *cm.PaginationModel) (err error)
	// get employee task based on user id, corporate subscription, and task type
	GetEmployeeTaskByDepartmentAndSubsId(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity) (err error)
	// get employee task by user id and array of tsak type
	GetEmployeeTaskByUserIDAndTypes(ctx context.Context, query model.EmployeeTaskQuery, employeeTasks *[]entity.EmployeeTaskEntity) (err error)
	// get employee task based on Corporate Id and Task type
	GetTaskBySubscriptionIdAndType(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity) (err error)
	// get employee task based on user id, corporate subscription, and task type
	SearchEmployeeTaskByUserIDSubscription(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity, pagingModel *cm.PaginationModel) (err error)

	GetCorporateDetails(ctx context.Context, corporateID string) (entity.CorporateEntity, error)

	// GetEmployeeByUserID
	// get employee by user_id
	GetEmployeeByUserID(ctx context.Context, userID string) (output entity.CorporateEmployeeEntityWithDepartment, err error)

	// GetDepartment
	GetDepartment(ctx context.Context, departmentID string) (dashentity.CorporateDepartmentEntity, error)

	// GetEmployeeByCorporateID
	// get employee by corporateID
	GetEmployeeByCorporateID(ctx context.Context, corporateID string) (outputArr []entity.CorporateEmployeeEntityWithDepartment, err error)
	// GetEmployeeByCorporateIDDeptCode
	// get employee by corporateID, departmentCode
	GetEmployeeByCorporateIDDeptCode(ctx context.Context, corporateID, DepartmentCode string) (outputArr []entity.CorporateEmployeeEntityWithDepartment, err error)

	InsertEmployeeTaskProgramAssignment(ctx context.Context, tasks []entity.EmployeeTaskEntity) (err error)

	GetProgramAssignmentHistory(ctx context.Context, corporateID string, limit, offset int) (historyArr []entity.EmployeeProgramAssignmentHistory2, total int, err error)

	InsertEmployeeProgramAssignmentHistory(ctx context.Context, assigningMode entity.TypeAssigningMode, programIDArr []string, corporateID, relatedID, createdBy string) (err error)

	GetEmployeeTaskProgramAssignment(ctx context.Context, taskProgramAssignmentArr *cm.PaginationResponseModel, corporateID string) (err error)

	GetEmployeeTaskByUserIDProgramIDTaskType(ctx context.Context, userIDArr []string, programIDArr []string, taskTypeArr []entity.EMPLOYEE_TASK_TYPE, taskStatusArr []entity.EMPLOYEE_TASK_STATUS) ([]entity.EmployeeTaskEntityForConsumer, error)

	GetEnrollmentByUserIDArr(ctx context.Context, userIDArr []string) ([]fcommonmodel.ProgramEnrollmentHistory, error)

	SearchPrograms(ctx context.Context, keyword string) ([]consumermodel.MProgram, error)

	GetPrograms(ctx context.Context) ([]consumermodel.MProgram, error)

	GetCorporateServices(ctx context.Context, params model.MCorporateServiceRequest) (outputArr []entity.MCorporateService, totalData int, err error)

	GetCorporateServiceByID(ctx context.Context, id string) (entity.MCorporateService, error)

	InsertCorporateServiceBooking(ctx context.Context, serviceID, corporateID string, scheduledAt time.Time) (err error)

	GetCorporateServiceBooking(ctx context.Context, corporateID string, limit, offset int) ([]entity.CorporateServiceBooking, int, error)

	SearchCorporateService(ctx context.Context, keyword string) ([]entity.MCorporateService, error)

	// insert new well being caused by scheduler trigger
	RenewAllWellBeingTask(ctx context.Context, start, end time.Time) (err error)
}
