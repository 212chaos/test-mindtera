package http_corporatevalue

import "github.com/gin-gonic/gin"

type CorporateValueHandler interface {
	GetCorporateValueHandler(ctx *gin.Context)
	GetCorporateValueBySMMHandler(ctx *gin.Context)
}
