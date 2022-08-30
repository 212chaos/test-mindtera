package repo_employeetask

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	cm "github.com/mindtera/go-common-module/common/v2/model"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	fcommonmodel "github.com/mindtera/go-common-module/fajar/model"
	"gorm.io/gorm/clause"

	consumermodel "github.com/mindtera/consumer-service/model"
	dashentity "github.com/mindtera/dashboard-auth/entity"
)

type EmployeeTaskRepositoryImpl struct {
	sugar     logger.CustomLogger
	pgConfig  gormpg.PostgresConfig
	assert    assert.Assert
	commonSvc commonsvc.CommonService
}

// constructor function to create new employee task repository DI
func NewEmployeeTaskRepository(sugar logger.CustomLogger,
	assert assert.Assert, pgConfig gormpg.PostgresConfig,
	commonSvc commonsvc.CommonService) EmployeeTaskRepository {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating employee task value table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.EmployeeTaskEntity{})
	}

	return &EmployeeTaskRepositoryImpl{
		assert:    assert,
		sugar:     sugar,
		pgConfig:  pgConfig,
		commonSvc: commonSvc,
	}
}

// upsert batch employee task repository
func (e *EmployeeTaskRepositoryImpl) UpsertBatchEmployeeTaskRepository(ctx context.Context, employeeTasks *[]entity.EmployeeTaskEntity) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)

	// caluse if corporate program reco still active
	if len(*employeeTasks) == 1 {
		firstElm := (*employeeTasks)[0]
		if firstElm.TaskType == entity.TYPE_MINDTERA_PROGRAM_RECOM && firstElm.TaskStatus == "ACTIVE" {
			var uid []uuid.UUID
			tx.Model(&entity.EmployeeTaskEntity{}).
				Select("task_related_id").
				Where("task_status = 'ACTIVE' AND record_flag = 'ACTIVE'").
				Where("user_id = ? AND subscription_id = ?", firstElm.UserID, firstElm.SubscriptionID).
				Order("created_at DESC").
				First(&uid)
			if len(uid) > 0 {
				if uid[0] == firstElm.TaskRelatedID {
					return nil
				}
			}
		}
	}

	tx.Model(&entity.EmployeeTaskEntity{}).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "subscription_id"},
					{Name: "user_id"},
					{Name: "task_type"},
					{Name: "assign_date"},
				},
				UpdateAll: true,
			},
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "id"},
				},
				UpdateAll: true,
			}).
		Create(employeeTasks)

	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Error(err.Error())

		tx.Rollback()
		return err
	}
	return err
}

// get employee task based on user id, corporate subscription, and task type
func (e *EmployeeTaskRepositoryImpl) GetEmployeeTaskByUserIDSubscription(ctx context.Context, employeeTask *entity.EmployeeTaskEntity,
	employeeTasks *[]entity.EmployeeTaskEntity, pagingModel *cm.PaginationModel) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)

	var total int64
	if pagingModel.Page != -1 {
		tx.Model(&entity.EmployeeTaskEntity{}).
			Select("count(1)").
			Where(employeeTask).
			Where(`(employee_tasks.record_flag = 'ACTIVE' OR employee_tasks.record_flag = 'FINISHED') 
				AND employee_tasks.deleted_at IS NULL`).
			Joins(`left join corporate_employee_entities on 
				employee_tasks.email = corporate_employee_entities.email`).
			Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
			Count(&total)
	}

	var tempMap []map[string]any
	tx.Model(&entity.EmployeeTaskEntity{}).
		Scopes(e.pgConfig.PaginateQuery(pagingModel)).
		Select(e.getColumnsNeededJoins()).
		Where(employeeTask).
		Where(`(employee_tasks.record_flag in ('ACTIVE', 'FINISHED')) 
			AND employee_tasks.deleted_at IS NULL`).
		Joins(`left join corporate_employee_entities on 
			employee_tasks.email = corporate_employee_entities.email`).
		Preload("DepartmentDetail", "record_flag = 'ACTIVE'").
		Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
		Order("employee_tasks.created_at ASC").
		Find(&employeeTasks)
	e.commonSvc.ObjectMapper(&tempMap, employeeTasks)

	// calculate count for employee
	if pagingModel.Page != -1 {
		pagingModel.DataPerPage = len(*employeeTasks)
		pagingModel.TotalData = int(total)
	}

	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Error(err.Error())
		tx.Rollback()
		return err
	}
	return err
}

