package http_corporate

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"

	corpemplsvc "github.com/mindtera/corporate-service/service/corporate-employee"
	corpsubssvc "github.com/mindtera/corporate-service/service/corporate-subscription"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	eq "github.com/mindtera/quiz-assessment-service/entity"
	qm "github.com/mindtera/quiz-assessment-service/model"
)

type CorporateHandlerImpl struct {
	sugar         logger.CustomLogger
	commonService commonsvc.CommonService
	corpSubsSvc   corpsubssvc.CorporateSubscriptionService
	corpEmplSvc   corpemplsvc.CorporateEmployeeService
	trx           trx.Transaction
}

func NewCorporateHandler(sugar logger.CustomLogger,
	trx trx.Transaction,
	commonService commonsvc.CommonService,
	corpSubsSvc corpsubssvc.CorporateSubscriptionService,
	corpEmplSvc corpemplsvc.CorporateEmployeeService) CorporateHandler {
	return &CorporateHandlerImpl{
		sugar:         sugar,
		trx:           trx,
		commonService: commonService,
		corpSubsSvc:   corpSubsSvc,
		corpEmplSvc:   corpEmplSvc,
	}
}

func (c *CorporateHandlerImpl) GetCorporateStatus(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	tcid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("context is empty")
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}
	cid := tcid.(uuid.UUID)

	// procees service
	var corporateStatus model.CorporateStatus
	c.sugar.WithContext(ctx).Infof("getting status for:%v", cid)
	if err := c.corpEmplSvc.GetCorporateEmployeeStatusService(ctx, cid, &corporateStatus); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate status service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
		return
	}
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateStatus)
}

func (c *CorporateHandlerImpl) GetCorporateDashboardResult(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// query
	var queryForm model.CorpQuery
	if err := ctx.ShouldBindQuery(&queryForm); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx,
			commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, "query is not valid")
		return
	}

	// check corporate_id context
	tcid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("context is empty")
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}
	queryForm.CorporateId = tcid.(uuid.UUID)
	// get quiz assessment result
	var assessmentResult qm.CorporateAssessmentResponseModel
	if errMsg := c.corpEmplSvc.GetCorporateEmployeeDashboardService(ctx, &queryForm, &assessmentResult); errMsg.Error != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate status service:%v", errMsg.Error.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), errMsg.Error)
		c.commonService.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		return
	}

	// handler if data nill
	if len(assessmentResult.Data) <= 0 {
		tempData := make([]eq.CorporateAssessmentSummary, 0)
		assessmentResult.Data = tempData
	}

	c.commonService.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE, assessmentResult)
}

func (c *CorporateHandlerImpl) UpsertCorporateSubscription(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	//binding payloads
	var subscription entity.CorporateSubscriptionEntity
	if err := ctx.ShouldBind(&subscription); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when binding payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}
	// process the payload
	if err := c.corpSubsSvc.UpsertCorporateSubscriptionService(ctx, &subscription); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when processing payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// update active user
	if err := c.corpEmplSvc.UpdateCorpEmployeeByCorpAndSubsId(ctx, subscription, "ACTIVE"); err != nil {
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err)
		return
	}

	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]any{
			"message": "success to add subscription",
			"payload": subscription})
}

// Callback expiration subscription
func (c *CorporateHandlerImpl) CallbackExpiredCorporateSubscription(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// binding request payload
	var schedulerCallback model.CorporateSubscriptionScheduler
	if err := ctx.ShouldBind(&schedulerCallback); err != nil {
		c.sugar.WithContext(ctx).Errorf("error binding callback payload: %v", err)
		err = errors.New("bad request")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err)
		return
	}

	// transform to subscription
	corporateSubs := schedulerCallback.TransformToSubscriptionEntity()
	if err := c.corpSubsSvc.ExpiredCallbackCorporateSubscriptionService(ctx, &corporateSubs); err != nil {
		c.sugar.WithContext(ctx).Errorf("error processing expired callback:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		err = errors.New("error processing payload")
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err)
		return
	}

	// update employee
	if err := c.corpEmplSvc.UpdateCorpEmployeeByCorpAndSubsId(ctx, corporateSubs, "EXPIRED"); err != nil {
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err)
		return
	}

	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), map[string]any{
		"message": "success update to expired",
		"raw":     corporateSubs})
}

// get corporate assessment result status
func (c *CorporateHandlerImpl) GetAssessmentStatus(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// processing payload
	var assessStatus model.AssessmentStatus
	if err := c.corpEmplSvc.GetCorpEmployeeAssessmentStatus(ctx, &assessStatus); err != nil {
		c.sugar.WithContext(ctx).Errorf("error processing assessment status:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		err = errors.New("error processing payload")
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err)
		return
	}
	c.commonService.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), assessStatus)
}

// get corporate well being result status
func (c *CorporateHandlerImpl) GetWellBeingStatus(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// query
	queryForm := struct {
		Department string `form:"department_code"`
		EndTime    int64  `form:"end_date"`
		StartTime  int64  `form:"start_date"`
	}{}
	if err := ctx.ShouldBindQuery(&queryForm); err != nil {
		c.sugar.WithContext(ctx).Errorf("error binding query:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx,
			commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, "query is not valid")
		return
	}

	if queryForm.EndTime <= 0 || queryForm.StartTime <= 0 || (queryForm.EndTime < queryForm.StartTime) {
		err := errors.New("invalid end start date query")
		c.sugar.WithContext(ctx).Errorf("error binding query:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}

	wellBeingStatus := model.WellBeingStatus{
		Department: queryForm.Department,
		StartDate:  time.UnixMilli(queryForm.StartTime),
		EndDate:    time.UnixMilli(queryForm.EndTime),
	}
	if errMsg := c.corpEmplSvc.GetWellBeingAssessmentStatus(ctx, &wellBeingStatus); errMsg.Error != nil {
		c.sugar.WithContext(ctx).Errorf("error processing assessment status:%v", errMsg.Error)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), errMsg.Error)
		c.commonService.CommonErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error)
		return
	}
	c.commonService.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), wellBeingStatus)
}
