package corporateroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	corporatehandler "github.com/mindtera/corporate-service/handler/http/corporate"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	gingroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type CorporateRouteImpl struct {
	ginGroup         gingroup.GinGroup
	authMiddleware   authmidware.AuthMiddleware
	corporateHandler corporatehandler.CorporateHandler
}

func NewCorporateRouter(ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	corporateHandler corporatehandler.CorporateHandler) CorporateRoute {
	return &CorporateRouteImpl{
		ginGroup:         ginRouter.GROUP(fmt.Sprintf("%v/corporate", common.API_URL)),
		corporateHandler: corporateHandler,
		authMiddleware:   authMiddleware}
}

func (c *CorporateRouteImpl) getCorporateStatus() {
	c.ginGroup.GET("/status", c.authMiddleware.AuthUserCorporate, c.corporateHandler.GetCorporateStatus)
	c.ginGroup.GET("/assessment-status", c.authMiddleware.AuthUserCorporate, c.corporateHandler.GetAssessmentStatus)
	c.ginGroup.GET("/well-being-status", c.authMiddleware.AuthUserCorporate, c.corporateHandler.GetWellBeingStatus)
	c.ginGroup.GET("/dashboard-result", c.authMiddleware.AuthUserCorporate, c.corporateHandler.GetCorporateDashboardResult)
}

func (c *CorporateRouteImpl) upsertCorporateSubscription() {
	c.ginGroup.POST("/subscription", c.authMiddleware.AuthUser,
		c.authMiddleware.AdminRequest, c.corporateHandler.UpsertCorporateSubscription)
}

func (c *CorporateRouteImpl) Routes() {
	c.getCorporateStatus()
	c.upsertCorporateSubscription()
}