// get employee task based on user id, corporate subscription, and task type
func (e *EmployeeTaskRepositoryImpl) SearchEmployeeTaskByUserIDSubscription(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity, pagingModel *cm.PaginationModel) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)

	var total int64
	if pagingModel.Page != -1 {
		tx.Model(&entity.EmployeeTaskEntity{}).
			Select(e.getColumnsNeededJoins()).
			Where(`LOWER(employee_tasks.email) LIKE ? OR LOWER(corporate_employee_entities.department) LIKE ? OR LOWER(corporate_employee_entities.department) LIKE ?`, employeeTask.Email, employeeTask.Department, employeeTask.EmployeeName).
			Where("employee_tasks.task_type = ?", employeeTask.TaskType).
			Where("employee_tasks.subscription_id = ? AND (employee_tasks.record_flag IN ('ACTIVE', 'FINISHED'))  AND employee_tasks.deleted_at IS NULL", employeeTask.SubscriptionID).
			Joins(`left join corporate_employee_entities on 
				employee_tasks.email = corporate_employee_entities.email`).
			Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
			Count(&total)
	}

	var abc []map[string]any
	tx.Model(&entity.EmployeeTaskEntity{}).
		Scopes(e.pgConfig.PaginateQuery(pagingModel)).
		Select(e.getColumnsNeededJoins()).
		Where("employee_tasks.task_type = ?", employeeTask.TaskType).
		Where(`LOWER(employee_tasks.email) LIKE ? OR LOWER(corporate_employee_entities.department) LIKE ? OR LOWER(corporate_employee_entities.department) LIKE ?`, employeeTask.Email, employeeTask.Department, employeeTask.EmployeeName).
		Where("employee_tasks.subscription_id = ? AND (employee_tasks.record_flag IN ('ACTIVE', 'FINISHED')) AND employee_tasks.deleted_at IS NULL", employeeTask.SubscriptionID).
		Joins(`left join corporate_employee_entities on 
			employee_tasks.email = corporate_employee_entities.email`).
		Preload("DepartmentDetail").
		Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
		Order("employee_tasks.created_at ASC").
		Find(&employeeTasks)

	emChan := make(chan []entity.EmployeeTaskEntity)
	emErrChan := make(chan error)
	go func() {
		b, er := json.Marshal(abc)
		if er != nil {
			emErrChan <- er
			return
		}
		var em []entity.EmployeeTaskEntity
		if er = json.Unmarshal(b, &em); er != nil {
			emErrChan <- er
			return
		}
		emErrChan <- nil
		emChan <- em
	}()

	// calculate count for employee
	if pagingModel.Page != -1 {
		pagingModel.DataPerPage = len(abc)
		pagingModel.TotalData = int(total)
	}

	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Error(err.Error())
		tx.Rollback()
		return err
	}

	if err = <-emErrChan; err != nil {
		e.sugar.WithContext(ctx).Error(err.Error())
		return err
	}
	(*employeeTasks) = <-emChan

	return err
}

// get employee task based on corp id and Department
func (e *EmployeeTaskRepositoryImpl) GetEmployeeTaskByDepartmentAndSubsId(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity) (err error) {
	types := strings.Split(string(employeeTask.TaskType), ",")

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.EmployeeTaskEntity{}).
		Select(e.getColumnsNeededJoins()).
		Where("employee_tasks.task_type IN ? AND employee_tasks.subscription_id = ?", types, employeeTask.SubscriptionID).
		Where("employee_tasks.record_flag = 'ACTIVE' AND employee_tasks.deleted_at IS NULL").
		Joins(`left join corporate_employee_entities on 
			employee_tasks.email = corporate_employee_entities.email`).
		Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
		Order("employee_tasks.created_at ASC")

	// department selector
	if !e.assert.IsEmpty(employeeTask.Department) {
		tx = tx.Where("corporate_employee_entities.department = ?", employeeTask.Department)
	}
	tx.Find(employeeTasks)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return err
	}
	return err
}

// get employee task by user id and array of tsak type
func (e *EmployeeTaskRepositoryImpl) GetEmployeeTaskByUserIDAndTypes(ctx context.Context, query model.EmployeeTaskQuery, employeeTasks *[]entity.EmployeeTaskEntity) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)

	var additionalQuery string
	if strings.Contains(strings.Join(query.Types, ","),
		string(entity.TYPE_ONBOARDING_QUIZ)) {
		additionalQuery = "OR employee_tasks.due_date IS NULL"
	}

	tx.Model(&entity.EmployeeTaskEntity{}).
		Select(e.getColumnsNeeded()).
		Where(fmt.Sprintf(`employee_tasks.task_type IN ? and employee_tasks.user_id = ? 
			and ((employee_tasks.due_date > ? and employee_tasks.due_date <= ?) %v)`,
			additionalQuery),
			query.Types, query.UserId, query.StartDate, query.EndDate).
		Where("employee_tasks.record_flag = 'ACTIVE' OR employee_tasks.task_status = 'ACTIVE'").
		Order("employee_tasks.created_at ASC").
		Find(employeeTasks)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
	}
	return err
}

