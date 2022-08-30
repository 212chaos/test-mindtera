package grpsserver_corporateemployee

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	corpemplsvc "github.com/mindtera/corporate-service/service/corporate-employee"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateEmployeeServerImpl struct {
	sugar       logger.CustomLogger
	trx         trx.Transaction
	grpcServer  conn.GRPCServer
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	corpEmplSvc corpemplsvc.CorporateEmployeeService
	pb.UnimplementedCorporateEmployeeServiceServer
}

func NewCorporateEmployeeServer(grpcServer conn.GRPCServer,
	sugar logger.CustomLogger,
	trx trx.Transaction,
	assert assert.Assert,
	commonSvc commonsvc.CommonService, corpEmplSvc corpemplsvc.CorporateEmployeeService) CorporateEmployeeServer {
	// register proto buff object to grpc server
	sugar.WithContext(context.Background()).Info("initiating corporate employee grpc server")
	result := &CorporateEmployeeServerImpl{
		grpcServer:  grpcServer,
		sugar:       sugar,
		assert:      assert,
		trx:         trx,
		commonSvc:   commonSvc,
		corpEmplSvc: corpEmplSvc,
	}
	pb.RegisterCorporateEmployeeServiceServer(grpcServer.GetClient(), result)
	return result
}

func (c *CorporateEmployeeServerImpl) UpdateCorporatePublicUser(ctx context.Context, publicUser *pb.PublicUser) (result *pb.PublicUser, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.PublicUser{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := c.trx.GormEndTransaction(ctx); errTrx != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = c.commonSvc.UpsertCorrelationId(ctx)

	c.sugar.WithContext(ctx).Infof("update for new registered user for corporate with id:%v", publicUser.PublicId)
	//validate payload
	if err = c.emailValidation(ctx, publicUser); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.PublicUser{}, err
	}
	// call service
	corporateEmployee := entity.CorporateEmployeeEntity{
		Email: publicUser.Email,
	}
	if err = c.corpEmplSvc.GetCorporateEmployeeByEmailService(ctx, &corporateEmployee); err != nil {
		c.sugar.WithContext(ctx).Errorf("cannot fetch corporate employee:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.PublicUser{}, err
	}

	// check if user already exist, no need to update
	if !c.assert.IsUUIDEmpty(corporateEmployee.PublicUserID.String()) ||
		strings.EqualFold(corporateEmployee.RecordFlag, "ACTIVE") {
		c.sugar.WithContext(ctx).Infof("user already registered")
		// mapping result
		result = publicUser
		result.CorporateEmployeeId = corporateEmployee.ID.String()
		return result, nil
	}

	// update corporate entity value
	corporateEmployee.PublicUserID = uuid.MustParse(publicUser.PublicId)
	corporateEmployee.RecordFlag = "ACTIVE"
	corporateEmployee.UpdatedBy = publicUser.Email

	// setting the context
	ctx = context.WithValue(ctx, "corporate_id", corporateEmployee.CorporateID)
	ctx = context.WithValue(ctx, common.PUBLIC_KEY, "public")

	// update corporate employee
	c.sugar.WithContext(ctx).Infof("updating user with corporate employee id:%v", corporateEmployee.ID)
	corporateEmployees := []entity.CorporateEmployeeEntity{corporateEmployee}
	if errRecords, err := c.corpEmplSvc.UpsertCorporateEmployeeService(ctx, &corporateEmployees); err != nil || len(errRecords) > 0 {
		c.sugar.WithContext(ctx).Errorf("error to update user with err:%v err records:%v", err, errRecords)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.PublicUser{}, err
	}

	// mapping result
	c.sugar.WithContext(ctx).Infof("success updating user with corporate employee id:%v", corporateEmployee.ID)
	result = publicUser
	result.CorporateEmployeeId = corporateEmployee.ID.String()
	return result, err
}

func (c *CorporateEmployeeServerImpl) GetCorporatePublicUserByEmail(ctx context.Context, publicUser *pb.PublicUser) (result *pb.CorporateEmployee, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.CorporateEmployee{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	ctx = c.commonSvc.UpsertCorrelationId(ctx)

	c.sugar.WithContext(ctx).Infof("getting corporate employee with public id:%v", publicUser.PublicId)
	//validate payload
	if err = c.emailValidation(ctx, publicUser); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.CorporateEmployee{}, err
	}

	// call service
	corporateEmployee := entity.CorporateEmployeeEntity{
		Email: publicUser.Email,
	}
	if err = c.corpEmplSvc.GetCorporateEmployeeByEmailService(ctx, &corporateEmployee); err != nil {
		c.sugar.WithContext(ctx).Errorf("cannot fetch corporate employee:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.CorporateEmployee{}, err
	}

	if corporateEmployee.RecordFlag == "UNREGISTERED" {
		err = errors.New("user not registered")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.CorporateEmployee{}, err
	}

	// marshalling json try
	b, err := json.Marshal(corporateEmployee)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error marshaling payload json:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.CorporateEmployee{}, err
	}

	// unmarshal and mapping to result
	err = json.Unmarshal(b, &result)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error unmarshaling payload json:%v", err.Error())
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return &pb.CorporateEmployee{}, err
	}

	return result, err
}

func (c *CorporateEmployeeServerImpl) emailValidation(ctx context.Context, publicUser *pb.PublicUser) (err error) {
	// mapping request
	if c.assert.IsEmpty(publicUser.Email) || c.assert.IsUUIDEmpty(publicUser.PublicId) {
		// bad request
		c.sugar.WithContext(ctx).Errorf("cannot proceed payload email %v public id %v", publicUser.Email, publicUser.PublicId)
		err = errors.New("BAD_REQUEST")
		return err
	}
	return err
}
