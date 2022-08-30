package http_corporatevaluerelation

import "github.com/gin-gonic/gin"

type CorporateValueRelationHandler interface {
	GetCorporateValueRelationHandler(ctx *gin.Context)
	GetCorporateValueRelationPublicHandler(ctx *gin.Context)
	UpsertCorporateValueRelationHandler(ctx *gin.Context)
	CalculateCorporateValueRelationHandler(ctx *gin.Context)
	DeleteCorporateValueRelationHandler(ctx *gin.Context)
}
