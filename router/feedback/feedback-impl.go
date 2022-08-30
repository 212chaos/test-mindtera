package route_feedback

import (
	"fmt"

	"github.com/mindtera/corporate-service/common"
	feedbackhdl "github.com/mindtera/corporate-service/handler/http/feedback"
	authmidware "github.com/mindtera/dashboard-auth/middleware/auth"
	ggroup "github.com/mindtera/go-common-module/common/v2/configuration/gin/group"
	ginrouter "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

type FeedbackRouterImpl struct {
	ginGroup       ggroup.GinGroup
	authMiddleware authmidware.AuthMiddleware
	feedbackHdl    feedbackhdl.FeedbackHttpHandler
}

func NewFeedbackRouter(ginRouter ginrouter.GinRouter,
	authMiddleware authmidware.AuthMiddleware,
	feedbackHdl feedbackhdl.FeedbackHttpHandler) FeedbackRouter {
	return &FeedbackRouterImpl{
		ginGroup:       ginRouter.GROUP(fmt.Sprintf("%v/feedback", common.API_URL)),
		authMiddleware: authMiddleware,
		feedbackHdl:    feedbackHdl,
	}
}

func (f *FeedbackRouterImpl) getRequest() {
	f.ginGroup.GET("", f.authMiddleware.AuthUserCorporate,
		f.feedbackHdl.GetFeedbackHdl)
	f.ginGroup.GET("/category", f.authMiddleware.AuthUserCorporate,
		f.feedbackHdl.GetFeedbackCategoryHdl)
	f.ginGroup.GET("/shown-category", f.authMiddleware.AuthUserCorporate,
		f.feedbackHdl.GetFeedbackShownCategoryHdl)
}

func (f *FeedbackRouterImpl) putRequest() {
	f.ginGroup.PUT("/:feedback_id", f.authMiddleware.AuthUserCorporate,
		f.feedbackHdl.MarkFeedbackAsReadHdl)
}

func (f *FeedbackRouterImpl) Router() {
	f.getRequest()
	f.putRequest()
}
