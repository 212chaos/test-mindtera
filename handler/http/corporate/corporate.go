package http_corporate

import "github.com/gin-gonic/gin"

type CorporateHandler interface {
	GetCorporateStatus(ctx *gin.Context)
	UpsertCorporateSubscription(ctx *gin.Context)
	GetCorporateDashboardResult(ctx *gin.Context)
	// Callback expiration subscription
	CallbackExpiredCorporateSubscription(ctx *gin.Context)
	// get corporate assessment result status
	GetAssessmentStatus(ctx *gin.Context)
	// get corporate assessment result status
	GetWellBeingStatus(ctx *gin.Context)
}
