package service_employeetask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	consumermodel "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"

	corpprogramrepo "github.com/mindtera/corporate-service/repository/corporate-program"
	repo "github.com/mindtera/corporate-service/repository/employee-task"
	corpsubssvc "github.com/mindtera/corporate-service/service/corporate-subscription"
	mailjetclient "github.com/mindtera/go-common-module/common/client/mailjet-client"
	"github.com/mindtera/go-common-module/common/logger"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
)

type EmployeeTaskServiceImpl struct {
	sugar            logger.CustomLogger
	repository       repo.EmployeeTaskRepository
	emplProgramRepo  corpprogramrepo.CorporateProgram
	redisSvc         redissvc.RedisSvc
	commonSvc        commonsvc.CommonService
	assert           assert.Assert
	corporateSubsSvc corpsubssvc.CorporateSubscriptionService
	mailjet          mailjetclient.MailjetClient
}

var (
	key = "EMPLOYEE_TASK"
)

//constructor service to create new employee task service DI
func NewEmployeeTaskService(sugar logger.CustomLogger,
	redisSvc redissvc.RedisSvc,
	repository repo.EmployeeTaskRepository,
	emplProgramRepo corpprogramrepo.CorporateProgram,
	commonSvc commonsvc.CommonService,
	assert assert.Assert,
	corporateSubsSvc corpsubssvc.CorporateSubscriptionService,
	mailjet mailjetclient.MailjetClient) EmployeeTaskService {
	return &EmployeeTaskServiceImpl{
		sugar:            sugar,
		repository:       repository,
		emplProgramRepo:  emplProgramRepo,
		redisSvc:         redisSvc,
		commonSvc:        commonSvc,
		assert:           assert,
		corporateSubsSvc: corporateSubsSvc,
		mailjet:          mailjet,
	}
}

// get employee task based on query needed
func (e *EmployeeTaskServiceImpl) GetEmployeeTaskByUserIDSubscriptionService(ctx context.Context, empTask *entity.EmployeeTaskEntity, empTaskPagination *commonmodel.PaginationResponseModel) (err error) {
	// subscription query
	if e.assert.IsUUIDEmpty(empTask.SubscriptionID.String()) {
		subsID, err := e.getSubscriptionIDFromContext(ctx)
		if err != nil {
			return err
		}
		empTask.SubscriptionID = subsID
	}

	// pagination payload
	paginationPayload := empTaskPagination.MetaData

	// check redis
	// keyImpl := fmt.Sprintf("%v_%v_%v_%v_%v", key,
	// 	empTask.SubscriptionID, empTask.TaskType,
	// 	paginationPayload.PageSize, paginationPayload.Page)

	// if !e.assert.IsEmpty(empTask.Email) {
	// 	keyImpl = fmt.Sprintf("%v_%v_%v_%v_%v_%v", key,
	// 		empTask.SubscriptionID, empTask.Email, empTask.TaskType,
	// 		paginationPayload.PageSize, paginationPayload.Page)
	// }

	// e.sugar.WithContext(ctx).Infof("getting employee task with key:%v", keyImpl)
	// if err = e.redisSvc.Get(ctx, keyImpl, empTaskPagination); err != nil {
	e.sugar.WithContext(ctx).Infof("getting employee task from database")
	var empTasks []entity.EmployeeTaskEntity
	if err = e.repository.GetEmployeeTaskByUserIDSubscription(ctx, empTask, &empTasks, &paginationPayload); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when fetching employee task:%v", err)
		err = errors.New("error processing payload")
		return err
	}
	empTaskPagination.MetaData = paginationPayload
	empTaskPagination.RawData = empTasks
	e.sugar.WithContext(ctx).Info("successfully fetch employee task")
	// }

	// go func() {
	// 	e.sugar.WithContext(ctx).Info("set key value in redis:%v", keyImpl)
	// 	if err = e.redisSvc.Set(context.Background(), keyImpl, empTaskPagination, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
	// 		e.sugar.WithContext(ctx).Errorf("error when set key value in redis:%v", err)
	// 		return
	// 	}
	// 	e.sugar.WithContext(ctx).Infof("get item for key :%v amount:%v", keyImpl, empTaskPagination.MetaData.DataPerPage)
	// }()
	return err
}

