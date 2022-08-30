package http_feedback

import "github.com/gin-gonic/gin"

type FeedbackHttpHandler interface {
	MarkFeedbackAsReadHdl(ctx *gin.Context)
	GetFeedbackHdl(ctx *gin.Context)
	GetFeedbackCategoryHdl(ctx *gin.Context)
	GetFeedbackShownCategoryHdl(ctx *gin.Context)
}
