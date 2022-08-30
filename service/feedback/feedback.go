package svc_feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
)

type Feedback interface {
	GetFeedbackByFilterSvc(ctx context.Context, feedbackFilter model.FeedbackFilter, feedbackPaging *mdl.PaginationResponseModel) (err mdl.ErrorMessage)
	UpsertFeedbackSvc(ctx context.Context, Feedback *entity.Feedback) (err mdl.ErrorMessage)
	MarkFeedbackAsReadSvc(ctx context.Context, feedbackId uuid.UUID) (err mdl.ErrorMessage)
	GetFeedbackCategorySvc(ctx context.Context, feedbackCategory *[]entity.FeedbackCategory) (err mdl.ErrorMessage)
	GetFeedbackShownCategorySvc(ctx context.Context, feedbackCategory *[]entity.FeedbackShownCategory) (err mdl.ErrorMessage)
}
