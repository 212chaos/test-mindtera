package grpsserver_corporateemployee

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	empltasksvc "github.com/mindtera/corporate-service/service/employee-task"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type EmployeeTaskServerImpl struct {
	sugar       logger.CustomLogger
	trx         trx.Transaction
	grpcServer  conn.GRPCServer
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	emplTaskSvc empltasksvc.EmployeeTaskService
	pb.UnimplementedEmployeeTaskServiceServer
}

func NewEmployeeTaskServer(sugar logger.CustomLogger,
	grpcServer conn.GRPCServer,
	trx trx.Transaction,
	commonSvc commonsvc.CommonService,
	emplTaskSvc empltasksvc.EmployeeTaskService,
	assert assert.Assert) EmployeeTaskServer {
	sugar.WithContext(context.Background()).Info("initiating employee task grpc server")
	result := &EmployeeTaskServerImpl{
		grpcServer:  grpcServer,
		sugar:       sugar,
		assert:      assert,
		trx:         trx,
		commonSvc:   commonSvc,
		emplTaskSvc: emplTaskSvc,
	}
	pb.RegisterEmployeeTaskServiceServer(grpcServer.GetClient(), result)
	return result
}

func (em *EmployeeTaskServerImpl) UpsertEmployeeTask(ctx context.Context, empTask *pb.EmployeeTaskUpdate) (result *pb.Empty, err error) {
	ctx = em.commonSvc.GetCorrelationIdFromGrpc(ctx)
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, em.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := em.trx.GormEndTransaction(ctx); errTrx != nil {
			em.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = em.commonSvc.UpsertCorrelationId(ctx)

	// mapping employee task payload
	endTime := empTask.EndDate.AsTime()

	var taskRelatedId uuid.UUID
	if !em.assert.IsEmpty(empTask.TaskRelatedId) {
		taskRelatedId = uuid.MustParse(empTask.TaskRelatedId)
	}

	emplTasks := []entity.EmployeeTaskEntity{{
		ID:             uuid.MustParse(empTask.Id),
		UserID:         uuid.MustParse(empTask.UserId),
		TaskRelatedID:  taskRelatedId,
		EndDate:        empTask.EndDate.AsTime(),
		StartDate:      empTask.StartDate.AsTime(),
		TaskType:       entity.EMPLOYEE_TASK_TYPE(empTask.TaskType),
		TaskStatus:     entity.EMPLOYEE_TASK_STATUS(empTask.TaskStatus),
		DueDate:        &endTime,
		SubscriptionID: uuid.MustParse(empTask.SubscriptionId),
		Email:          empTask.Email,
		RecordFlag:     empTask.TaskStatus,
	}}
	if err = em.emplTaskSvc.UpsertEmployeeInBatchService(ctx, &emplTasks); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		em.sugar.WithContext(ctx).Errorf("error upsert employee task  err:%v", err.Error())
	}
	return &pb.Empty{}, err
}

func (em *EmployeeTaskServerImpl) GetEmployeeTask(ctx context.Context, empTask *pb.EmployeeTaskRequst) (result *pb.EmployeeTasks, err error) {
	ctx = em.commonSvc.GetCorrelationIdFromGrpc(ctx)
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, em.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := em.trx.GormEndTransaction(ctx); errTrx != nil {
			em.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = em.commonSvc.UpsertCorrelationId(ctx)

	// mapping employee task payload
	task := model.EmployeeTaskQuery{
		UserId:    uuid.MustParse(empTask.UserId),
		EndDate:   empTask.EndDate.AsTime(),
		StartDate: empTask.StartDate.AsTime(),
		Status:    empTask.TaskStatus,
		Types:     strings.Split(empTask.TaskType, ","),
	}

	var emplTasks []entity.EmployeeTaskEntity
	if err = em.emplTaskSvc.GetEmployeeTaskByUserIdAndType(ctx, &task, &emplTasks); err != nil {
		em.sugar.WithContext(ctx).Errorf("error fetching employee task for user id:%v err:%v", task.UserId, err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.EmployeeTasks{}, err
	}
	// mapping result
	var tasks []*pb.EmployeeTask
	for _, v := range emplTasks {

		// check due date
		var ddate *timestamppb.Timestamp
		if v.DueDate != nil {
			ddate = timestamppb.New(*v.DueDate)
		}

		var adate *timestamppb.Timestamp
		if v.AssignDate != nil {
			adate = timestamppb.New(*v.AssignDate)
		}

		tasks = append(tasks, &pb.EmployeeTask{
			Id:             v.ID.String(),
			UserId:         v.UserID.String(),
			SubscriptionId: v.SubscriptionID.String(),
			EndDatePb:      ddate,
			StartDatePb:    adate,
			Email:          v.Email,
			TaskStatus:     string(v.TaskStatus),
			TaskType:       string(v.TaskType),
		})
	}
	// mapping to result
	result = &pb.EmployeeTasks{
		Tasks: tasks,
	}
	return result, err
}

func (em *EmployeeTaskServerImpl) CreateEmployeeTask(ctx context.Context, in *pb.EmployeeTaskCreate) (result *pb.Empty, err error) {
	ctx = em.commonSvc.GetCorrelationIdFromGrpc(ctx)
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, em.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := em.trx.GormEndTransaction(ctx); errTrx != nil {
			em.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = em.commonSvc.UpsertCorrelationId(ctx)

	// mapping task
	var task entity.EmployeeTaskEntity
	if err = em.commonSvc.ObjectMapper(in, &task); err != nil {
		em.sugar.WithContext(ctx).Errorf("error mapping task")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.Empty{}, err
	}
	dueDate := in.GetDueDatePb().AsTime()
	task.DueDate = &dueDate

	// upsert employee task
	if err = em.emplTaskSvc.UpsertEmployeeInBatchService(ctx, &[]entity.EmployeeTaskEntity{task}); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		em.sugar.WithContext(ctx).Errorf("error upsert employee task  err:%v", err.Error())
	}
	return &pb.Empty{}, err
}

func (em *EmployeeTaskServerImpl) GetEmployeeProgram(ctx context.Context, in *pb.ProgramQuery) (result *pb.Programs, err error) {
	ctx = em.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.Programs{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, em.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := em.trx.GormEndTransaction(ctx); errTrx != nil {
			em.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = em.commonSvc.UpsertCorrelationId(ctx)

	// mapping input
	var empl entity.EmployeeTaskEntity
	if err = em.commonSvc.ObjectMapper(in, &empl); err != nil {
		em.sugar.WithContext(ctx).Errorf("error when mapping object:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// get employee program processing
	programs, err := em.emplTaskSvc.GetEmployeeProgramTask(ctx, empl)
	if err != nil {
		em.sugar.WithContext(ctx).Errorf("error getting program task:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	var tempResult []*pb.Program
	if err = em.commonSvc.ObjectMapper(&programs, &tempResult); err != nil {
		em.sugar.WithContext(ctx).Errorf("error when mapping object:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}
	result = &pb.Programs{Values: tempResult}
	return result, err
}

func (em *EmployeeTaskServerImpl) RenewWellBeingTrigger(ctx context.Context, in *pb.Empty) (result *pb.Empty, err error) {
	ctx = em.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.Empty{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, em.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := em.trx.GormEndTransaction(ctx); errTrx != nil {
			em.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = em.commonSvc.UpsertCorrelationId(ctx)

	if err = em.emplTaskSvc.RenewAllWellBeingTaskSvc(ctx); err != nil {
		em.sugar.WithContext(ctx).Errorf("error when mapping object:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
	}
	return &pb.Empty{}, err
}