// get employee task based on Corporate Id and Task type
func (e *EmployeeTaskRepositoryImpl) GetTaskBySubscriptionIdAndType(ctx context.Context, employeeTask *entity.EmployeeTaskEntity, employeeTasks *[]entity.EmployeeTaskEntity) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.EmployeeTaskEntity{}).
		Select(e.getColumnsNeededJoins()).
		Where(employeeTask).
		Where("employee_tasks.record_flag = 'ACTIVE' AND employee_tasks.deleted_at IS NULL").
		Joins(`left join corporate_employee_entities on 
			employee_tasks.email = corporate_employee_entities.email`).
		Where("corporate_employee_entities.record_flag = 'ACTIVE' AND corporate_employee_entities.deleted_at IS NULL").
		Order("employee_tasks.created_at ASC").
		Find(employeeTasks)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
	}
	return err
}

func (e *EmployeeTaskRepositoryImpl) GetCorporateDetails(ctx context.Context, corporateID string) (entity.CorporateEntity, error) {
	var corporate entity.CorporateEntity

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Where("id = ? AND record_flag = 'ACTIVE'", corporateID).
		Find(&corporate)

	if err := tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return entity.CorporateEntity{}, err
	}

	return corporate, nil
}

func (e *EmployeeTaskRepositoryImpl) GetEmployeeByUserID(ctx context.Context, userID string) (output entity.CorporateEmployeeEntityWithDepartment, err error) {
	var employeeGormable model.CorporateEmployeeEntityGormable2

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Raw(getEmployeeByUserID(userID)).
		Find(&employeeGormable)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return entity.CorporateEmployeeEntityWithDepartment{}, err
	}

	if employeeGormable == (model.CorporateEmployeeEntityGormable2{}) {
		return entity.CorporateEmployeeEntityWithDepartment{}, nil
	}

	return employeeGormable.ToCorporateEmployeeEntity(), err
}

func (e *EmployeeTaskRepositoryImpl) GetEmployeeTaskByUserIDProgramIDTaskType(ctx context.Context, userIDArr []string, programIDArr []string, taskTypeArr []entity.EMPLOYEE_TASK_TYPE, taskStatusArr []entity.EMPLOYEE_TASK_STATUS) ([]entity.EmployeeTaskEntityForConsumer, error) {
	var tasks []entity.EmployeeTaskEntityForConsumer

	tx := e.pgConfig.GenerateTransaction(ctx)

	tx.Where("user_id IN ? AND task_type IN ? AND task_related_id IN ? AND task_status IN ? AND record_flag = 'ACTIVE'",
		userIDArr, taskTypeArr, programIDArr, taskStatusArr,
	).Find(&tasks)
	if err := tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return nil, err
	}

	return tasks, nil
}

func (e *EmployeeTaskRepositoryImpl) GetEnrollmentByUserIDArr(ctx context.Context, userIDArr []string) ([]fcommonmodel.ProgramEnrollmentHistory, error) {
	var enrollmentArr []fcommonmodel.ProgramEnrollmentHistory

	tx := e.pgConfig.GenerateTransaction(ctx)

	tx.Where("user_id IN ? AND record_flag = 'ACTIVE'",
		userIDArr,
	).Find(&enrollmentArr)
	if err := tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return nil, err
	}

	return enrollmentArr, nil
}

func (e *EmployeeTaskRepositoryImpl) GetEmployeeByCorporateID(ctx context.Context, corporateID string) (outputArr []entity.CorporateEmployeeEntityWithDepartment, err error) {
	var employeeGormables []model.CorporateEmployeeEntityGormable2

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Raw(getEmployeeByCorporateID(corporateID)).
		Find(&employeeGormables)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return nil, err
	}

	outputArr = make([]entity.CorporateEmployeeEntityWithDepartment, len(employeeGormables))
	for i, v := range employeeGormables {
		outputArr[i] = v.ToCorporateEmployeeEntity()
	}

	return outputArr, err
}

