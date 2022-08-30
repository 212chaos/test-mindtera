package http_employeetask

import "github.com/gin-gonic/gin"

type EmployeeTaskHandler interface {
	// get employee task by query
	GetEmployeeTask(ctx *gin.Context)

	// upsert employee task in batch
	UpsertEmployeeTask(ctx *gin.Context)

	// notifier employee task in batch
	NotifyEmployeeTask(ctx *gin.Context)

	// callback employee task from scheduler
	CallbackEmployeeTask(ctx *gin.Context)

	InsertEmployeeTaskProgramAssignment(ctx *gin.Context)

	GetEmployeeTaskProgramAssignment(ctx *gin.Context)

	SearchProgram(ctx *gin.Context)

	GetCorporateServices(ctx *gin.Context)

	GetCorporateServiceByID(ctx *gin.Context)

	InsertCorporateServiceBooking(ctx *gin.Context)

	GetCurrentCorporateServiceBooking(ctx *gin.Context)

	SearchCorporateService(ctx *gin.Context)

	RenewWellBeingInternal(ctx *gin.Context)
}
