package http_employeetask

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	empltasksvc "github.com/mindtera/corporate-service/service/employee-task"
	"github.com/mindtera/go-common-module/common/logger"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	trx "github.com/mindtera/go-common-module/common/v2/repository/transaction"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"

	consumermodel "github.com/mindtera/consumer-service/model"
)

type EmployeeTaskHandlerImpl struct {
	sugar       logger.CustomLogger
	commonSvc   commonsvc.CommonService
	assert      assert.Assert
	trx         trx.Transaction
	emplTaskSvc empltasksvc.EmployeeTaskService
}

func NewEmployeeTaskService(sugar logger.CustomLogger,
	trx trx.Transaction,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	emplTaskSvc empltasksvc.EmployeeTaskService,
) EmployeeTaskHandler {
	return &EmployeeTaskHandlerImpl{
		sugar:       sugar,
		assert:      assert,
		trx:         trx,
		commonSvc:   commonSvc,
		emplTaskSvc: emplTaskSvc}
}

// get employee task by query
func (e *EmployeeTaskHandlerImpl) GetEmployeeTask(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	// query employee task needed
	var taskQuery entity.EmployeeTaskEntity
	if err := ctx.ShouldBindQuery(&taskQuery); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding query from url:%v", err)
		err = errors.New("bad query request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}

	qs, exist := ctx.GetQuery("subscription_id")
	if exist {
		taskQuery.SubscriptionID = uuid.MustParse(qs)
	}

	// query pagination needed
	// fetch query
	var pagination commonmodel.PaginationModel
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when binding query payload:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	pagination.ValidatePaging()
	employeeTask := commonmodel.PaginationResponseModel{MetaData: pagination}

	// validation whether search is exist or not
	q, isExist := ctx.GetQuery("q")

	// processing payload
	switch {
	case isExist:
		if err := e.emplTaskSvc.
			SearchEmployeeTaskByUserIDSubscriptionService(ctx, strings.ToLower(q),
				&taskQuery, &employeeTask); err != nil {
			// response
			ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
			e.commonSvc.CommonErrorResponseSwitcher(ctx,
				(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
				err.Error())
			return
		}
	default:
		if err := e.emplTaskSvc.
			GetEmployeeTaskByUserIDSubscriptionService(ctx, &taskQuery, &employeeTask); err != nil {

			// response
			ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
			e.commonSvc.CommonErrorResponseSwitcher(ctx,
				(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
				err.Error())
			return
		}
	}
	// success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		employeeTask)
}

// upsert employee task in batch
func (e *EmployeeTaskHandlerImpl) UpsertEmployeeTask(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var employeeTasks model.EmployeeTaskRequest
	if err := ctx.ShouldBind(&employeeTasks); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	//processing payload
	employeeTaskArr := employeeTasks.EmployeeTasks
	if err := e.emplTaskSvc.UpsertEmployeeInBatchService(ctx, &employeeTaskArr); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]any{
			"message": "success upsert task",
			"raw":     employeeTaskArr})
}

// notifier employee task in batch
func (e *EmployeeTaskHandlerImpl) NotifyEmployeeTask(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	// binding payload
	var task entity.EmployeeTaskEntity
	if err := ctx.ShouldBind(&task); err != nil {
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
	}

	// checking payload
	if e.assert.IsEmpty(string(task.TaskType)) {
		err := errors.New("bad task type payload")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
	}

	// process service
	if err := e.emplTaskSvc.TaskReminderNotifier(ctx, task); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}

	// response
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		"successfully sending email")
}

// callback employee task from scheduler
func (e *EmployeeTaskHandlerImpl) CallbackEmployeeTask(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	// binding request payload
	var schedulerCallback model.CorporateSubscriptionScheduler
	if err := ctx.ShouldBind(&schedulerCallback); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding callback payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	ctx.Set("corporate_id", schedulerCallback.CorporateID)
	// processing payload
	template := schedulerCallback.TransformToEmployeeTask()
	if err := e.emplTaskSvc.CallbackEmployeeTaskForAssessment(ctx, template); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]any{
			"message": "success update task",
			"raw":     template})
}

// upsert employee task in batch
func (e *EmployeeTaskHandlerImpl) InsertEmployeeTaskProgramAssignment(ctx *gin.Context) {

	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var paRequest model.ProgramAssignmentRequest
	if err := ctx.ShouldBind(&paRequest); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}

	paCreateData, err := paRequest.GenerateCreateModel()
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}

	//processing payload
	if err := e.emplTaskSvc.InsertEmployeeTaskProgramAssignment(ctx,
		paCreateData.AssigningMode,
		paCreateData.RelatedID,
		paCreateData.AssignDate,
		paCreateData.ProgramID,
		paCreateData.ConfirmIgnoreTakenProgram,
	); err != nil {

		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		if err.Error() == string(commonmodel.ERR_PROGRAMS_HAVE_BEEN_TAKEN) {
			e.commonSvc.CommonErrorResponseSwitcher(
				ctx,
				commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE,
				string(commonmodel.ERR_PROGRAMS_HAVE_BEEN_TAKEN),
			)
			return
		} else if err.Error() == string(commonmodel.ERR_NO_EMPLOYEE_TASK_ADDED) {
			e.commonSvc.CommonErrorResponseSwitcher(
				ctx,
				commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE,
				string(commonmodel.ERR_NO_EMPLOYEE_TASK_ADDED),
			)
			return
		}
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]any{
			"message": "success upsert task",
		})
}

