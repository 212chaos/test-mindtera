package grpcclient_corporatesubscription

import pb "github.com/mindtera/go-common-module/common/pb"

type CorporateSubscriptionClient interface {
	// get client registration for grpc
	GetClient() pb.CorporateSubscriptionServiceClient
}