func (e *EmployeeTaskServiceImpl) SearchEmployeeTaskByUserIDSubscriptionService(ctx context.Context, search string, empTask *entity.EmployeeTaskEntity, empTaskPagination *commonmodel.PaginationResponseModel) (err error) {

	// subscription query
	if e.assert.IsUUIDEmpty(empTask.SubscriptionID.String()) {
		subsID, err := e.getSubscriptionIDFromContext(ctx)
		if err != nil {
			return err
		}
		empTask.SubscriptionID = subsID
	}

	// pagination payload
	paginationPayload := empTaskPagination.MetaData

	// process payload needed for query
	search = "%" + search + "%"
	empTask.EmployeeName = search
	empTask.Department = search
	empTask.Email = search

	// check redis
	keyImpl := fmt.Sprintf("%v_%v_%v_%v_%v_%v", key,
		empTask.SubscriptionID, search,
		empTask.TaskType,
		paginationPayload.PageSize, paginationPayload.Page)

	e.sugar.WithContext(ctx).Infof("getting employee task with key search:%v", keyImpl)
	if err = e.redisSvc.Get(ctx, keyImpl, empTaskPagination); err != nil {
		e.sugar.WithContext(ctx).Infof("getting employee tack from database")
		var empTasks []entity.EmployeeTaskEntity
		if err = e.repository.SearchEmployeeTaskByUserIDSubscription(ctx, empTask, &empTasks, &paginationPayload); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when fetching employee task:%v", err)
			err = errors.New("error processing payload")
			return err
		}
		empTaskPagination.MetaData = paginationPayload
		empTaskPagination.RawData = empTasks
		e.sugar.WithContext(ctx).Info("successfully fetch employee task")
	}

	go func() {
		e.sugar.WithContext(ctx).Info("set key value in redis:%v", keyImpl)
		if err = e.redisSvc.Set(context.Background(), keyImpl, empTaskPagination, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when set key value in redis:%v", err)
			return
		}
		e.sugar.WithContext(ctx).Infof("get item for key :%v amount:%v", keyImpl, empTaskPagination.MetaData.DataPerPage)
	}()
	return err
}

// Get employee task by user id and task type
func (e *EmployeeTaskServiceImpl) GetEmployeeTaskByUserIdAndType(ctx context.Context, empTask *model.EmployeeTaskQuery, empTasks *[]entity.EmployeeTaskEntity) (err error) {
	// search multiple task
	e.sugar.WithContext(ctx).Infof("searching task for user id:%v with types:%v", empTask.UserId, empTask.Types)
	if err = e.repository.GetEmployeeTaskByUserIDAndTypes(ctx, *empTask, empTasks); err != nil {
		return err
	}
	if len(*empTasks) <= 0 {
		err = errors.New("employee not assigned to searched task")
	}
	return err
}

