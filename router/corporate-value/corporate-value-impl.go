package corporatevalueroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	corporatevaluehandler "github.com/mindtera/corporate-service/handler/http/corporate-value"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type CorporateValueRouteImpl struct {
	ginGroup              ggroup.GinGroup
	authMiddleware        authmidware.AuthMiddleware
	corporateValueHandler corporatevaluehandler.CorporateValueHandler
}

func NewCorporateValueRouter(ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	corporateValueHandler corporatevaluehandler.CorporateValueHandler) CorporateValueRoute {
	return &CorporateValueRouteImpl{
		ginGroup:              ginRouter.GROUP(fmt.Sprintf("%v/corporate-value", common.API_URL)),
		authMiddleware:        authMiddleware,
		corporateValueHandler: corporateValueHandler,
	}
}

func (c *CorporateValueRouteImpl) getCorporateValues() {
	c.ginGroup.GET("",
		c.authMiddleware.AuthUserCorporate,
		c.corporateValueHandler.GetCorporateValueHandler)
	c.ginGroup.GET("/smm",
		c.authMiddleware.AuthUserCorporate,
		c.corporateValueHandler.GetCorporateValueBySMMHandler)
}

func (c *CorporateValueRouteImpl) Routes() {
	c.getCorporateValues()
}
