package grpcclient_feedback

import (
	"context"

	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
)

type FeedbackClientImpl struct {
	sugar  logger.CustomLogger
	client pb.FeedbackServiceClient
}

// constructor function to build new object
func NewFeedbackClient(sugar logger.CustomLogger,
	clientConnection conn.GRPCClient) FeedbackClient {
	// create new connection for client
	client := pb.NewFeedbackServiceClient(clientConnection.GetClient())
	sugar.WithContext(context.Background()).Info("creating grpc connection for client feedback")
	return &FeedbackClientImpl{
		sugar:  sugar,
		client: client,
	}
}

// get client registration for grpc
func (c *FeedbackClientImpl) GetClient() pb.FeedbackServiceClient {
	c.sugar.WithContext(context.Background()).Info("getting client for feedback connection")
	return c.client
}