// upsert employee in batch
func (e *EmployeeTaskServiceImpl) UpsertEmployeeInBatchService(ctx context.Context, empTasks *[]entity.EmployeeTaskEntity) (err error) {

	// get subsid for corporate
	var subsID uuid.UUID
	if len(*empTasks) > 0 {
		if !e.assert.IsUUIDEmpty((*empTasks)[0].SubscriptionID.String()) {
			subsID = (*empTasks)[0].SubscriptionID
		}
	}

	// delete all cache first (not recommended but no choice for now)
	chanErr := make(chan error)
	go func(ce chan error) {
		e.sugar.WithContext(ctx).Info("deleting employee task cache")
		if err := e.redisSvc.Delete(context.Background(), fmt.Sprintf("%v_%v_*", key, subsID)); err != nil {
			e.sugar.WithContext(ctx).Errorf("error deleting cache from redis:%v", err)
			err = errors.New("error processing payload")
			ce <- err
		}
		ce <- nil
	}(chanErr)

	// checking payload for employee tasks should contains user id and task type
	wg, mx := e.commonSvc.GenerateWaitGroupAndMutex()
	chanErrTask := make(chan error, len(*empTasks))
	for _, val := range *empTasks {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, v entity.EmployeeTaskEntity) {
			defer e.commonSvc.EndWaitGroupAndMutex(wg_, mx)
			mx.Lock()
			chanErrTask <- e.validateEmployeeTaskPayload(ctx, v)
		}(&wg, val)
	}
	wg.Wait()

	// check all errors from channel
	for i := 0; i < len(*empTasks); i++ {
		if err := <-chanErrTask; err != nil {
			e.sugar.WithContext(ctx).Errorf("error when validating the payload: uuid or task empty")
			return err
		}
	}
	// check deleting cache error
	if err = <-chanErr; err != nil {
		e.sugar.WithContext(ctx).Errorf("error deleting cache in redis:%v", err)
		err = errors.New("error processing payload")
		return err
	}

	//upsert batch payload
	if err = e.repository.UpsertBatchEmployeeTaskRepository(ctx, empTasks); err != nil {
		e.sugar.WithContext(ctx).Errorf("error upserting batch payload:%v", err)
		err = errors.New("error store payload")
		return err
	}
	return err
}

// upsert employee from corporate employee entity
func (e *EmployeeTaskServiceImpl) UpsertEmployeeTaskByEmployeeEntityService(ctx context.Context, employeeEntities *[]entity.CorporateEmployeeEntity, templateTask entity.EmployeeTaskEntity) (err error) {

	// generate payload for employee task
	wg, mx := e.commonSvc.GenerateWaitGroupAndMutex()
	// loop all employee payload
	var empTasks []entity.EmployeeTaskEntity
	for _, val := range *employeeEntities {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, v entity.CorporateEmployeeEntity) {
			defer e.commonSvc.EndWaitGroupAndMutex(wg_, mx)
			mx.Lock()

			// append payload
			if !e.assert.IsUUIDEmpty(v.PublicUserID.String()) {
				templateTask.UserID = v.PublicUserID
				templateTask.Email = v.Email
				empTasks = append(empTasks, templateTask)
			}
		}(&wg, val)
	}
	wg.Wait()

	if len(empTasks) > 0 {
		// store generated employee task array
		if err = e.UpsertEmployeeInBatchService(ctx, &empTasks); err != nil {
			e.sugar.WithContext(ctx).Errorf("error upserting corporate employee task")
			err = errors.New("error storing payload")
			return err
		}
	}
	return err
}

// callback employee task
func (e *EmployeeTaskServiceImpl) CallbackEmployeeTaskForAssessment(ctx context.Context, templateTask entity.EmployeeTaskEntity) (err error) {

	// generate query for getting assigned employee
	newType := templateTask.TaskType
	additionalTimeDur := time.Now().Sub(time.Now().AddDate(0, -1, 0))

	switch newType {
	case "POST_ASSESSMENT":
		templateTask.TaskType = "PRE_ASSESSMENT"
		additionalTimeDur *= 2
	case "FOLLOWUP_ASSESSMENT":
		templateTask.TaskType = "POST_ASSESSMENT"
	}

	// check if development environment
	if strings.EqualFold(common.ENV, "DEVELOPMENT") ||
		strings.EqualFold(common.ENV, "DEV") {
		additionalTimeDur = 5 * time.Minute
	}

	// getting employees first
	paginationModel := commonmodel.PaginationResponseModel{
		MetaData: commonmodel.PaginationModel{
			Page:     -1,
			PageSize: -1}}
	if err = e.GetEmployeeTaskByUserIDSubscriptionService(ctx, &templateTask, &paginationModel); err != nil {
		e.sugar.WithContext(ctx).Errorf("error fetching get employee assigned task:%v", err)
		err = errors.New("error fetching payload")
		return err
	}

	// transform employee task
	var empTask []entity.EmployeeTaskEntity
	b, err := json.Marshal(paginationModel.RawData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &empTask)
	if err != nil {
		return err
	}

	// loop over employee task
	if len(empTask) > 0 {
		e.sugar.WithContext(ctx).Info("generating new task assessment in batch for subs id:%v", templateTask.SubscriptionID)
		wg, mx := e.commonSvc.GenerateWaitGroupAndMutex()
		for idx := range empTask {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, i int) {
				defer e.commonSvc.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				endTime := empTask[i].DueDate.Add(additionalTimeDur)
				empTask[i].AssignDate = empTask[i].DueDate
				empTask[i].ID = uuid.New()
				empTask[i].TaskType = newType
				empTask[i].TaskStatus = "ACTIVE"
				empTask[i].DueDate = &endTime
			}(&wg, idx)
		}
		wg.Wait()

		// upsert batch
		e.sugar.WithContext(ctx).Info("upserting new task assessment in batch for subs id:%v", templateTask.SubscriptionID)
		if err = e.UpsertEmployeeInBatchService(ctx, &empTask); err != nil {
			e.sugar.WithContext(ctx).Errorf("error upserting for subs id:%v", templateTask.SubscriptionID)
			return err
		}
	}

	return err
}

