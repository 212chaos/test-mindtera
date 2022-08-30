package corporatevaluerelationroute

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	corpvalrelhandler "github.com/mindtera/corporate-service/handler/http/corporate-value-relation"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type CorporateValueRelationRouteImpl struct {
	ginGroup          ggroup.GinGroup
	ginGroupPublic    ggroup.GinGroup
	authMiddleware    authmidware.AuthMiddleware
	corpValRelHandler corpvalrelhandler.CorporateValueRelationHandler
}

func NewCorporateValueRelationRouter(
	ginGroup ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	corpValRelHandler corpvalrelhandler.CorporateValueRelationHandler) CorporateValueRelationRoute {
	return &CorporateValueRelationRouteImpl{
		ginGroup:          ginGroup.GROUP(fmt.Sprintf("%v/corporate-value-relation", common.API_URL)),
		ginGroupPublic:    ginGroup.GROUP(fmt.Sprintf("%v/public/corporate-value-relation", common.API_PUBLIC_URL)),
		authMiddleware:    authMiddleware,
		corpValRelHandler: corpValRelHandler,
	}
}

func (c *CorporateValueRelationRouteImpl) getCorporateValueRelation() {
	c.ginGroup.GET("", c.authMiddleware.AuthUserCorporate,
		c.corpValRelHandler.GetCorporateValueRelationHandler)
	c.ginGroupPublic.GET("", c.authMiddleware.AuthUserPublic,
		c.corpValRelHandler.GetCorporateValueRelationPublicHandler)
}

func (c *CorporateValueRelationRouteImpl) upsertCorporateValueRelation() {
	c.ginGroup.POST("", c.authMiddleware.AuthUserCorporate,
		c.corpValRelHandler.UpsertCorporateValueRelationHandler)
	c.ginGroup.POST("/calculate-result", c.authMiddleware.AuthUserCorporate,
		c.corpValRelHandler.CalculateCorporateValueRelationHandler)
}

func (c *CorporateValueRelationRouteImpl) deleteCorporateValueRelation() {
	c.ginGroup.DELETE("", c.corpValRelHandler.DeleteCorporateValueRelationHandler)
}

func (c *CorporateValueRelationRouteImpl) calculateCorporateValueRelation() {

}

func (c *CorporateValueRelationRouteImpl) Routes() {
	c.getCorporateValueRelation()
	c.upsertCorporateValueRelation()
	c.calculateCorporateValueRelation()
	c.deleteCorporateValueRelation()
}
