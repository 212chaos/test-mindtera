package schedulercallbackroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	corpsvchandler "github.com/mindtera/corporate-service/handler/http/corporate"
	emptaskhandler "github.com/mindtera/corporate-service/handler/http/employee-task"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	internalrequestmiddleware "github.com/mindtera/go-common-module/common/middleware/internal-request-middleware"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type SchedulerCallbackRouterImpl struct {
	ginGroup        ggroup.GinGroup
	authMiddleware  authmidware.AuthMiddleware
	corpSvcHandler  corpsvchandler.CorporateHandler
	empTaskHandler  emptaskhandler.EmployeeTaskHandler
	internalMidware internalrequestmiddleware.InternalRequestMiddleware
}

// route constructor
func NewSchedulerCallbackRouter(
	ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	corpSvcHandler corpsvchandler.CorporateHandler,
	internalMidware internalrequestmiddleware.InternalRequestMiddleware,
	empTaskHandler emptaskhandler.EmployeeTaskHandler) SchedulerCallbackRoute {
	return &SchedulerCallbackRouterImpl{
		ginGroup:        ginRouter.GROUP(fmt.Sprintf("%v/scheduler", common.API_URL)),
		authMiddleware:  authMiddleware,
		corpSvcHandler:  corpSvcHandler,
		internalMidware: internalMidware,
		empTaskHandler:  empTaskHandler}
}

func (s *SchedulerCallbackRouterImpl) postAssessmentCallback() {
	s.ginGroup.POST("/assessment-callback", s.internalMidware.CheckServiceKey,
		s.empTaskHandler.CallbackEmployeeTask)
}
func (s *SchedulerCallbackRouterImpl) postExpiredSubscriptionCallback() {
	s.ginGroup.POST("/corporate-end-period", s.internalMidware.CheckServiceKey,
		s.corpSvcHandler.CallbackExpiredCorporateSubscription)
}

func (s *SchedulerCallbackRouterImpl) Routes() {
	s.postAssessmentCallback()
	s.postExpiredSubscriptionCallback()
}
