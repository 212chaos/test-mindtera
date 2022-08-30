package repo_corporateemployee

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	ce "github.com/mindtera/consumer-service/model"
	publicuser "github.com/mindtera/consumer-service/model"
	usersubs "github.com/mindtera/corporate-service/repository/user-subscription"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	cm "github.com/mindtera/go-common-module/common/v2/model"
)

type CorporateEmployeeRepositoryImpl struct {
	sugar        logger.CustomLogger
	pgConfig     gormpg.PostgresConfig
	userSubsRepo usersubs.UserSubscriptionRepo
}

func NewCorporateEmployeeRepository(sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig,
	userSubsRepo usersubs.UserSubscriptionRepo) CorporateEmployeeRepository {
	param := pgConfig.GetParam()

	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate assessment relation table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateEmployeeEntity{})
	}

	return &CorporateEmployeeRepositoryImpl{
		sugar:        sugar,
		pgConfig:     pgConfig,
		userSubsRepo: userSubsRepo,
	}
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeRepository(ctx context.Context, corporateID uuid.UUID,
	corporateEmployees *[]entity.CorporateEmployeeEntity, pagingModel *cm.PaginationModel) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	var total int64
	tx.Table("corporate_employee_entities").
		Select("count(1)").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where(`corporate_id = ? `, corporateID).
		Order("created_at DESC").
		Count(&total)

	//count amount of fetched data as well
	tx.Table("corporate_employee_entities").
		Select("id, public_user_id, email, corporate_id, department, role, employee_name, record_flag").
		Scopes(c.pgConfig.PaginateQuery(pagingModel)).
		Preload("RoleDetail").
		Preload("DepartmentDetail").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("corporate_id = ?", corporateID).
		Order("created_at DESC").
		Find(corporateEmployees)

	(*pagingModel).DataPerPage = len(*corporateEmployees)
	(*pagingModel).TotalData = int(total)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
		return err
	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetEmployeeByCorpAndSubsId(ctx context.Context, corpId, subsId uuid.UUID, corpEmpls *[]entity.CorporateEmployeeEntity) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	//count amount of fetched data as well
	tx.Table("corporate_employee_entities").
		Select("id, public_user_id, email, corporate_id, department, role, employee_name, record_flag, corporate_subscription_id").
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Where("corporate_id = ? AND corporate_subscription_id = ?", corpId, subsId).
		Order("created_at DESC").
		Find(corpEmpls)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())
	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeByName(ctx context.Context,
	filter map[string]string,
	corporateEmployees *[]entity.CorporateEmployeeEntity,
	pagingModel *cm.PaginationModel) (err error) {

	tx := c.pgConfig.GenerateTransaction(ctx)
	//count amount of fetched data as well
	likeFilter := "%" + filter["query"] + "%"

	// total amount of employee
	var total int64
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("count(1)").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where(`corporate_id = ? and
			(LOWER(employee_name) like ? or
				LOWER(email) like ? or LOWER(role) like ? or lower(department) like ?)`, filter["corporate_id"],
			likeFilter, likeFilter, likeFilter, likeFilter).
		Order("created_at DESC").
		Count(&total)

	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("id, public_user_id, email, corporate_id, department, role, employee_name, record_flag").
		Scopes(c.pgConfig.PaginateQuery(pagingModel)).
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where(`corporate_id = ? and 
			(LOWER(employee_name) like ? or 
				LOWER(email) like ? or LOWER(role) like ? or lower(department) like ?)`, filter["corporate_id"],
			likeFilter, likeFilter, likeFilter, likeFilter).
		Preload("RoleDetail").
		Preload("DepartmentDetail").
		Order("created_at DESC").
		Find(corporateEmployees)

	(*pagingModel).DataPerPage = len(*corporateEmployees)
	(*pagingModel).TotalData = int(total)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeByPublicID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error) {

	//count amount of fetched data as well
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("id, email, employee_name, public_user_id, role, department, corporate_id, record_flag").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("public_user_id = ?", corporateEmployees.PublicUserID).
		Find(corporateEmployees)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

		return err
	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeByID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error) {
	//count amount of fetched data as well
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("id, email, corporate_id, record_flag").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("id = ?", corporateEmployees.ID).
		Find(corporateEmployees)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

		return err
	}
	return err
}