func (e *EmployeeTaskHandlerImpl) GetEmployeeTaskProgramAssignment(ctx *gin.Context) {

	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var paginationQuery model.AssignedProgramRequestForm
	if err := ctx.ShouldBindQuery(&paginationQuery); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	paginationQuery.ValidateRequest()
	paging := commonmodel.PaginationResponseModel{
		MetaData: paginationQuery.PaginationModel,
	}
	//processing payload
	if err := e.emplTaskSvc.GetProgramAssignmentHistory(ctx,
		&paging,
	); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		paging)
}

func (e *EmployeeTaskHandlerImpl) SearchProgram(ctx *gin.Context) {
	var err error
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var searchRequest model.SearchProgramRequest
	if err := ctx.ShouldBindQuery(&searchRequest); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding search program request payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	searchQuery, ok := searchRequest.GenerateModel()
	if !ok {

		e.sugar.WithContext(ctx).Errorf("error generating serach query request")
		err := errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return

	}

	var programs []consumermodel.MProgram
	if searchQuery == "" {
		programs, err = e.emplTaskSvc.GetPrograms(ctx)
	} else {
		programs, err = e.emplTaskSvc.SearchPrograms(ctx, searchQuery)
	}
	//processing payload
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}

	totalData := len(programs)
	paging := commonmodel.PaginationResponseModel{
		MetaData: commonmodel.PaginationModel{
			Page:        1,
			PageSize:    5,
			DataPerPage: totalData,
			TotalData:   0,
		},
		RawData: programs,
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		paging)
}

func (e *EmployeeTaskHandlerImpl) GetCorporateServices(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var paginationQuery model.MCorporateServiceRequest
	if err := ctx.ShouldBindQuery(&paginationQuery); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	paginationQuery.ValidatePaging()
	//processing payload
	outputArr, totalData, err := e.emplTaskSvc.GetCorporateServices(ctx, paginationQuery)
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	paging := commonmodel.PaginationResponseModel{
		MetaData: commonmodel.PaginationModel{
			Page:        paginationQuery.Page,
			PageSize:    paginationQuery.PageSize,
			TotalData:   totalData,
			DataPerPage: len(outputArr),
		},
		RawData: outputArr,
	}
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		paging)
}

func (e *EmployeeTaskHandlerImpl) GetCorporateServiceByID(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	corporateServiceIDStr := ctx.Param("corporate_service_id")
	_, err := uuid.Parse(corporateServiceIDStr)
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//processing payload
	result, err := e.emplTaskSvc.GetCorporateServiceByID(ctx, corporateServiceIDStr)
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		result)
}

// upsert employee task in batch
func (e *EmployeeTaskHandlerImpl) InsertCorporateServiceBooking(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var bookingJson model.CorporateServiceBookingJson
	if err := ctx.ShouldBind(&bookingJson); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	bookingCreationData, err := bookingJson.GenerateCreateModel()
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	//processing payload
	if err := e.emplTaskSvc.InsertCorporateServiceBooking(ctx, bookingCreationData.ServiceID, bookingCreationData.ScheduledAt); err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		map[string]any{
			"message": "success insert booking",
		})
}

func (e *EmployeeTaskHandlerImpl) GetCurrentCorporateServiceBooking(ctx *gin.Context) {

	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var paginationQuery commonmodel.PaginationModel
	if err := ctx.ShouldBindQuery(&paginationQuery); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding employee task payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}
	paginationQuery.ValidatePaging()
	//processing payload

	output, totalData, err := e.emplTaskSvc.GetCorporateServiceBooking(ctx, paginationQuery.PageSize, paginationQuery.GetOffset())
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	paging := commonmodel.PaginationResponseModel{
		MetaData: commonmodel.PaginationModel{
			Page:        paginationQuery.Page,
			PageSize:    paginationQuery.PageSize,
			TotalData:   totalData,
			DataPerPage: len(output),
		},
		RawData: output,
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		paging)
}

func (e *EmployeeTaskHandlerImpl) SearchCorporateService(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	//binding payload
	var searchRequest model.SearchCorporateServiceRequest
	if err := ctx.ShouldBindQuery(&searchRequest); err != nil {
		e.sugar.WithContext(ctx).Errorf("error binding search program request payload: %v", err)
		err = errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return
	}

	ok := searchRequest.ValidateRequest()
	if !ok {
		e.sugar.WithContext(ctx).Errorf("error generating serach query request")
		err := errors.New("bad request")
		// response
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_BAD_REQUEST_TYPE),
			err.Error())
		return

	}
	//processing payload
	corporateServices, err := e.emplTaskSvc.SearchCorporateService(ctx, searchRequest.Query)
	if err != nil {
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}

	paging := commonmodel.PaginationResponseModel{
		// MetaData: paginationQuery,
		RawData: corporateServices,
	}
	//success
	e.commonSvc.CommonResponseSwitcher(ctx,
		(commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE),
		paging)
}

func (e *EmployeeTaskHandlerImpl) RenewWellBeingInternal(ctx *gin.Context) {
	// set transaction
	ctx.Set(string(commonmodel.TRANSACTION_KEY), e.trx.GormBeginTransaction(ctx))
	defer func() {
		// close transaction
		if err := e.trx.GormEndTransaction(ctx); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when process payload:%v", err)
		}
	}()
	e.commonSvc.SetCorrelationIdFromHeader(ctx)

	if err := e.emplTaskSvc.RenewAllWellBeingTaskSvc(ctx); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when processing service:%v", err)
		ctx.Set(string(commonmodel.TRANSACTION_ERROR_KEY), err)
		e.commonSvc.CommonErrorResponseSwitcher(ctx,
			(commonmodel.ERR_STANDARD_INTERNAL_TYPE),
			err.Error())
		return
	}
	e.commonSvc.CommonResponseSwitcher(ctx, commonmodel.SUCCESS_STANDARD_SUCCESS_TYPE, "success renew well begin")
}
