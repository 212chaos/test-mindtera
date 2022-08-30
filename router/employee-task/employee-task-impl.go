package employeetaskroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	emptaskhandler "github.com/mindtera/corporate-service/handler/http/employee-task"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	internalmidware "github.com/mindtera/go-common-module/common/middleware/internal-request-middleware"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type EmployeeTaskRouteImpl struct {
	ginGroup            ggroup.GinGroup
	internalMidware     internalmidware.InternalRequestMiddleware
	authMiddleware      authmidware.AuthMiddleware
	employeeTaskHandler emptaskhandler.EmployeeTaskHandler
}

func NewEmployeeTaskRouter(
	ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	internalMidware internalmidware.InternalRequestMiddleware,
	employeeTaskHandler emptaskhandler.EmployeeTaskHandler) EmployeeTaskRoute {
	return &EmployeeTaskRouteImpl{
		ginGroup:            ginRouter.GROUP(fmt.Sprintf("%v/employee-task", common.API_URL)),
		authMiddleware:      authMiddleware,
		employeeTaskHandler: employeeTaskHandler,
		internalMidware:     internalMidware,
	}
}

func (e *EmployeeTaskRouteImpl) getEmployeeTask() {
	e.ginGroup.GET("",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.GetEmployeeTask)

	// for internal communication
	e.ginGroup.GET("/internal",
		e.internalMidware.CheckServiceKey,
		e.employeeTaskHandler.GetEmployeeTask)
}

func (e *EmployeeTaskRouteImpl) postEmployeeTask() {
	e.ginGroup.POST("",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.UpsertEmployeeTask)

	// for internal communication
	e.ginGroup.POST("/internal",
		e.internalMidware.CheckServiceKey,
		e.employeeTaskHandler.UpsertEmployeeTask)

	// notifier
	e.ginGroup.POST("/notifier",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.NotifyEmployeeTask)

	// renew all
	e.ginGroup.PUT("/renew-well-being",
		e.authMiddleware.AuthUser,
		e.authMiddleware.AdminRequest,
		e.employeeTaskHandler.RenewWellBeingInternal)
}

func (e *EmployeeTaskRouteImpl) employeeProgramAssignment() {
	e.ginGroup.POST("/assign-program",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.InsertEmployeeTaskProgramAssignment)
	e.ginGroup.GET("/assign-program",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.GetEmployeeTaskProgramAssignment,
	)
	e.ginGroup.GET("/search-program",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.SearchProgram,
	)
}

func (e *EmployeeTaskRouteImpl) Routes() {
	e.getEmployeeTask()
	e.postEmployeeTask()
	e.employeeProgramAssignment()
}
