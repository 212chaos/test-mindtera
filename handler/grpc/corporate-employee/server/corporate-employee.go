package grpsserver_corporateemployee

import (
	"context"

	pb "github.com/mindtera/go-common-module/common/pb"
)

type CorporateEmployeeServer interface {
	UpdateCorporatePublicUser(ctx context.Context, publicUser *pb.PublicUser) (result *pb.PublicUser, err error)
	GetCorporatePublicUserByEmail(ctx context.Context, publicUser *pb.PublicUser) (result *pb.CorporateEmployee, err error)
}
