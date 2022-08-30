package grpcclient_feedback

import pb "github.com/mindtera/go-common-module/common/pb"

type FeedbackClient interface {
	// get client registration for grpc
	GetClient() pb.FeedbackServiceClient
}
