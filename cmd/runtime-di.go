package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/mindtera/corporate-service/common"
	corpemplsvr "github.com/mindtera/corporate-service/handler/grpc/corporate-employee/server"
	corpsubssvr "github.com/mindtera/corporate-service/handler/grpc/corporate-subscription/server"
	corpvalsvr "github.com/mindtera/corporate-service/handler/grpc/corporate-value/server"
	empltasksvr "github.com/mindtera/corporate-service/handler/grpc/employee-task/server"
	feedbacksvr "github.com/mindtera/corporate-service/handler/grpc/feedback/server"
	logger "github.com/mindtera/go-common-module/common/logger"

	corphandler "github.com/mindtera/corporate-service/handler/http/corporate"
	corpemplhandler "github.com/mindtera/corporate-service/handler/http/corporate-employee"
	corpvalhandler "github.com/mindtera/corporate-service/handler/http/corporate-value"
	corpvalrelhandler "github.com/mindtera/corporate-service/handler/http/corporate-value-relation"
	empltaskhandler "github.com/mindtera/corporate-service/handler/http/employee-task"
	feedbackhandler "github.com/mindtera/corporate-service/handler/http/feedback"

	corpassessrelrepo "github.com/mindtera/corporate-service/repository/corporate-assessment-relation"
	corpemplrepo "github.com/mindtera/corporate-service/repository/corporate-employee"
	corpprogramrepo "github.com/mindtera/corporate-service/repository/corporate-program"
	corpsmmrepo "github.com/mindtera/corporate-service/repository/corporate-smm"
	corpsubsrepo "github.com/mindtera/corporate-service/repository/corporate-subscription"
	corpvalrepo "github.com/mindtera/corporate-service/repository/corporate-value"
	cforpvalrelrepo "github.com/mindtera/corporate-service/repository/corporate-value-relation"
	empltaskrepo "github.com/mindtera/corporate-service/repository/employee-task"
	feedbackrepo "github.com/mindtera/corporate-service/repository/feedback"
	usersubs "github.com/mindtera/corporate-service/repository/user-subscription"

	// get from dashboard
	corprolerepo "github.com/mindtera/dashboard-auth/repository/corporate-role"
	corpdeptrepo "github.com/mindtera/dashboard-auth/repository/department"
	corprolesvc "github.com/mindtera/dashboard-auth/service/corporate-role"
	corpdeptsvc "github.com/mindtera/dashboard-auth/service/department"

	bookingservicerouter "github.com/mindtera/corporate-service/router/booking-service"
	corprouter "github.com/mindtera/corporate-service/router/corporate"
	corpemplrouter "github.com/mindtera/corporate-service/router/corporate-employee"
	corporatevaluerouter "github.com/mindtera/corporate-service/router/corporate-value"
	corpvalrelrouter "github.com/mindtera/corporate-service/router/corporate-value-relation"
	employeetaskrouter "github.com/mindtera/corporate-service/router/employee-task"
	feedbackrouter "github.com/mindtera/corporate-service/router/feedback"
	schedullerrouter "github.com/mindtera/corporate-service/router/scheduler-callback"

	corpemplsvc "github.com/mindtera/corporate-service/service/corporate-employee"
	corpsubssvc "github.com/mindtera/corporate-service/service/corporate-subscription"
	corpvalsvc "github.com/mindtera/corporate-service/service/corporate-value"
	corpvalrelsvc "github.com/mindtera/corporate-service/service/corporate-value-relation"
	empltasksvc "github.com/mindtera/corporate-service/service/employee-task"
	feedbacksvc "github.com/mindtera/corporate-service/service/feedback"
	usersubssvc "github.com/mindtera/corporate-service/service/user-subscription"

	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	middlewareexternal "github.com/mindtera/dashboard-auth/middleware/external"
	authcommonsvc "github.com/mindtera/dashboard-auth/service/common"

	mailjet "github.com/mindtera/go-common-module/common/client/mailjet-client"
	schedulerclient "github.com/mindtera/go-common-module/common/client/scheduler-client"
	corsmiddleware "github.com/mindtera/go-common-module/common/middleware/cors-middleware"
	internalmidware "github.com/mindtera/go-common-module/common/middleware/internal-request-middleware"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
	googlestorage "github.com/mindtera/go-common-module/common/v2/configuration/google-cloud"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	grpcserver "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	redisconf "github.com/mindtera/go-common-module/common/v2/configuration/redis"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonservice "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
	util "github.com/mindtera/go-common-module/common/v2/service/util"
)

// initiate all grouped DI

func repoDependencies() []any {
	return []any{
		empltaskrepo.NewEmployeeTaskRepository, corprolerepo.NewCorporateRoleRepository,
		corpvalrepo.NewCorporateValueRepository, corpassessrelrepo.NewCorporateAssessmentRelationRepo,
		cforpvalrelrepo.NewCorporateValueRelationRepo, usersubs.NewUserSubscriptionRepository,
		corpsmmrepo.NewCorporateSMMRepository, corpemplrepo.NewCorporateEmployeeRepository,
		corpsubsrepo.NewCorporateSubscriptionRepository, commonservice.NewCommonService,
		corpdeptrepo.NewCorporateDepartmentRepository, corpprogramrepo.NewCorporateProgram,
		feedbackrepo.NewFeedbackRepository}
}

