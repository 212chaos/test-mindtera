package http_corporatevalue

import (
	"github.com/gin-gonic/gin"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	corporatevalueservice "github.com/mindtera/corporate-service/service/corporate-value"
	"github.com/mindtera/go-common-module/common/logger"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	commonservice "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateValueHandlerImpl struct {
	sugar                 logger.CustomLogger
	trx                   trx.Transaction
	commonService         commonservice.CommonService
	corporateValueService corporatevalueservice.CorporateValueService
}

func NewCorporateValueHandler(sugar logger.CustomLogger,
	commonService commonservice.CommonService,
	trx trx.Transaction,
	corporateValueService corporatevalueservice.CorporateValueService) CorporateValueHandler {
	return &CorporateValueHandlerImpl{
		sugar:                 sugar,
		trx:                   trx,
		commonService:         commonService,
		corporateValueService: corporateValueService,
	}
}

func (c *CorporateValueHandlerImpl) GetCorporateValueHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	// procees service
	var corporateValues []entity.CorporateValueEntity
	if err := c.corporateValueService.GetCorporateValueService(ctx, &corporateValues); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, commonmodel.ErrResponseType(err.Error()), err.Error())
		return
	}
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateValues)
}

func (c *CorporateValueHandlerImpl) GetCorporateValueBySMMHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonService.SetCorrelationIdFromHeader(ctx)

	var corporateValueModel []model.CorporateValueModel
	if err := c.corporateValueService.GetCorporateValueBySMMService(ctx, &corporateValueModel); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate value service:%v", err.Error())
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonService.CommonErrorResponseSwitcher(ctx, commonmodel.ErrResponseType(err.Error()), err.Error())
		return
	}
	c.commonService.CommonResponseSwitcher(ctx, (commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE), corporateValueModel)
}
