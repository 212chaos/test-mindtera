package grpcclient_employeetask

import pb "github.com/mindtera/go-common-module/common/pb"

type EmployeeTaskClient interface {
	// get client registration for grpc
	GetClient() pb.EmployeeTaskServiceClient
}