// get all active public id by email
func (c *CorporateEmployeeRepositoryImpl) GetUsersByEmail(ctx context.Context, email []string, employees *[]entity.CorporateEmployeeEntity) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("id, email").
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("email in ?", email).
		Find(employees)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeIDRepository(ctx context.Context, corporateID uuid.UUID, corporateEmployees *[]entity.CorporateEmployeeEntity, pagingModel *cm.PaginationModel) (err error) {

	//count amount of fetched data as well
	var count int64
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("id, email, record_flag").
		Scopes(c.pgConfig.PaginateQuery(pagingModel)).
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("corporate_id = ?", corporateID).
		Find(corporateEmployees)

	tx.Model(&entity.CorporateEmployeeEntity{}).Select("count(1)").
		Scopes(c.pgConfig.PaginateQuery(pagingModel)).
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Where("corporate_id = ?", corporateID).
		Count(&count)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmployeeWithQuery(ctx context.Context, corporateEmployee *entity.CorporateEmployeeEntity) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateEmployeeEntity{}).
		Select(`
		corporate_employee_entities.id, 
		corporate_employee_entities.public_user_id, 
		corporate_employee_entities.employee_name, 
		corporate_employee_entities.corporate_id, 
		corporate_employee_entities.department, 
		corporate_employee_entities.role, 
		corporate_employee_entities.corporate_subscription_id, 
		corporate_employee_entities.record_flag, 
		corporate_employee_entities.email`).
		Preload("CorporateDetail", common.ACTIVE_RECORD_FLAG_QUERY,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			}).
		Preload("DepartmentDetail", common.ACTIVE_RECORD_FLAG_QUERY,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, code")
			}).
		Joins(`LEFT JOIN corporate_subscription_entities
			ON corporate_subscription_entities.id = corporate_employee_entities.corporate_subscription_id`).
		Where(corporateEmployee).
		Where("corporate_subscription_entities.record_flag = 'ACTIVE'").
		Find(corporateEmployee)

	// change location for corporate
	if corporateEmployee.CorporateDetail != nil {
		corporateEmployee.Corporate = corporateEmployee.CorporateDetail.Name
	}

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetCorporateEmplByPublicAndSubsID(ctx context.Context, corporateEmployees *entity.CorporateEmployeeEntity) (err error) {
	//count amount of fetched data as well
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select(`corporate_employee_entities.id, 
			corporate_employee_entities.email, 
			corporate_employee_entities.corporate_id, 
			corporate_employee_entities.public_user_id,
			corporate_employee_entities.record_flag`).
		Joins(`LEFT JOIN corporate_subscription_entities cs 
			ON corporate_employee_entities.corporate_subscription_id = cs.id`).
		Where(`corporate_employee_entities.record_flag = 'ACTIVE'`).
		Where("cs.record_flag = 'ACTIVE'").
		Where(`corporate_employee_entities.corporate_subscription_id = ? 
			AND corporate_employee_entities.corporate_id = ?  
			AND cs.end_period > ?`,
			corporateEmployees.CorporateSubscriptionID,
			corporateEmployees.CorporateID,
			time.Now()).
		Find(corporateEmployees)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Errorf("error to query:%v", err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetEmployeeSum(ctx context.Context, corporateID uuid.UUID, count *int64) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateEmployeeEntity{}).
		Select("count(1)").
		Where("corporate_id = ?", corporateID).
		Where(common.ACTIVE_RECORD_FLAG_QUERY_OR_UNREGISTERED).
		Count(count)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetEmployeeDeleted(ctx context.Context, corporateID uuid.UUID, count *int64) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(entity.CorporateEmployeeEntity{}).
		Unscoped().
		Select("count(1)").
		Where("corporate_id = ?", corporateID).
		Where("record_flag = 'DELETED'").
		Count(count)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) UpsertCorporateEmployeeRepository(ctx context.Context, corporateEmployees *[]entity.CorporateEmployeeEntity) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			UpdateAll: true,
		}).
		Create(corporateEmployees)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) DeleteEmployee(ctx context.Context, employeeEntity *entity.CorporateEmployeeEntity) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	if employeeEntity.RecordFlag == "" {
		err = errors.New("entity not found")
		return err
	}

	status := "DEACTIVATED"
	switch employeeEntity.RecordFlag {
	case "ACTIVE":
		employeeEntity.RecordFlag = common.DEACTIVATE_FLAG
	case "DEACTIVATED":
		employeeEntity.RecordFlag = common.DELETE_FLAG
		employeeEntity.DeletedAt = gorm.DeletedAt{
			Valid: true,
			Time:  time.Now()}
	case "UNREGISTERED":
		employeeEntity.RecordFlag = common.HARD_DELETE_FLAG
		employeeEntity.DeletedAt = gorm.DeletedAt{
			Valid: true,
			Time:  time.Now()}
	default:
		employeeEntity.RecordFlag = "UNREGISTERED"
	}
	c.sugar.WithContext(ctx).Infof("deleting employee corporate id:%v employee status :%v", employeeEntity.ID, employeeEntity.RecordFlag)
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"record_flag", "deleted_at"}),
		}).
		Create(employeeEntity).
		Select("id, public_user_id, corporate_subscription_id, email").
		Where("id = ?", employeeEntity.ID).
		Find(employeeEntity)

	if employeeEntity.RecordFlag == common.DELETE_FLAG || employeeEntity.RecordFlag == common.HARD_DELETE_FLAG {
		// deactive all employee tasks
		c.sugar.WithContext(ctx).Infof("deleting employee task employee id:%v employee status :%v", employeeEntity.ID, employeeEntity.RecordFlag)
		tx.Table("employee_tasks").
			Where("email = ? AND subscription_id = ?",
				employeeEntity.Email,
				employeeEntity.CorporateSubscriptionID).
			Updates(map[string]any{
				"deleted_at":  time.Now(),
				"record_flag": "DELETED",
				"task_status": "DELETED"})
	}

	// upsert public subscription status
	subs := ce.Subscription{
		UserID:             employeeEntity.PublicUserID,
		SubscriptionPlanID: &(employeeEntity.CorporateSubscriptionID)}
	c.userSubsRepo.UpdateUserSubscriptionStatusWithTx(ctx, tx, &subs, status)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetEmployeeSumWithTx(ctx context.Context, tx *gorm.DB, corporateID uuid.UUID, count *int64) (err error) {
	tx.Model(entity.CorporateEmployeeEntity{}).
		Where("corporate_id = ?", corporateID).
		Count(count)
	return err
}

