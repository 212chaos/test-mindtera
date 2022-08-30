package grpcclient_employeetask

import (
	"context"

	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
)

type EmployeeTaskClientImpl struct {
	sugar  logger.CustomLogger
	client pb.EmployeeTaskServiceClient
}

// constructor function to build new object
func NewEmployeeTaskClient(sugar logger.CustomLogger,
	clientConnection conn.GRPCClient) EmployeeTaskClient {
	// create new connection for client
	client := pb.NewEmployeeTaskServiceClient(clientConnection.GetClient())
	sugar.WithContext(context.Background()).Info("creating grpc connection for client corporate employee")
	return &EmployeeTaskClientImpl{
		sugar:  sugar,
		client: client,
	}
}

// get client registration for grpc
func (c *EmployeeTaskClientImpl) GetClient() pb.EmployeeTaskServiceClient {
	c.sugar.WithContext(context.Background()).Info("getting client for corporate employee connection")
	return c.client
}
