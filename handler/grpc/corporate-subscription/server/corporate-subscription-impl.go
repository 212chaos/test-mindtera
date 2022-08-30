package grpcserver_corporatesubscription

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/model"
	corpemplsvc "github.com/mindtera/corporate-service/service/corporate-employee"
	"github.com/mindtera/go-common-module/common/logger"
	"github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateSubscriptionSrvImpl struct {
	sugar       logger.CustomLogger
	trx         trx.Transaction
	grpcServer  conn.GRPCServer
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	corpEmplSvc corpemplsvc.CorporateEmployeeService
	pb.UnimplementedCorporateSubscriptionServiceServer
}

func NewCorporateSubscriptionSrv(sugar logger.CustomLogger,
	trx trx.Transaction,
	grpcServer conn.GRPCServer,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	corpEmplSvc corpemplsvc.CorporateEmployeeService) CorporateSubscriptionSrv {

	result := &CorporateSubscriptionSrvImpl{
		sugar:       sugar,
		trx:         trx,
		commonSvc:   commonSvc,
		assert:      assert,
		corpEmplSvc: corpEmplSvc,
		grpcServer:  grpcServer}
	pb.RegisterCorporateSubscriptionServiceServer(grpcServer.GetClient(), result)
	return result
}

func (c *CorporateSubscriptionSrvImpl) GetCorporateSubscription(ctx context.Context, corp *pb.Corporate) (result *pb.Empty, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.Empty{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := c.trx.GormEndTransaction(ctx); errTrx != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()

	// binding payload
	var corpId uuid.UUID
	if err = c.commonSvc.ObjectMapper(corp.GetId(), &corpId); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		c.sugar.WithContext(ctx).Errorf("error when mapping grpc payload to uuid:%v", err.Error())
		return result, err
	}
	// check uuid valid or not
	if c.assert.IsUUIDEmpty(corpId.String()) {
		err = errors.New("uuid is invalid")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		c.sugar.WithContext(ctx).Errorf("error uuid:%v", err.Error())
		return result, err
	}
	// get subscription
	var status model.CorporateStatus
	if err = c.corpEmplSvc.GetCorporateEmployeeStatusService(ctx, corpId, &status); err != nil {
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		c.sugar.WithContext(ctx).Errorf("error processing payload to get status:%v", err.Error())
		return result, err
	}
	//check subscription
	if c.assert.IsUUIDEmpty(status.SubscriptionDetail.ID.String()) {
		err = errors.New("subscription uuid is invalid")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		c.sugar.WithContext(ctx).Errorf("error subscription uuid is invalid:%v", err.Error())
	}
	return result, err
}
