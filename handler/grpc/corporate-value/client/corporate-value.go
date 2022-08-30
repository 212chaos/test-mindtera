package grpcclient_corporatevalue

import pb "github.com/mindtera/go-common-module/common/pb"

type CorporateValueClient interface {
	// get client registration for grpc
	GetClient() pb.CorporateValueServiceClient
}
