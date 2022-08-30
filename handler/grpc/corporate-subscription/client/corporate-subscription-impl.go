package grpcclient_corporatesubscription

import (
	"context"

	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
)

type CorporateSubscriptionClientImpl struct {
	sugar  logger.CustomLogger
	client pb.CorporateSubscriptionServiceClient
}

// constructor function to build new object
func NewCorporateSubscriptionClient(sugar logger.CustomLogger, clientConnection conn.GRPCClient) CorporateSubscriptionClient {
	// create new connection for client
	client := pb.NewCorporateSubscriptionServiceClient(clientConnection.GetClient())
	sugar.WithContext(context.Background()).Info("creating grpc connection for client corporate subscription")
	return &CorporateSubscriptionClientImpl{
		sugar:  sugar,
		client: client,
	}
}

// get client registration for grpc
func (c *CorporateSubscriptionClientImpl) GetClient() pb.CorporateSubscriptionServiceClient {
	c.sugar.WithContext(context.Background()).Info("getting client for corporate subscription connection")
	return c.client
}
