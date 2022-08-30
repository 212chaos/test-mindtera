package repo_feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
)

type FeedbackRepository interface {
	GetFeedbackByFilter(ctx context.Context, feedbackFilter model.FeedbackFilter, pagingMeta *mdl.PaginationModel, feedback *[]entity.Feedback) (err error)
	UpsertFeedback(ctx context.Context, Feedback *entity.Feedback) (err error)
	MarkFeedbackAsRead(ctx context.Context, feedbackId, companyId uuid.UUID) (err error)
	GetFeedbackCategory(ctx context.Context, feedbackCategory *[]entity.FeedbackCategory) (err error)
	GetFeedbackShownCategory(ctx context.Context, feedbackCategory *[]entity.FeedbackShownCategory) (err error)
}