func svcDependencies() []any {
	return []any{
		corpdeptsvc.NewDepartmentService, corpsubssvc.NewCorporateSubscriptionService,
		corprolesvc.NewCorporateRoleService, authcommonsvc.NewCommonService,
		corpvalsvc.NewCorporateValueService, corpvalrelsvc.NewCorporateValueRelationService,
		corpemplsvc.NewCorporateEmployeeService, empltasksvc.NewEmployeeTaskService,
		usersubssvc.NewPublicUserSubscriptionService, feedbacksvc.NewFeedbackSvc}
}

func commonDependencies() []any {
	return []any{
		logger.NewWrappedZapLogger, logger.NewCustomLogger, googlestorage.NewGoogleStorageConfig,
		gormpg.NewPostgresConfig, redisconf.NewRedisConfig,
		schedulerclient.NewSchedulerClient, mailjet.NewMailjetClient,
		redissvc.NewRedisSvc, assert.NewAssert,
		util.NewUtil, trx.NewTransaction}
}

func handlerDependencies() []any {
	return []any{
		corsmiddleware.NewCorsMiddleware, internalmidware.NewInternalRequestMiddleware,
		// authmidware.NewAuthMiddleware,
		middlewareexternal.NewAuthMiddlewareExternal,
		corpvalhandler.NewCorporateValueHandler,
		corpvalrelhandler.NewCorporateValueRelationHandler, corpemplhandler.NewCorporateEmployeeHandler,
		corphandler.NewCorporateHandler, empltaskhandler.NewEmployeeTaskService,
		corpemplsvr.NewCorporateEmployeeServer, empltasksvr.NewEmployeeTaskServer,
		grpcserver.NewGRPCServer, ginrouter.NewGinRouter, corpvalsvr.NewCoporateValueServer,
		feedbackhandler.NewFeedbackHttpHandler, feedbacksvr.NewFeedbackServer,
		corpsubssvr.NewCorporateSubscriptionSrv}

}

func BuildInRuntime() (g ginrouter.GinRouter, grpcServer grpcserver.GRPCServer, err error) {
	c := dig.New()
	// define all generic
	var constructor []any
	constructor = append(constructor, commonDependencies()...)
	constructor = append(constructor, repoDependencies()...)
	constructor = append(constructor, svcDependencies()...)
	constructor = append(constructor, handlerDependencies()...)

	// provide all generic
	for _, service := range constructor {
		if err := c.Provide(service); err != nil {
			return nil, nil, err
		}
	}
	// invoked function needed
	if err = c.Invoke(func(
		gr ginrouter.GinRouter,
		grpc grpcserver.GRPCServer,
		cors corsmiddleware.CorsMiddleware,
		authMidware authmidware.AuthMiddleware,
		internalMidware internalmidware.InternalRequestMiddleware,
		corpValHandler corpvalhandler.CorporateValueHandler,
		corpValRelHandler corpvalrelhandler.CorporateValueRelationHandler,
		corpEmplHandler corpemplhandler.CorporateEmployeeHandler,
		corpHandler corphandler.CorporateHandler,
		emplTaskHandler empltaskhandler.EmployeeTaskHandler,
		corpEmplSvr corpemplsvr.CorporateEmployeeServer,
		emplTaskSvr empltasksvr.EmployeeTaskServer,
		corpValSvr corpvalsvr.CorporateValueServer,
		corpSubsSvr corpsubssvr.CorporateSubscriptionSrv,
		feedbackSvr feedbacksvr.FeedbackServer,
		feedbackHdl feedbackhandler.FeedbackHttpHandler) {
		//populate result
		g = gr
		grpcServer = grpc

		// init middleware etc
		g.USE(gin.Recovery(),
			gin.Logger(),
			cors.CommonRequest)

		// health check
		g.GET("/", func(ctx *gin.Context) {
			ctx.JSONP(http.StatusOK, map[string]any{
				"success":          true,
				"application-name": common.SERVICE_NAME,
				"access-time":      time.Now()})
		})

		//corporate route
		corpemplrouter.NewCorporateEmployeeRouter(g, authMidware, corpEmplHandler).Routes()
		corpvalrelrouter.NewCorporateValueRelationRouter(g, authMidware, corpValRelHandler).Routes()
		corporatevaluerouter.NewCorporateValueRouter(g, authMidware, corpValHandler).Routes()
		corprouter.NewCorporateRouter(g, authMidware, corpHandler).Routes()
		schedullerrouter.NewSchedulerCallbackRouter(g, authMidware, corpHandler, internalMidware, emplTaskHandler).Routes()
		employeetaskrouter.NewEmployeeTaskRouter(g, authMidware, internalMidware, emplTaskHandler).Routes()
		bookingservicerouter.NewBookingServiceRouter(g, authMidware, internalMidware, emplTaskHandler).Routes()
		feedbackrouter.NewFeedbackRouter(g, authMidware, feedbackHdl).Router()
	}); err != nil {
		panic(err)
	}
	return g, grpcServer, err
}