// get employee public user data
func (c *CorporateEmployeeRepositoryImpl) GetUserByUserEmail(ctx context.Context, emails []string, publicUser *[]publicuser.User) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Table("users").
		Select("id, email").
		Where("email IN ?", emails).
		Where("record_flag = 'ACTIVE'").
		Find(&publicUser)

	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())

	}
	return err
}

func (c *CorporateEmployeeRepositoryImpl) GetUsersByDepartment(ctx context.Context, corporate_id uuid.UUID, department string, publicId *[]string) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)
	switch department {
	case "ALL":
		tx.Model(&entity.CorporateEmployeeEntity{}).
			Select("public_user_id").
			Where("record_flag = 'ACTIVE' AND corporate_id = ?", corporate_id).
			Find(publicId)
	default:
		tx.Model(&entity.CorporateEmployeeEntity{}).
			Select("public_user_id").
			Where("record_flag = 'ACTIVE' and department = ? AND corporate_id = ?", department, corporate_id).
			Find(&publicId)
	}
	if err = tx.Error; err != nil {
		c.sugar.WithContext(ctx).Error(err.Error())
	}
	return err
}

// get amount filled employee
func (c *CorporateEmployeeRepositoryImpl) GetEmployeeAssessStatus(ctx context.Context, subsId uuid.UUID, assessStatus *model.AssessmentStatus) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	var uid []uuid.UUID
	tx.Model(&entity.CorporateEmployeeEntity{}).
		Select("public_user_id").
		Where(`record_flag = 'ACTIVE'
			AND corporate_subscription_id = ?`, subsId).
		Find(&uid)

	var abc []map[string]any
	tx.Table("employee_tasks").
		Select("task_type, count(1) as amount").
		Where(`task_status = 'FINISHED' 
			AND subscription_id = ?
			AND user_id IN ? 
			AND (task_status != 'DELETED' OR record_flag != 'DELETED')
			AND task_type IN ('PRE_ASSESSMENT','POST_ASSESSMENT','FOLLOWUP_ASSESSMENT')`, subsId, uid).
		Group("task_type").
		Find(&abc)

	for _, v := range abc {
		switch v["task_type"].(string) {
		case "PRE_ASSESSMENT":
			assessStatus.PreAssessmentAmount = int(v["amount"].(int64))
		case "POST_ASSESSMENT":
			assessStatus.PostAssessmentAmount = int(v["amount"].(int64))
		case "FOLLOWUP_ASSESSMENT":
			assessStatus.FollowupAssessmentAmount = int(v["amount"].(int64))
		}
	}

	return tx.Error
}

// get amount filled employee
func (c *CorporateEmployeeRepositoryImpl) GetEmployeeWellBeingAssessStatus(ctx context.Context, subsId uuid.UUID, wellBeing *model.WellBeingStatus) (err error) {
	tx := c.pgConfig.GenerateTransaction(ctx)

	var abc []map[string]any

	switch wellBeing.Department {
	case "ALL":
		tx.Table("employee_tasks").
			Select("task_type, task_status").
			Joins(`LEFT JOIN corporate_employee_entities 
			ON employee_tasks.user_id = corporate_employee_entities.public_user_id`).
			Where(`subscription_id = ?
			AND task_type = 'WELL_BEING'
			AND (employee_tasks.task_status != 'DELETED' OR employee_tasks.record_flag != 'DELETED')
			AND (due_date > ? AND due_date <= ?) 
			AND corporate_employee_entities.record_flag = 'ACTIVE'`, subsId, wellBeing.StartDate, wellBeing.EndDate).
			Find(&abc)
	default:
		tx.Table("employee_tasks").
			Select("task_type, task_status").
			Joins(`LEFT JOIN corporate_employee_entities 
			ON employee_tasks.user_id = corporate_employee_entities.public_user_id`).
			Where(`subscription_id = ?
			AND employee_tasks.task_type = 'WELL_BEING'
			AND (due_date > ? AND due_date <= ?)
			AND (employee_tasks.task_status != 'DELETED' OR employee_tasks.record_flag != 'DELETED')
			AND corporate_employee_entities.department = ?
			AND corporate_employee_entities.record_flag = 'ACTIVE'`,
				subsId, wellBeing.StartDate, wellBeing.EndDate, wellBeing.Department).
			Find(&abc)
	}

	var finished int
	for _, v := range abc {
		if strings.EqualFold(v["task_status"].(string), "FINISHED") {
			finished++
		}
	}
	wellBeing.Status = fmt.Sprintf("%v/%v", finished, len(abc))

	return tx.Error
}
