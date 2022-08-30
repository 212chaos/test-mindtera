package corporateemployeeroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	corporateemployeehandler "github.com/mindtera/corporate-service/handler/http/corporate-employee"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type CorporateEmployeeRouteImpl struct {
	ginGroup                 ggroup.GinGroup
	authMiddleware           authmidware.AuthMiddleware
	corporateEmployeeHandler corporateemployeehandler.CorporateEmployeeHandler
}

func NewCorporateEmployeeRouter(ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	corporateEmployeeHandler corporateemployeehandler.CorporateEmployeeHandler) CorporateEmployeeRoute {
	return &CorporateEmployeeRouteImpl{
		authMiddleware:           authMiddleware,
		corporateEmployeeHandler: corporateEmployeeHandler,
		ginGroup:                 ginRouter.GROUP(fmt.Sprintf("%v/corporate-employee", common.API_URL)),
	}
}

func (c *CorporateEmployeeRouteImpl) getCorporateEmployee() {
	c.ginGroup.GET("", c.authMiddleware.AuthUserCorporate,
		c.corporateEmployeeHandler.GetCorporateEmployeeHandler)
	c.ginGroup.GET("/public", c.authMiddleware.AuthUserPublic,
		c.corporateEmployeeHandler.GetCorporateEmployeePublicIDHandler)
}

func (c *CorporateEmployeeRouteImpl) upsertCorporateEmployee() {
	c.ginGroup.POST("", c.authMiddleware.AuthUserCorporate,
		c.corporateEmployeeHandler.UpsertCorporateEmployeeHandler)
}

func (c *CorporateEmployeeRouteImpl) deleteCorporateEmployee() {
	c.ginGroup.DELETE("", c.authMiddleware.AuthUserCorporate,
		c.corporateEmployeeHandler.DeleteCorporateEmployeeHandler)
}

func (c *CorporateEmployeeRouteImpl) Routes() {
	c.getCorporateEmployee()
	c.upsertCorporateEmployee()
	c.deleteCorporateEmployee()
}
