package grpcserver_corporatevalue

import (
	"context"
	"errors"

	"github.com/mindtera/corporate-service/model"
	corpvalrel "github.com/mindtera/corporate-service/service/corporate-value-relation"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateValueServerImpl struct {
	sugar         logger.CustomLogger
	trx           trx.Transaction
	grpcServer    conn.GRPCServer
	commonSvc     commonsvc.CommonService
	assert        assert.Assert
	corpValRelSvc corpvalrel.CorporateValueRelationService
	pb.UnimplementedCorporateValueServiceServer
}

func NewCoporateValueServer(sugar logger.CustomLogger,
	grpcServer conn.GRPCServer,
	trx trx.Transaction,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	corpValRelSvc corpvalrel.CorporateValueRelationService) CorporateValueServer {
	sugar.WithContext(context.Background()).Info("initiating corp value grpc server")
	result := &CorporateValueServerImpl{
		grpcServer:    grpcServer,
		sugar:         sugar,
		assert:        assert,
		trx:           trx,
		commonSvc:     commonSvc,
		corpValRelSvc: corpValRelSvc,
	}
	pb.RegisterCorporateValueServiceServer(grpcServer.GetClient(), result)
	return result
}

func (c *CorporateValueServerImpl) GetCorpValueIdByCorpId(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValueIds, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.CorpValueIds{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := c.trx.GormEndTransaction(ctx); errTrx != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = c.commonSvc.UpsertCorrelationId(ctx)

	// checking payload
	if c.assert.IsUUIDEmpty(corpInfo.SubsId) {
		err = errors.New("BAD_REQUEST")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// processing payload
	var corpValIds []model.CorpValueId
	if err = c.corpValRelSvc.GetCorporateValueIds(ctx, corpInfo, &corpValIds); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching data:%v", err.Error())
		err = errors.New("error fetching corporate value data")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// transform value
	var values []*pb.CorpValueId
	if err := c.commonSvc.ObjectMapper(&corpValIds, &values); err != nil {
		c.sugar.WithContext(ctx).Errorf("error mapping data from entity to pb:%v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}
	result = &pb.CorpValueIds{Values: values}
	return result, err
}

func (c *CorporateValueServerImpl) GetCorpValueDetailByCorpId(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValues, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.CorpValues{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := c.trx.GormEndTransaction(ctx); errTrx != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = c.commonSvc.UpsertCorrelationId(ctx)

	// checking payload
	if c.assert.IsUUIDEmpty(corpInfo.SubsId) {
		err = errors.New("BAD_REQUEST")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// processing payload
	var corpValIds []model.CorpValue
	if err = c.corpValRelSvc.GetCorporateValueDetail(ctx, corpInfo, &corpValIds); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching data:%v", err.Error())
		err = errors.New("error fetching corporate value data")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// transform value
	var values []*pb.CorpValue
	if err := c.commonSvc.ObjectMapper(&corpValIds, &values); err != nil {
		c.sugar.WithContext(ctx).Errorf("error mapping data from entity to pb:%v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}
	result = &pb.CorpValues{Values: values}
	return result, err
}

func (c *CorporateValueServerImpl) GetCorpValueName(ctx context.Context, corpInfo *pb.CorpInfo) (result *pb.CorpValueNameArr, err error) {
	ctx = c.commonSvc.GetCorrelationIdFromGrpc(ctx)
	result = &pb.CorpValueNameArr{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTrx := c.trx.GormEndTransaction(ctx); errTrx != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", errTrx)
		}
	}()
	ctx = c.commonSvc.UpsertCorrelationId(ctx)

	// checking payload
	if c.assert.IsUUIDEmpty(corpInfo.SubsId) {
		err = errors.New("BAD_REQUEST")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}

	// get value
	var corpValues []model.CorpValueName
	if err = c.corpValRelSvc.GetCorporateValueName(ctx, corpInfo, &corpValues); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching data:%v", err.Error())
		err = errors.New("error fetching corporate value data")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}
	// transform value
	var values []*pb.CorpValueName
	if err := c.commonSvc.ObjectMapper(&corpValues, &values); err != nil {
		c.sugar.WithContext(ctx).Errorf("error mapping data from entity to pb:%v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return result, err
	}
	result = &pb.CorpValueNameArr{Values: values}
	return result, err
}
