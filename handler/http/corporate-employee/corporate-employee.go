package http_corporateemployee

import "github.com/gin-gonic/gin"

type CorporateEmployeeHandler interface {
	GetCorporateEmployeeHandler(ctx *gin.Context)
	GetCorporateEmployeePublicIDHandler(ctx *gin.Context)
	UpsertCorporateEmployeeHandler(ctx *gin.Context)
	DeleteCorporateEmployeeHandler(ctx *gin.Context)
}
