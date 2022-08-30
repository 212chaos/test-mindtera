package grpcserver_corporatevalue

import (
	"context"

	"github.com/mindtera/go-common-module/common/pb"
)

type CorporateValueServer interface {
	GetCorpValueIdByCorpId(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValueIds, err error)
	GetCorpValueDetailByCorpId(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValues, err error)
	GetCorpValueName(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValueNameArr, err error)
}
