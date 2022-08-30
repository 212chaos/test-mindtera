package http_corporateemployee

import (
	"errors"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	corpemplsvc "github.com/mindtera/corporate-service/service/corporate-employee"
	"github.com/mindtera/go-common-module/common/logger"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type CorporateEmployeeHandlerImpl struct {
	sugar       logger.CustomLogger
	commonSvc   commonsvc.CommonService
	trx         trx.Transaction
	assert      assert.Assert
	corpEmplSvc corpemplsvc.CorporateEmployeeService
}

func NewCorporateEmployeeHandler(
	sugar logger.CustomLogger,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	trx trx.Transaction,
	corpEmplSvc corpemplsvc.CorporateEmployeeService) CorporateEmployeeHandler {
	return &CorporateEmployeeHandlerImpl{
		sugar:       sugar,
		assert:      assert,
		trx:         trx,
		commonSvc:   commonSvc,
		corpEmplSvc: corpEmplSvc,
	}
}

func (c *CorporateEmployeeHandlerImpl) GetCorporateEmployeeHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonSvc.SetCorrelationIdFromHeader(ctx)

	// get corporate id first
	cid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("corporate id is not found")
		c.sugar.WithContext(ctx).Errorf("error when getting id payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}
	// get corporate code
	ccode, exist := ctx.Get("corporate_code")
	if !exist {
		err := errors.New("corporate code is not found")
		c.sugar.WithContext(ctx).Errorf("error when getting id payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}

	c.sugar.WithContext(ctx).Infof("getting employee for:%v, %v", cid, ccode)

	//fetch query
	var paginationQuery commonmodel.PaginationModel
	if err := ctx.ShouldBindQuery(&paginationQuery); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when binding query payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}
	paginationQuery.ValidatePaging()

	// process payload
	employeePaging := commonmodel.PaginationResponseModel{MetaData: paginationQuery}

	// check employee query
	query := strings.ToLower(ctx.Query("q"))
	if c.assert.IsEmpty(query) {
		if err := c.corpEmplSvc.GetCorporateEmployeeService(ctx, &employeePaging); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
			ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
			c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
			return
		}
	} else {
		if err := c.corpEmplSvc.GetCorporateEmployeeByNameService(ctx, query, &employeePaging); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
			ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
			c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
			return
		}
	}

	// final response
	c.commonSvc.CommonResponseSwitcher(ctx,
		commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
		employeePaging)
}

func (c *CorporateEmployeeHandlerImpl) GetCorporateEmployeePublicIDHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonSvc.SetCorrelationIdFromHeader(ctx)

	userPublicID := ctx.GetString("user_id_public")
	if c.assert.IsEmpty(userPublicID) {
		err := errors.New("user id is not found")
		c.sugar.WithContext(ctx).Errorf("user id is not found from context:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}

	// proceed the payload
	var corporateInfo model.CorporateInformation
	corporateInfo.EmployeeInformation.PublicUserID = uuid.MustParse(userPublicID)

	if err := c.corpEmplSvc.GetCorporateEmployeeByPublicIDService(ctx, &corporateInfo); err != nil {
		err := errors.New("error processing payload")
		c.sugar.WithContext(ctx).Errorf("error processing payload service:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
		return
	}

	// success response
	c.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
		corporateInfo)
}

func (c *CorporateEmployeeHandlerImpl) UpsertCorporateEmployeeHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonSvc.SetCorrelationIdFromHeader(ctx)

	// get corporate id first
	cid, exist := ctx.Get("corporate_id")
	if !exist {
		err := errors.New("corporate id is not found")
		c.sugar.WithContext(ctx).Errorf("error when getting id payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}
	// get corporate code
	ccode, exist := ctx.Get("corporate_code")
	if !exist {
		err := errors.New("corporate code is not found")
		c.sugar.WithContext(ctx).Errorf("error when getting id payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}
	c.sugar.WithContext(ctx).Infof("upsert employee for:%v, %v", cid, ccode)

	// binding
	employeeReq := struct {
		EmployeeRecord []entity.CorporateEmployeeEntity `json:"employee_record"`
	}{}
	if err := ctx.ShouldBind(&employeeReq); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when binding payload")
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}

	// process
	errRecord, err := c.corpEmplSvc.UpsertCorporateEmployeeService(ctx, &employeeReq.EmployeeRecord)
	if errRecord != nil {
		sort.Slice(errRecord, func(i, j int) bool {
			return errRecord[i].Number < errRecord[j].Number
		})
	}

	if err != nil {
		if err.Error() == "BAD_REQUEST" {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
			ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
			c.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
				map[string]any{
					"type":         "PARTIAL_SUCCESS",
					"error_record": errRecord})
			return
		}
		c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
		return
	}

	// final response
	if errRecord != nil {
		c.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
			map[string]any{
				"type":         "PARTIAL_SUCCESS",
				"error_record": errRecord})
		return
	}
	c.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
		map[string]any{"type": "FULLY_SUCCESS"})
}

func (c *CorporateEmployeeHandlerImpl) DeleteCorporateEmployeeHandler(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), c.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := c.trx.GormEndTransaction(ctx); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	c.commonSvc.SetCorrelationIdFromHeader(ctx)

	// get corporate id first
	cid, exist := ctx.GetQuery("id")
	if !exist {
		err := errors.New("ID is not found")
		c.sugar.WithContext(ctx).Errorf("error when getting id payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE, err.Error())
		return
	}

	employee := entity.CorporateEmployeeEntity{
		ID: uuid.MustParse(cid)}

	if err := c.corpEmplSvc.DeleteCorporateEmployeeService(ctx, &employee); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		c.commonSvc.CommonErrorResponseSwitcher(ctx, commonmodel.ERR_STANDARD_INTERNAL_TYPE, err.Error())
		return
	}

	c.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE,
		map[string]any{
			"message": "success delete employee",
			"payload": employee})
}
