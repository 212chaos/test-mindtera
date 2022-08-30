package grpsserver_corporateemployee

import (
	"context"

	pb "github.com/mindtera/go-common-module/common/pb"
)

type EmployeeTaskServer interface {
	UpsertEmployeeTask(ctx context.Context, empTask *pb.EmployeeTaskUpdate) (result *pb.Empty, err error)
	GetEmployeeTask(ctx context.Context, empTask *pb.EmployeeTaskRequst) (result *pb.EmployeeTasks, err error)
	CreateEmployeeTask(ctx context.Context, in *pb.EmployeeTaskCreate) (*pb.Empty, error)
	GetEmployeeProgram(ctx context.Context, in *pb.ProgramQuery) (*pb.Programs, error)
	RenewWellBeingTrigger(ctx context.Context, in *pb.Empty) (result *pb.Empty, err error)
}
