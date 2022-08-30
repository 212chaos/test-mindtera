package bookingserviceroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	emptaskhandler "github.com/mindtera/corporate-service/handler/http/employee-task"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	internalmidware "github.com/mindtera/go-common-module/common/middleware/internal-request-middleware"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type BookingServiceRouteImpl struct {
	ginGroup            ggroup.GinGroup
	internalMidware     internalmidware.InternalRequestMiddleware
	authMiddleware      authmidware.AuthMiddleware
	employeeTaskHandler emptaskhandler.EmployeeTaskHandler
}

func NewBookingServiceRouter(
	ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	internalMidware internalmidware.InternalRequestMiddleware,
	employeeTaskHandler emptaskhandler.EmployeeTaskHandler) BookingServiceRoute {
	return &BookingServiceRouteImpl{
		ginGroup:            ginRouter.GROUP(fmt.Sprintf("%v/mindtera-service", common.API_URL)),
		authMiddleware:      authMiddleware,
		employeeTaskHandler: employeeTaskHandler,
		internalMidware:     internalMidware,
	}
}

func (e *BookingServiceRouteImpl) bookingServices() {
	e.ginGroup.GET("",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.GetCorporateServices,
	)

	e.ginGroup.POST("",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.InsertCorporateServiceBooking)

	e.ginGroup.GET("/:corporate_service_id",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.GetCorporateServiceByID,
	)

	e.ginGroup.GET("/search-service",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.SearchCorporateService,
	)

	e.ginGroup.GET("/my-bookings",
		e.authMiddleware.AuthUserCorporate,
		e.employeeTaskHandler.GetCurrentCorporateServiceBooking,
	)
}

func (e *BookingServiceRouteImpl) Routes() {
	e.bookingServices()
}