// notifier for task reminder
func (e *EmployeeTaskServiceImpl) TaskReminderNotifier(ctx context.Context, task entity.EmployeeTaskEntity) (err error) {

	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err = errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)

	// get subscription
	e.sugar.WithContext(ctx).Infof("getting subscription for corporate:%v type:%v", corporateID, task.TaskType)
	corporateSubscription := entity.CorporateSubscriptionEntity{CorporateID: corporateID}
	if e.corporateSubsSvc.GetCorporateSubscriptionService(ctx, &corporateSubscription); err != nil {
		e.sugar.WithContext(ctx).Errorf("error get corporate subs:%v", err)
		return errors.New("internal server error to get corporate detail")
	}

	//check subscription id
	if e.assert.IsUUIDEmpty(corporateSubscription.ID.String()) {
		err = errors.New("corporate does not have subscription")
		return err
	}
	task.SubscriptionID = corporateSubscription.ID

	// getting all employee for subscription id
	e.sugar.WithContext(ctx).Infof("getting employees for subsid:%v", corporateSubscription.ID)
	var employees []entity.EmployeeTaskEntity
	if err = e.repository.GetTaskBySubscriptionIdAndType(ctx, &task, &employees); err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}

	//sending email
	go e.validateTaskReminderEmail(ctx, employees)

	return err
}

// get employee program tasks
func (e *EmployeeTaskServiceImpl) GetEmployeeProgramTask(ctx context.Context, employee entity.EmployeeTaskEntity) (programs []model.EmployeeProgram, err error) {
	e.sugar.WithContext(ctx).Infof("getting program for employee id:%v", employee.UserID)
	// getting from repo
	programs, err = e.emplProgramRepo.GetCorporateEmployeeProgram(ctx, employee)
	if len(programs) <= 0 {
		err = errors.New("user is not assigned to any programs")
	}
	return programs, err
}