// GetEmployee3
// Get employee list where corporateID, departmentCode
func (e *EmployeeTaskRepositoryImpl) GetEmployeeByCorporateIDDeptCode(ctx context.Context, corporateID, DepartmentCode string) (outputArr []entity.CorporateEmployeeEntityWithDepartment, err error) {
	var employeeGormables []model.CorporateEmployeeEntityGormable2

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Raw(getEmployeeByCorporateIDDepartmentCode(corporateID, DepartmentCode)).
		Find(&employeeGormables)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return nil, err
	}

	outputArr = make([]entity.CorporateEmployeeEntityWithDepartment, len(employeeGormables))
	for i, v := range employeeGormables {
		outputArr[i] = v.ToCorporateEmployeeEntity()
	}

	return outputArr, err
}

func (e *EmployeeTaskRepositoryImpl) GetDepartment(ctx context.Context, departmentID string) (output dashentity.CorporateDepartmentEntity, err error) {

	tx := e.pgConfig.GenerateTransaction(ctx)
	tx.Model(&dashentity.CorporateDepartmentEntity{}).
		Where("record_flag = 'ACTIVE'").
		Where("id = ?", departmentID).
		Find(&output)

	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		return dashentity.CorporateDepartmentEntity{}, err
	}

	return output, err
}

func (e *EmployeeTaskRepositoryImpl) InsertEmployeeTaskProgramAssignment(ctx context.Context, tasks []entity.EmployeeTaskEntity) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	err = tx.Exec(insertTaskEmployee(tasks)).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return err
	}
	return err
}

func (e *EmployeeTaskRepositoryImpl) InsertEmployeeProgramAssignmentHistory(ctx context.Context, assigningMode entity.TypeAssigningMode, programIDArr []string, corporateID, relatedID, createdBy string) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	err = tx.Exec(insertProgramAssignmentHistory(programIDArr, corporateID, assigningMode, relatedID, createdBy)).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return err
	}
	return err
}

func (e *EmployeeTaskRepositoryImpl) GetProgramAssignmentHistory(ctx context.Context, corporateID string, limit, offset int) ([]entity.EmployeeProgramAssignmentHistory2, int, error) {
	var (
		gormables []model.EmployeeProgramAssignmentHistoryGormable2
		totalData int
	)
	tx := e.pgConfig.GenerateTransaction(ctx)
	err := tx.
		Raw(getProgramAssignmentHistory(corporateID, limit, offset)).
		Find(&gormables).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, 0, err
	}

	err = tx.Raw(countAssignedProgramHistory(corporateID)).
		Scan(&totalData).
		Error
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, 0, err
	}
	output := make([]entity.EmployeeProgramAssignmentHistory2, len(gormables))
	for i, v := range gormables {
		output[i] = v.ToEmployeeProgramAssignmentHistory2()
	}
	return output, totalData, nil
}

func (e *EmployeeTaskRepositoryImpl) GetEmployeeTaskProgramAssignment(ctx context.Context, taskProgramAssignment *cm.PaginationResponseModel, corporateID string) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	var employeeTasksGormable []model.EmployeeTaskEntityGormable
	err = tx.
		Raw(getAssignedProgram(corporateID, entity.TYPE_PROGRAM_ASSIGNMENT, entity.STATUS_ACTIVE, taskProgramAssignment.MetaData.PageSize, taskProgramAssignment.MetaData.GetOffset())).
		Find(&employeeTasksGormable).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return err
	}

	err = tx.Raw(countAssignedProgram(corporateID, entity.TYPE_PROGRAM_ASSIGNMENT, entity.STATUS_ACTIVE)).
		Scan(&taskProgramAssignment.MetaData.TotalData).
		Error
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return err
	}
	employeeTasksGormableLen := len(employeeTasksGormable)
	output := make([]entity.EmployeeTaskEntityForConsumer, employeeTasksGormableLen)
	for i, v := range employeeTasksGormable {
		output[i] = v.ToEmployeeTaskEntityForConsumer()
	}

	taskProgramAssignment.RawData = output
	taskProgramAssignment.MetaData.DataPerPage = employeeTasksGormableLen
	return err
}

func (e *EmployeeTaskRepositoryImpl) SearchPrograms(ctx context.Context, keyword string) ([]consumermodel.MProgram, error) {
	var err error
	tx := e.pgConfig.GenerateTransaction(ctx)
	var mPrograms []consumermodel.MProgram
	err = tx.
		Raw(searchPrograms(keyword)).
		Find(&mPrograms).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, err
	}

	return mPrograms, err
}

