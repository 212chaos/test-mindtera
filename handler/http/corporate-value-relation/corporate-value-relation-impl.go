package http_corporatevaluerelation

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	corporateemployeeservice "github.com/mindtera/corporate-service/service/corporate-employee"
	corporatevaluerelationservice "github.com/mindtera/corporate-service/service/corporate-value-relation"
	"github.com/mindtera/go-common-module/common/logger"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	commonservice "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateValueRelationHandlerImpl struct {
	sugar                         logger.CustomLogger
	trx                           trx.Transaction
	commonService                 commonservice.CommonService
	corporateValueRelationService corporatevaluerelationservice.CorporateValueRelationService
	corporateEmployeeService      corporateemployeeservice.CorporateEmployeeService
}

func NewCorporateValueRelationHandler(sugar logger.CustomLogger,
	trx trx.Transaction,
	commonService commonservice.CommonService,
	corporateValueRelationService corporatevaluerelationservice.CorporateValueRelationService,
	corporateEmployeeService corporateemployeeservice.CorporateEmployeeService) CorporateValueRelationHandler {
	return &CorporateValueRelationHandlerImpl{
		sugar:                         sugar,
		trx:                           trx,
		commonService:                 commonService,
		corporateValueRelationService: corporateValueRelationService,
		corporateEmployeeService:      corporateEmployeeService,
	}
}

func (c *CorporateValueRelationHandlerImpl) GetCorporateValueRelationHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	cid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("context is empty")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}

	// processing payload
	var corporateRelation model.CorporateValueRelationResultModel
	if err := c.corporateValueRelationService.GetCorporateValueByCorporateID(ctx, cid.(uuid.UUID), &corporateRelation); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value relation service:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// success
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateRelation)
}

func (c *CorporateValueRelationHandlerImpl) UpsertCorporateValueRelationHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	cid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("context is empty")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}

	// get payload
	payload := struct {
		CorporateValue []entity.CorporateValueEntity `json:"corporate_value" binding:"required"`
	}{}
	// binding payload
	if err := ctx.ShouldBind(&payload); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when binding payload:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
	}

	// process payload
	if err := c.corporateValueRelationService.UpsertCorporateValueByCorporateID(ctx, cid.(uuid.UUID), &payload.CorporateValue); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value relation service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// success
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]interface{}{
			"corporate_id":    cid,
			"corporate_value": payload.CorporateValue})
}

func (c *CorporateValueRelationHandlerImpl) CalculateCorporateValueRelationHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	cid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("context is empty")
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}
	corporateID := cid.(uuid.UUID)

	// get payload
	payload := struct {
		CorporateValue []entity.CorporateValueEntity `json:"corporate_value" binding:"required"`
	}{}
	// binding payload
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when binding payload:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
	}

	// make corporate value relation
	var corporateValueRelation []entity.CorporateValueRelation
	var corporateValueID []uuid.UUID

	// sync
	wg, mx := c.commonService.GenerateWaitGroupAndMutex()
	for _, val := range payload.CorporateValue {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, v entity.CorporateValueEntity) {
			// defer function
			defer c.commonService.EndWaitGroupAndMutex(wg_, mx)

			mx.Lock()
			corporateValueRelation = append(corporateValueRelation, entity.CorporateValueRelation{
				CorporateID:      corporateID,
				CorporateValueID: v.ID})
			corporateValueID = append(corporateValueID, v.ID)
		}(&wg, val)
	}
	wg.Wait()

	// process payload
	corporateRelationModel := model.CorporateValueRelationResultModel{RawData: corporateValueRelation}
	if err := c.corporateValueRelationService.CalculateCorporateValueByCorporateRelation(ctx, corporateValueID, &corporateRelationModel); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value relation service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// success
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateRelationModel)
}

func (c *CorporateValueRelationHandlerImpl) DeleteCorporateValueRelationHandler(ctx *gin.Context) {
	// get query
	email := ctx.Query("email")

	// check email
	var isWhitelist bool
	for _, v := range common.DELETION_VALUE_WHITELIST {
		if strings.EqualFold(v, email) {
			isWhitelist = true
			break
		}
	}

	// check condition
	fmt.Println(common.ENV)
	if !(strings.EqualFold(common.ENV, "DEVELOPMENT") || strings.EqualFold(common.ENV, "DEV") || isWhitelist) {
		err := errors.New("request is forbidden for user")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when trying to delete:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}

	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	corporateID := ctx.Query("id")
	if strings.EqualFold(corporateID, "") {
		err := errors.New("query is empty")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}

	if err := c.corporateValueRelationService.DeleteCorporateValueByCorporateID(ctx, uuid.MustParse(corporateID)); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value relation service delete:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}
	// success
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]interface{}{
			"corporate_id": corporateID,
			"message":      "success deleting corporate id relation"})
}

func (c *CorporateValueRelationHandlerImpl) GetCorporateValueRelationPublicHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// get user id public
	uid, isExist := ctx.Get("user_id_public")
	if !isExist {
		err := errors.New("context is empty")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.sugar.WithContext(ctx).Errorf("error when get context:%v", err.Error())
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE), err.Error())
		return
	}

	// check corporate employee first
	corporateEmployee := model.CorporateInformation{
		EmployeeInformation: entity.CorporateEmployeeEntity{
			PublicUserID: uuid.MustParse(uid.(string))}}
	// proceed corporate employee by public service
	if err := c.corporateEmployeeService.GetCorporateEmployeeByPublicIDService(ctx, &corporateEmployee); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate employee service get:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// proceed corporate value relation
	var corporateRelation model.CorporateValueRelationResultModel
	if err := c.corporateValueRelationService.GetCorporateValueByCorporateIDRawOnly(ctx, corporateEmployee.SubscriptionDetail.CorporateID, &corporateRelation); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value relation service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, (commonmodel.ERR_STANDARD_INTERNAL_TYPE), err.Error())
		return
	}

	// success
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateRelation)
}
