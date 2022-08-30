package service_employeetask

import (
	"context"
	"time"

	consumermodel "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
)

type EmployeeTaskService interface {
	// get employee task based on query needed
	GetEmployeeTaskByUserIDSubscriptionService(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTaskPagination *commonmodel.PaginationResponseModel) (err error)
	// Search employee task based on wild card query needed
	SearchEmployeeTaskByUserIDSubscriptionService(ctx context.Context, search string, employeeTask *entity.EmployeeTaskEntity, employeeTaskPagination *commonmodel.PaginationResponseModel) (err error)
	// Get employee task by user id and task type
	GetEmployeeTaskByUserIdAndType(ctx context.Context, empTask *model.EmployeeTaskQuery, empTasks *[]entity.EmployeeTaskEntity) (err error)
	// upsert employee in batch
	UpsertEmployeeInBatchService(ctx context.Context, employeeTasks *[]entity.EmployeeTaskEntity) (err error)
	// upsert employee from corporate employee entity
	UpsertEmployeeTaskByEmployeeEntityService(ctx context.Context, employeeEntities *[]entity.CorporateEmployeeEntity, templateTask entity.EmployeeTaskEntity) (err error)
	// callback employee task
	CallbackEmployeeTaskForAssessment(ctx context.Context, templateTask entity.EmployeeTaskEntity) (err error)
	// notifier for task reminder
	TaskReminderNotifier(ctx context.Context, task entity.EmployeeTaskEntity) (err error)
	// get employee program tasks
	GetEmployeeProgramTask(ctx context.Context, employee entity.EmployeeTaskEntity) (programs []model.EmployeeProgram, err error)

	InsertEmployeeTaskProgramAssignment(ctx context.Context, assigningMode entity.TypeAssigningMode, relatedID string, assignDate time.Time, programIDArr []string, confirmIgnoreTakenProgram bool) (err error)

	GetEmployeeTaskProgramAssignment(ctx context.Context, assignedProgramPaging *commonmodel.PaginationResponseModel) (err error)

	GetProgramAssignmentHistory(ctx context.Context, programAssignmentHistory *commonmodel.PaginationResponseModel) (err error)

	SearchPrograms(ctx context.Context, keyword string) ([]consumermodel.MProgram, error)

	GetPrograms(ctx context.Context) ([]consumermodel.MProgram, error)

	GetCorporateServices(ctx context.Context, params model.MCorporateServiceRequest) (outputArr []entity.MCorporateService, totalData int, err error)

	GetCorporateServiceByID(ctx context.Context, id string) (entity.MCorporateService, error)

	InsertCorporateServiceBooking(ctx context.Context, serviceID string, scheduledAt time.Time) (err error)

	GetCorporateServiceBooking(ctx context.Context, limit, offset int) ([]entity.CorporateServiceBooking, int, error)

	SearchCorporateService(ctx context.Context, keyword string) ([]entity.MCorporateService, error)

	// insert new well being caused by scheduler trigger
	RenewAllWellBeingTaskSvc(ctx context.Context) (err error)
}