// EmployeeTaskAssignProgram
// relatedID could be:
// corporate_id, user_id, department_code
func (e *EmployeeTaskServiceImpl) InsertEmployeeTaskProgramAssignment(ctx context.Context, assigningMode entity.TypeAssigningMode, relatedID string, assignDate time.Time, programIDArr []string, confirmIgnoreTakenProgram bool) (err error) {
	tempCorporateID := ctx.Value("corporate_id")
	if tempCorporateID == nil {
		e.sugar.WithContext(ctx).Errorf("error retrieveing corporate_id")
		return errors.New("internal server error to get employee task")

	}

	corporateID := tempCorporateID.(uuid.UUID)
	historyRelatedID := ""

	var employees []entity.CorporateEmployeeEntityWithDepartment
	if assigningMode == entity.ASSIGN_MODE_CORPORATE {
		employees, err = e.repository.GetEmployeeByCorporateID(ctx, corporateID.String())
		if err != nil {
			e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
			return errors.New("internal server error to get employee task")
		}
		if len(employees) == 0 {
			e.sugar.WithContext(ctx).Errorf("error employees from task: no employee available")
			return errors.New("empty employee")
		}

		historyRelatedID = corporateID.String()

	} else if assigningMode == entity.ASSIGN_MODE_DEPARTMENT {

		tempEmployee, err := e.repository.GetEmployeeByCorporateIDDeptCode(ctx, corporateID.String(), relatedID)
		if err != nil {
			e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
			return errors.New("internal server error to get employee task")
		}
		if len(tempEmployee) == 0 {
			e.sugar.WithContext(ctx).Errorf("error employees from task: no employee available")
			return errors.New("employee list is empty")
		}
		employees = make([]entity.CorporateEmployeeEntityWithDepartment, len(tempEmployee))
		for i, v := range tempEmployee {
			employees[i] = v
		}

		historyRelatedID = tempEmployee[0].CorporateDepartment.ID.String()
	} else {
		employee, err := e.repository.GetEmployeeByUserID(ctx, relatedID)
		if err != nil {
			e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
			return errors.New("internal server error to get employee task")
		}

		if employee == (entity.CorporateEmployeeEntityWithDepartment{}) {
			e.sugar.WithContext(ctx).Errorf("employee does not exist")
			return errors.New("employee does not exist")
		}
		historyRelatedID = employee.PublicUserID.String()
		employees = append(employees, employee)
	}
	userIDArr := make([]string, len(employees))
	for i, v := range employees {
		userIDArr[i] = v.PublicUserID.String()
	}

	existingTasks, err := e.repository.GetEmployeeTaskByUserIDProgramIDTaskType(ctx, userIDArr, programIDArr, entity.GetProgramRelatedTaskType(), entity.GetAllTaskStatus())
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}

	var inserTaskArr []entity.EmployeeTaskEntity
	// list of enrollment to filter out
	enrolledExistingTask, err := e.repository.GetEnrollmentByUserIDArr(ctx, userIDArr)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}
	flagHasTaskExists, flagHasTaskProgramEnrolled := false, false
	for _, employee := range employees {
		for _, programID := range programIDArr {
			programUUID, err := uuid.Parse(programID)

			if err != nil {
				e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
				return errors.New("internal server error to get employee task")
			}
			// filter out already existing tasks
			flagThisTaskExist, flagThisTaskProgramEnrolled := false, false
			for _, existTask := range existingTasks {
				if existTask.UserID == employee.PublicUserID && existTask.TaskRelatedID == programUUID {
					flagThisTaskExist, flagHasTaskExists = true, true
					break
				}
			}

			// filter out tasks that the program is being enrolled with

			for _, enrollment := range enrolledExistingTask {
				if employee.PublicUserID == enrollment.UserID && enrollment.ProgramID == programUUID {
					flagThisTaskProgramEnrolled, flagHasTaskProgramEnrolled = true, true
					break
				}

			}
			if !flagThisTaskExist && !flagThisTaskProgramEnrolled {
				inserTaskArr = append(inserTaskArr, entity.EmployeeTaskEntity{
					UserID:        employee.PublicUserID,
					Email:         employee.Email,
					Department:    employee.Department,
					EmployeeName:  employee.EmployeeName,
					AssignDate:    &assignDate,
					TaskType:      entity.TYPE_PROGRAM_ASSIGNMENT,
					TaskRelatedID: programUUID,
					TaskStatus:    entity.STATUS_ACTIVE,
					RecordFlag:    "ACTIVE",
				})
			}
		}
	}

	if confirmIgnoreTakenProgram == false && (flagHasTaskExists || flagHasTaskProgramEnrolled) {
		e.sugar.WithContext(ctx).Error("some programs have been taken")
		return errors.New(string(commonmodel.ERR_PROGRAMS_HAVE_BEEN_TAKEN))
	}

	if len(inserTaskArr) == 0 {
		e.sugar.WithContext(ctx).Errorf("no tasks appended")
		return errors.New(string(commonmodel.ERR_NO_EMPLOYEE_TASK_ADDED))
	}

	err = e.repository.InsertEmployeeTaskProgramAssignment(ctx, inserTaskArr)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}

	err = e.repository.InsertEmployeeProgramAssignmentHistory(ctx, assigningMode, programIDArr, corporateID.String(), historyRelatedID, corporateID.String())
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}
	return nil
	// get corporate_code from context
}

