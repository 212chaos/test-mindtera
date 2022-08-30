package grpcserver_feedback

import (
	"context"

	"github.com/mindtera/go-common-module/common/pb"
)

type FeedbackServer interface {
	UpsertFeedback(context.Context, *pb.Feedback) (*pb.Empty, error)
	GetFeedbackCategory(context.Context, *pb.Empty) (*pb.FeedbackCategories, error)
	GetFeedbackShownCategory(context.Context, *pb.Empty) (*pb.FeedbackShownCategories, error)
}
