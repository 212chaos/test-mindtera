package grpcserver_corporatesubscription

import (
	"context"

	"github.com/mindtera/go-common-module/common/pb"
)

type CorporateSubscriptionSrv interface {
	GetCorporateSubscription(context.Context, *pb.Corporate) (*pb.Empty, error)
}