// GetAllEmployeeTask
// relatedID could be:
// corporate_id, user_id, department_code
func (e *EmployeeTaskServiceImpl) GetEmployeeTaskProgramAssignment(ctx context.Context, assignedProgramPaging *commonmodel.PaginationResponseModel) (err error) {
	tempCorporateID := ctx.Value("corporate_id")
	if tempCorporateID == nil {
		e.sugar.WithContext(ctx).Errorf("error retrieveing corporate_id")
		return errors.New("internal server error to get employee task")

	}

	corporateID := tempCorporateID.(uuid.UUID)

	err = e.repository.GetEmployeeTaskProgramAssignment(ctx, assignedProgramPaging, corporateID.String())
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}

	return nil
	// get corporate_code from context

}

// GetAllEmployeeTask
// relatedID could be:
// corporate_id, user_id, department_code
func (e *EmployeeTaskServiceImpl) GetProgramAssignmentHistory(ctx context.Context, programAssignmentHistory *commonmodel.PaginationResponseModel) (err error) {
	tempCorporateID := ctx.Value("corporate_id")
	if tempCorporateID == nil {
		e.sugar.WithContext(ctx).Errorf("error retrieveing corporate_id")
		return errors.New("internal server error to get employee task")

	}

	corporateID := tempCorporateID.(uuid.UUID)

	programAssignmentHistoryArr, totalData, err := e.repository.GetProgramAssignmentHistory(ctx, corporateID.String(), programAssignmentHistory.MetaData.PageSize, programAssignmentHistory.MetaData.GetOffset())
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return errors.New("internal server error to get employee task")
	}
	for i, history := range programAssignmentHistoryArr {
		if history.AssigningType == entity.ASSIGN_MODE_DEPARTMENT {
			relatedData, err := e.repository.GetDepartment(ctx, history.RelatedID.String())
			if err != nil {
				e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
				return errors.New("internal server error to get program assignment history")
			}

			programAssignmentHistoryArr[i].RelatedData = entity.EmployeeProgramAssignmentHistoryRelatedData{
				ID:   relatedData.ID.String(),
				Name: relatedData.Name,
			}

		} else if history.AssigningType == entity.ASSIGN_MODE_INDIVIDUAL {
			relatedData, err := e.repository.GetEmployeeByUserID(ctx, history.RelatedID.String())
			if err != nil {
				e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
				return errors.New("internal server error to get program assignment history")
			}

			programAssignmentHistoryArr[i].RelatedData = entity.EmployeeProgramAssignmentHistoryRelatedData{
				ID:   relatedData.ID.String(),
				Name: relatedData.EmployeeName,
			}

		}
	}
	programAssignmentHistory.RawData = programAssignmentHistoryArr
	programAssignmentHistory.MetaData.TotalData = totalData
	programAssignmentHistory.MetaData.DataPerPage = len(programAssignmentHistoryArr)

	return nil
}

// GetAllEmployeeTask
// relatedID could be:
// corporate_id, user_id, department_code
func (e *EmployeeTaskServiceImpl) SearchPrograms(ctx context.Context, keyword string) ([]consumermodel.MProgram, error) {

	output, err := e.repository.SearchPrograms(ctx, strings.ToUpper(keyword))
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return nil, errors.New("internal server error to get employee task")
	}

	for i, v := range output {
		output[i].ThumbnailSource = e.getGooglePublicObject(common.EMPLOYEE_TASK_BUCKET, v.ThumbnailUrl)
	}

	return output, nil
	// get corporate_code from context
}

// GetPrograms
func (e *EmployeeTaskServiceImpl) GetPrograms(ctx context.Context) ([]consumermodel.MProgram, error) {

	output, err := e.repository.GetPrograms(ctx)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return nil, errors.New("internal server error to get employee task")
	}

	for i, v := range output {
		output[i].ThumbnailSource = e.getGooglePublicObject(common.EMPLOYEE_TASK_BUCKET, v.ThumbnailUrl)
	}

	return output, nil
}

