package grpcclient_corporateemployee

import (
	"context"

	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
)

type CorporateEmployeeClientImpl struct {
	sugar  logger.CustomLogger
	client pb.CorporateEmployeeServiceClient
}

// constructor function to build new object
func NewCorporateEmployeeClient(sugar logger.CustomLogger, clientConnection conn.GRPCClient) CorporateEmployeeClient {
	// create new connection for client
	client := pb.NewCorporateEmployeeServiceClient(clientConnection.GetClient())
	sugar.WithContext(context.Background()).Info("creating grpc connection for client corporate employee")
	return &CorporateEmployeeClientImpl{
		sugar:  sugar,
		client: client,
	}
}

// get client registration for grpc
func (c *CorporateEmployeeClientImpl) GetClient() pb.CorporateEmployeeServiceClient {
	c.sugar.WithContext(context.Background()).Info("getting client for corporate employee connection")
	return c.client
}
