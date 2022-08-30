package grpcclient_corporateemployee

import pb "github.com/mindtera/go-common-module/common/pb"

type CorporateEmployeeClient interface {
	// get client registration for grpc
	GetClient() pb.CorporateEmployeeServiceClient
}