func (e *EmployeeTaskServiceImpl) GetCorporateServices(ctx context.Context, params model.MCorporateServiceRequest) ([]entity.MCorporateService, int, error) {

	output, totalData, err := e.repository.GetCorporateServices(ctx, params)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return nil, 0, errors.New("internal server error to get employee task")
	}

	for i, v := range output {
		output[i].ImageURL = e.getGooglePublicObject(common.CORPORATE_SERVICE_BOOKING_BUCKET, v.ImageURLSource)
	}

	return output, totalData, nil
}

func (e *EmployeeTaskServiceImpl) GetCorporateServiceByID(ctx context.Context, id string) (entity.MCorporateService, error) {

	output, err := e.repository.GetCorporateServiceByID(ctx, id)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error employees from task:%v", err.Error())
		return entity.MCorporateService{}, errors.New("internal server error to get employee task")
	}

	output.ImageURL = e.getGooglePublicObject(common.CORPORATE_SERVICE_BOOKING_BUCKET, output.ImageURLSource)
	return output, nil
}

func (e *EmployeeTaskServiceImpl) InsertCorporateServiceBooking(ctx context.Context, serviceID string, scheduledAt time.Time) (err error) {
	tempCorporateID := ctx.Value("corporate_id")
	if tempCorporateID == nil {
		e.sugar.WithContext(ctx).Errorf("error retrieveing corporate_id")
		return errors.New("internal server error to get employee task")

	}

	corporateID := tempCorporateID.(uuid.UUID)

	err = e.repository.InsertCorporateServiceBooking(ctx, serviceID, corporateID.String(), scheduledAt)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error InsertCorporateServiceBooking:%v", err.Error())
		return errors.New("internal server error InsertCorporateServiceBooking")
	}

	err = e.sendBookingInfoEmail(ctx, serviceID, scheduledAt)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error InsertCorporateServiceBooking:%v", err.Error())
		return errors.New("internal server error InsertCorporateServiceBooking")
	}
	return nil
	// get corporate_code from context
}

func (e *EmployeeTaskServiceImpl) GetCorporateServiceBooking(ctx context.Context, limit, offset int) ([]entity.CorporateServiceBooking, int, error) {
	tempCorporateID := ctx.Value("corporate_id")
	if tempCorporateID == nil {
		e.sugar.WithContext(ctx).Errorf("error retrieveing corporate_id")
		return nil, 0, errors.New("internal server error to get employee task")

	}

	corporateID := tempCorporateID.(uuid.UUID)

	output, totalData, err := e.repository.GetCorporateServiceBooking(ctx, corporateID.String(), limit, offset)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error GetCorporateServiceBooking:%v", err.Error())
		return nil, 0, errors.New("internal server error GetCorporateServiceBooking")
	}

	return output, totalData, nil
}

func (e *EmployeeTaskServiceImpl) SearchCorporateService(ctx context.Context, keyword string) ([]entity.MCorporateService, error) {
	output, err := e.repository.SearchCorporateService(ctx, keyword)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error SearchCorporateService:%v", err.Error())
		return nil, errors.New("internal server error SearchCorporateService")
	}

	for i, v := range output {
		output[i].ImageURL = e.getGooglePublicObject(common.CORPORATE_SERVICE_BOOKING_BUCKET, v.ImageURLSource)
	}

	return output, nil
}

// insert new well being caused by scheduler trigger
func (e *EmployeeTaskServiceImpl) RenewAllWellBeingTaskSvc(ctx context.Context) (err error) {
	start, end := e.commonSvc.GenerateFirstEndMonth(time.Now().Add(56 * time.Hour))
	e.sugar.WithContext(ctx).Infof("updating well being with end period:%v", end.String())
	// upserting new well being
	if err = e.repository.RenewAllWellBeingTask(ctx, start, end); err != nil {
		e.sugar.WithContext(ctx).Errorf("error in upserting new well being:%v", err.Error())
	}
	return err
}
