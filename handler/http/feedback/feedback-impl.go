package http_feedback

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mindtera/go-common-module/common/logger"
	mdl "github.com/mindtera/go-common-module/common/v2/model"

	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	feedbacksvc "github.com/mindtera/corporate-service/service/feedback"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type FeedbackHttpHandlerImpl struct {
	sugar       logger.CustomLogger
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	trx         trx.Transaction
	feedbackSvc feedbacksvc.Feedback
}

func NewFeedbackHttpHandler(sugar logger.CustomLogger,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	trx trx.Transaction,
	feedbackSvc feedbacksvc.Feedback) FeedbackHttpHandler {
	return &FeedbackHttpHandlerImpl{
		sugar:       sugar,
		commonSvc:   commonSvc,
		assert:      assert,
		trx:         trx,
		feedbackSvc: feedbackSvc,
	}
}

func (f *FeedbackHttpHandlerImpl) MarkFeedbackAsReadHdl(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(mdl.TRANSACTION_KEY), f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := f.trx.GormEndTransaction(ctx); err != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	f.commonSvc.SetCorrelationIdFromHeader(ctx)

	corpid, exist := ctx.Get("corporate_id")
	if !exist {
		f.sugar.WithContext(ctx).Errorf("unauthorized: corporate id is not found")
		err := errors.New("corporate id is not found")
		f.commonSvc.CommonErrorResponseSwitcher(ctx, mdl.ERR_UNAUTHORIZED_TYPE, err.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), err)
		return
	}
	// getting feedback id
	feedbackId := ctx.Param("feedback_id")
	if f.assert.IsEmpty(feedbackId) {
		f.sugar.WithContext(ctx).Errorf("bad request: feedback id is not found")
		err := errors.New("feedback id is not found")
		f.commonSvc.CommonErrorResponseSwitcher(ctx, mdl.ERR_UNAUTHORIZED_TYPE, err.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), err)
		return
	}
	f.sugar.WithContext(ctx).Infof("getting request to mark as read from %v to id %v", corpid, feedbackId)

	// process payload
	if errMsg := f.feedbackSvc.MarkFeedbackAsReadSvc(ctx, uuid.MustParse(feedbackId)); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("process paylod: error processing payload:%v", errMsg)
		f.commonSvc.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), errMsg.Error)
		return
	}
	f.sugar.WithContext(ctx).Infof("success execute request %v to id %v", corpid, feedbackId)

	f.commonSvc.CommonResponseSwitcher(ctx, mdl.SUCCESS_STANDARD_SUCCESS_TYPE,
		map[string]any{
			"msg":         "success mark as read",
			"feedback_id": feedbackId})
}

func (f *FeedbackHttpHandlerImpl) GetFeedbackHdl(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(mdl.TRANSACTION_KEY), f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := f.trx.GormEndTransaction(ctx); err != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	f.commonSvc.SetCorrelationIdFromHeader(ctx)

	// check corporate first
	corpid, exist := ctx.Get("corporate_id")
	if !exist {
		f.sugar.WithContext(ctx).Errorf("unauthorized: corporate id is not found")
		err := errors.New("corporate id is not found")
		f.commonSvc.CommonErrorResponseSwitcher(ctx, mdl.ERR_UNAUTHORIZED_TYPE, err)
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), err)
		return
	}
	f.sugar.WithContext(ctx).Infof("getting feedback for corporate:%v", corpid)
	// binding header
	var feedbackFilter model.FeedbackFilter
	if err := ctx.ShouldBindQuery(&feedbackFilter); err != nil {
		f.sugar.WithContext(ctx).Errorf("bad request: failed to binding query")
		f.commonSvc.CommonErrorResponseSwitcher(ctx, mdl.ERR_STANDARD_BAD_REQUEST_TYPE, "failed to binding query")
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), err)
		return
	}
	// check time query
	if (feedbackFilter.EndDateInt < feedbackFilter.StartDateInt) ||
		feedbackFilter.StartDateInt == 0 || feedbackFilter.EndDateInt == 0 {
		f.sugar.WithContext(ctx).Errorf("bad request: end start time is invalid")
		err := errors.New("invalid query")
		f.commonSvc.CommonErrorResponseSwitcher(ctx, mdl.ERR_STANDARD_BAD_REQUEST_TYPE, "invalid query")
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), err)
		return
	}
	// binding pagination
	feedbackPaging := mdl.PaginationResponseModel{
		MetaData: mdl.PaginationModel{
			DataPerPage: feedbackFilter.DataPerPage,
			Page:        feedbackFilter.Page,
			PageSize:    feedbackFilter.PageSize,
			TotalData:   feedbackFilter.TotalData,
		},
	}
	feedbackPaging.MetaData.ValidatePaging()

	if errMsg := f.feedbackSvc.GetFeedbackByFilterSvc(ctx, feedbackFilter, &feedbackPaging); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("process paylod: error processing payload:%v", errMsg)
		f.commonSvc.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), errMsg.Error)
		return
	}

	f.commonSvc.CommonResponseSwitcher(ctx, mdl.SUCCESS_STANDARD_SUCCESS_TYPE, feedbackPaging)
}

func (f *FeedbackHttpHandlerImpl) GetFeedbackCategoryHdl(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(mdl.TRANSACTION_KEY), f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := f.trx.GormEndTransaction(ctx); err != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	f.commonSvc.SetCorrelationIdFromHeader(ctx)
	// processing payload
	var feedbacks []entity.FeedbackCategory
	if errMsg := f.feedbackSvc.GetFeedbackCategorySvc(ctx, &feedbacks); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("process paylod: error processing feedback category payload:%v", errMsg)
		f.commonSvc.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), errMsg.Error)
		return
	}

	f.commonSvc.CommonResponseSwitcher(ctx, mdl.SUCCESS_STANDARD_SUCCESS_TYPE, feedbacks)
}

func (f *FeedbackHttpHandlerImpl) GetFeedbackShownCategoryHdl(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(mdl.TRANSACTION_KEY), f.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := f.trx.GormEndTransaction(ctx); err != nil {
			f.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	f.commonSvc.SetCorrelationIdFromHeader(ctx)
	// processing payload
	var feedbacks []entity.FeedbackShownCategory
	if errMsg := f.feedbackSvc.GetFeedbackShownCategorySvc(ctx, &feedbacks); errMsg.Error != nil {
		f.sugar.WithContext(ctx).Errorf("process paylod: error processing feedback category payload:%v", errMsg)
		f.commonSvc.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		ctx.Set(string(mdl.TRANSACTION_ERROR_KEY), errMsg.Error)
		return
	}

	f.commonSvc.CommonResponseSwitcher(ctx, mdl.SUCCESS_STANDARD_SUCCESS_TYPE, feedbacks)
}