func (e *EmployeeTaskRepositoryImpl) GetPrograms(ctx context.Context) ([]consumermodel.MProgram, error) {
	var err error
	tx := e.pgConfig.GenerateTransaction(ctx)

	var (
		mPrograms        []consumermodel.MProgram
		programGormables []consumermodel.MProgramGormable2
	)

	err = tx.
		Raw(getPrograms()).
		Find(&programGormables).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, err
	}

	mPrograms = make([]consumermodel.MProgram, len(programGormables))
	for idx, progGormable := range programGormables {
		mPrograms[idx] = progGormable.ToMProgram()
	}

	return mPrograms, err
}

func (e *EmployeeTaskRepositoryImpl) GetCorporateServices(ctx context.Context, params model.MCorporateServiceRequest) ([]entity.MCorporateService, int, error) {
	var (
		err       error
		totalData int
	)
	tx := e.pgConfig.GenerateTransaction(ctx)
	var corporateServiceGormableArr []model.MCorporateServiceGormable
	tx.
		Raw(getCorporateServices(params.IncludeWebinar, params.IncludeExercise, params.IncludeCoaching, params.PageSize, params.GetOffset())).
		Find(&corporateServiceGormableArr)

	tx.
		Raw(countGetCorporateServices(params.IncludeWebinar, params.IncludeExercise, params.IncludeCoaching)).Scan(&totalData)

	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, 0, err
	}

	output := make([]entity.MCorporateService, len(corporateServiceGormableArr))
	for i, v := range corporateServiceGormableArr {
		output[i] = v.ToMCorporateService()
	}

	return output, totalData, nil
}

func (e *EmployeeTaskRepositoryImpl) GetCorporateServiceByID(ctx context.Context, id string) (entity.MCorporateService, error) {
	var err error
	tx := e.pgConfig.GenerateTransaction(ctx)
	var corporateServiceGormable model.MCorporateServiceGormable
	err = tx.
		Raw(getCorporateService(id)).
		Find(&corporateServiceGormable).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return entity.MCorporateService{}, err
	}

	return corporateServiceGormable.ToMCorporateService(), nil
}

func (e *EmployeeTaskRepositoryImpl) InsertCorporateServiceBooking(ctx context.Context, serviceID, corporateID string, scheduledAt time.Time) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	err = tx.Exec(insertCorporateServiceBooking([]entity.CorporateServiceBooking{
		{
			ServiceID:   serviceID,
			CorporateID: corporateID,
			ScheduledAt: scheduledAt,
			Status:      entity.WAITING_CONFIRM_STATUS,
		},
	})).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return err
	}
	return err
}

func (e *EmployeeTaskRepositoryImpl) GetCorporateServiceBooking(ctx context.Context, corporateID string, limit, offset int) ([]entity.CorporateServiceBooking, int, error) {
	var (
		err       error
		totalData int
	)
	tx := e.pgConfig.GenerateTransaction(ctx)
	var corporateServiceBookingGormableArr []model.CorporateServiceBookingsGormable2

	err = tx.
		Raw(getCorporateServiceBooking(corporateID, limit, offset)).
		Find(&corporateServiceBookingGormableArr).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, 0, err
	}

	err = tx.Raw(countCorporateServiceBooking(corporateID)).
		Scan(&totalData).
		Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, 0, err
	}

	output := make([]entity.CorporateServiceBooking, len(corporateServiceBookingGormableArr))
	for i, v := range corporateServiceBookingGormableArr {
		output[i] = v.ToCorporateServiceBookings()
	}

	return output, totalData, nil
}

func (e *EmployeeTaskRepositoryImpl) SearchCorporateService(ctx context.Context, keyword string) ([]entity.MCorporateService, error) {
	var err error
	tx := e.pgConfig.GenerateTransaction(ctx)
	var mCorporateServiceGormableArr []model.MCorporateServiceGormable
	err = tx.
		Raw(searchCorporateService(keyword)).
		Find(&mCorporateServiceGormableArr).Error

	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
		tx.Rollback()

		return nil, err
	}
	output := make([]entity.MCorporateService, len(mCorporateServiceGormableArr))
	for i, v := range mCorporateServiceGormableArr {
		output[i] = v.ToMCorporateService()
	}

	return output, err
}

// insert new well being caused by scheduler trigger
func (e *EmployeeTaskRepositoryImpl) RenewAllWellBeingTask(ctx context.Context, start, end time.Time) (err error) {
	tx := e.pgConfig.GenerateTransaction(ctx)
	var a any
	tx.Raw(e.insertNewWellBeing(start, end)).
		Scan(&a)
	tx.Raw(e.updateWellBeingPreviousMonth(start, end)).
		Scan(&a)
	if err = tx.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error in transaction:%v", err.Error())
	}
	e.sugar.WithContext(ctx).Infof("successfully updating WELL_BEING with affected row:%v", tx.RowsAffected)
	return err
}
