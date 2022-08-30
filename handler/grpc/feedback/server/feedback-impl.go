package grpcserver_feedback

import (
	"context"
	"errors"

	"github.com/mindtera/corporate-service/entity"
	feedbacksvc "github.com/mindtera/corporate-service/service/feedback"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	conn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/server"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type FeedbackServerImpl struct {
	sugar       logger.CustomLogger
	trx         trx.Transaction
	grpcServer  conn.GRPCServer
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	feedbackSvc feedbacksvc.Feedback
	pb.UnimplementedFeedbackServiceServer
}

func NewFeedbackServer(sugar logger.CustomLogger,
	trx trx.Transaction,
	grpcServer conn.GRPCServer,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	feedbackSvc feedbacksvc.Feedback) FeedbackServer {
	result := &FeedbackServerImpl{
		sugar:       sugar,
		trx:         trx,
		grpcServer:  grpcServer,
		commonSvc:   commonSvc,
		assert:      assert,
		feedbackSvc: feedbackSvc,
	}
	pb.RegisterFeedbackServiceServer(grpcServer.GetClient(), result)

	return result
}

func (f *FeedbackServerImpl) UpsertFeedback(ctx context.Context, in *pb.Feedback) (output *pb.Empty, err error) {
	output = &pb.Empty{}
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTx := f.trx.GormEndTransaction(ctx); errTx != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
	}()
	ctx = f.commonSvc.GetCorrelationIdFromGrpc(ctx)

	// transform to feedback entity
	var feedback entity.Feedback
	f.sugar.WithContext(ctx).Infof("transforming payload from: %v", in.GetUserId())
	if err = f.commonSvc.ObjectMapper(in, &feedback); err != nil {
		f.sugar.WithContext(ctx).Errorf("error when transforming payload payload:%v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return output, err
	}

	// check payload is valid or not
	f.sugar.WithContext(ctx).Infof("validate payload from: %v", feedback.UserId)
	if f.assert.IsUUIDEmpty(feedback.CorporateId.String()) ||
		f.assert.IsUUIDEmpty(feedback.UserId.String()) ||
		f.assert.IsEmpty(feedback.Body) ||
		f.assert.IsEmpty(feedback.Sender) ||
		f.assert.IsEmpty(feedback.CategoryCode) ||
		f.assert.IsEmpty(feedback.ShownCategoryCode) {

		f.sugar.WithContext(ctx).Errorf("bad request: INVALID PARAMETER")
		err = errors.New("BAD_REQUEST")
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return output, err
	}

	// process paylaod to insert
	f.sugar.WithContext(ctx).Infof("processing payload from: %v with id :%v", feedback.UserId, feedback.ID)
	if errMsg := f.feedbackSvc.UpsertFeedbackSvc(ctx, &feedback); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("error when processing payload: %v", errMsg)
		err = errMsg.Error
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
	}
	return output, err
}

func (f *FeedbackServerImpl) GetFeedbackCategory(ctx context.Context, in *pb.Empty) (output *pb.FeedbackCategories, err error) {
	output = &pb.FeedbackCategories{}
	ctx = f.commonSvc.GetCorrelationIdFromGrpc(ctx)
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTx := f.trx.GormEndTransaction(ctx); errTx != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
	}()

	// get data from service
	var categories []entity.FeedbackCategory
	if errMsg := f.feedbackSvc.GetFeedbackCategorySvc(ctx, &categories); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("error when processing payload: %v", errMsg)
		err = errMsg.Error
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return output, err
	}

	// transform to grpc payload
	feedbackCategories := struct {
		FeedbackCategories []entity.FeedbackCategory `json:"feedback_categories"`
	}{FeedbackCategories: categories}
	if err = f.commonSvc.ObjectMapper(&feedbackCategories, output); err != nil {
		f.sugar.WithContext(ctx).Errorf("error when transforming payload: %v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
	}
	return output, err
}

func (f *FeedbackServerImpl) GetFeedbackShownCategory(ctx context.Context, in *pb.Empty) (output *pb.FeedbackShownCategories, err error) {
	output = &pb.FeedbackShownCategories{}
	ctx = f.commonSvc.GetCorrelationIdFromGrpc(ctx)
	// set transaction
	ctx = context.WithValue(ctx, commonmodel.TRANSACTION_KEY, f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if errTx := f.trx.GormEndTransaction(ctx); errTx != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
	}()
	ctx = f.commonSvc.UpsertCorrelationId(ctx)
	// get data from service
	var categories []entity.FeedbackShownCategory
	if errMsg := f.feedbackSvc.GetFeedbackShownCategorySvc(ctx, &categories); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("error when processing payload: %v", errMsg)
		err = errMsg.Error
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
		return output, err
	}

	// transform to grpc payload
	feedbackCategories := struct {
		FeedbackShownCategories []entity.FeedbackShownCategory `json:"feedback_shown_categories"`
	}{FeedbackShownCategories: categories}
	if err = f.commonSvc.ObjectMapper(&feedbackCategories, output); err != nil {
		f.sugar.WithContext(ctx).Errorf("error when transforming payload: %v", err)
		ctx = context.WithValue(ctx, commonmodel.TRANSACTION_ERROR_KEY, err)
	}
	return output, err
}
