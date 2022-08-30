package grpcclient_corporatevalue

import (
	"context"

	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
)

type CorporateValueClientImpl struct {
	sugar  logger.CustomLogger
	client pb.CorporateValueServiceClient
}

// constructor function to build new object
func NewCorporateValueClient(sugar logger.CustomLogger, clientConnection conn.GRPCClient) CorporateValueClient {
	// create new connection for client
	client := pb.NewCorporateValueServiceClient(clientConnection.GetClient())
	sugar.WithContext(context.Background()).Info("creating grpc connection for client corporate value")
	return &CorporateValueClientImpl{
		sugar:  sugar,
		client: client,
	}
}

// get client registration for grpc
func (c *CorporateValueClientImpl) GetClient() pb.CorporateValueServiceClient {
	c.sugar.WithContext(context.Background()).Info("getting client for corporate employee connection")
	return c.client
}
