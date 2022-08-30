package service_corporateemployee

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"
	"github.com/mindtera/go-common-module/common/pb"

	cerepo "github.com/mindtera/corporate-service/repository/corporate-employee"
	scsvc "github.com/mindtera/corporate-service/service/corporate-subscription"
	employeetasksvc "github.com/mindtera/corporate-service/service/employee-task"
	usrsubssvc "github.com/mindtera/corporate-service/service/user-subscription"
	ca "github.com/mindtera/dashboard-auth/entity"
	corprolesvc "github.com/mindtera/dashboard-auth/service/corporate-role"
	deptsvc "github.com/mindtera/dashboard-auth/service/department"
	mailjet "github.com/mindtera/go-common-module/common/client/mailjet-client"
	grpcclnconn "github.com/mindtera/go-common-module/common/v2/configuration/grpc-connection/client"
	cm "github.com/mindtera/go-common-module/common/v2/model"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
	util "github.com/mindtera/go-common-module/common/v2/service/util"
	quizassesscln "github.com/mindtera/quiz-assessment-service/handler/grpc/corporate-assessment/client"
	qm "github.com/mindtera/quiz-assessment-service/model"
)

type CorporateEmployeeServiceImpl struct {
	sugar           logger.CustomLogger
	corpEmplRepo    cerepo.CorporateEmployeeRepository
	redisSvc        redissvc.RedisSvc
	commonSvc       commonsvc.CommonService
	assert          assert.Assert
	util            util.Util
	deptSvc         deptsvc.DepartmentService
	roleService     corprolesvc.CorporateRoleService
	corpSubsSvc     scsvc.CorporateSubscriptionService
	employeeTaskSvc employeetasksvc.EmployeeTaskService
	userSubsSvc     usrsubssvc.UserSubsService
	mailjet         mailjet.MailjetClient
	quizAssessCln   quizassesscln.CorporateAssessmentClient
}

var (
	keyAmount = "CORPORATE_EMPLOYEE_AMOUNT"
)

func NewCorporateEmployeeService(
	sugar logger.CustomLogger,
	corpEmplRepo cerepo.CorporateEmployeeRepository,
	redisSvc redissvc.RedisSvc,
	commonSvc commonsvc.CommonService,
	util util.Util,
	assert assert.Assert,
	deptSvc deptsvc.DepartmentService,
	roleService corprolesvc.CorporateRoleService,
	corpSubsSvc scsvc.CorporateSubscriptionService,
	employeeTaskSvc employeetasksvc.EmployeeTaskService,
	mailjet mailjet.MailjetClient,
	userSubsSvc usrsubssvc.UserSubsService) CorporateEmployeeService {

	// init grpc connection
	conn := grpcclnconn.NewGRPCClientConnection(sugar, os.Getenv("GRPC_QUIZ_ASSESS_SERVICE"))
	quizAssessCln := quizassesscln.NewCorporateAssessmentClient(sugar, conn)

	return &CorporateEmployeeServiceImpl{
		sugar:           sugar,
		mailjet:         mailjet,
		corpEmplRepo:    corpEmplRepo,
		redisSvc:        redisSvc,
		commonSvc:       commonSvc,
		util:            util,
		deptSvc:         deptSvc,
		roleService:     roleService,
		corpSubsSvc:     corpSubsSvc,
		assert:          assert,
		employeeTaskSvc: employeeTaskSvc,
		userSubsSvc:     userSubsSvc,
		quizAssessCln:   quizAssessCln,
	}
}

func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeService(ctx context.Context, corporatePaging *cm.PaginationResponseModel) (err error) {
	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err = errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)
	key := fmt.Sprintf("CORPORATE_EMPLOYEE_PAGE_%v_%v_%v", corporateID, corporatePaging.MetaData.Page, corporatePaging.MetaData.PageSize)
	// try to get data from redis first
	if err = c.redisSvc.Get(ctx, key, &corporatePaging); err != nil {
		c.sugar.WithContext(ctx).Infof("getting from database for employee corporate:%v", corporateID)
		var corpEmpls []entity.CorporateEmployeeEntity
		if err = c.corpEmplRepo.GetCorporateEmployeeRepository(ctx, corporateID, &corpEmpls, &corporatePaging.MetaData); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)

			return err
		}
		corporatePaging.RawData = corpEmpls
	}

	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func() {
		c.sugar.WithContext(ctxBack).Infof("set key value for redis corporate paging:%v", key)
		if err := ctx.Err(); err != nil {
			c.sugar.WithContext(ctxBack).Errorf("error in context when setting redis: %v", err)
		}
		if corporatePaging.MetaData.DataPerPage <= 0 {
			return
		}
		if err := c.redisSvc.Set(ctxBack, key, corporatePaging, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctxBack).Errorf("error when setting cache in redis:%v", err)
		}
	}()
	return err
}

func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeByNameService(ctx context.Context, query string, corporatePaging *cm.PaginationResponseModel) (err error) {

	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err = errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)

	// get employee name
	filter := map[string]string{
		"corporate_id": corporateID.String(),
		"query":        query}

	key := fmt.Sprintf("CORPORATE_EMPLOYEE_PAGE_%v_%v_%v_%v", corporateID, corporatePaging.MetaData.Page, corporatePaging.MetaData.PageSize, query)
	// try to get data from redis first
	if err = c.redisSvc.Get(ctx, key, &corporatePaging); err != nil {
		c.sugar.WithContext(ctx).Infof("getting from database for employee corporate with filter:%v, %v", corporateID, query)
		var corpEmpls []entity.CorporateEmployeeEntity
		if err = c.corpEmplRepo.GetCorporateEmployeeByName(ctx, filter, &corpEmpls, &corporatePaging.MetaData); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)

			return err
		}
		corporatePaging.RawData = corpEmpls
	}
	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func() {
		if err := ctxBack.Err(); err != nil {
			c.sugar.WithContext(ctxBack).Errorf("error in context when setting redis: %v", err)
		}
		if err := c.redisSvc.Set(ctxBack, key, corporatePaging, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when setting cache in redis:%v", err)
		}
	}()
	return err
}

func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeByPublicIDService(ctx context.Context, corporateEmployeeInfo *model.CorporateInformation) (err error) {
	corpEmpls := corporateEmployeeInfo.EmployeeInformation
	key := "CORPORATE_PUBLIC_EMPLOYEE_" + corpEmpls.PublicUserID.String()

	// check redis first
	if err = c.redisSvc.Get(ctx, key, &corpEmpls); err != nil {
		c.sugar.WithContext(ctx).Infof("getting user information from database:%v", corpEmpls.PublicUserID)
		if err = c.corpEmplRepo.GetCorporateEmployeeByPublicID(ctx, &corpEmpls); err != nil {
			c.sugar.WithContext(ctx).Errorf("error getting user information from database:%v", err)
			return err
		}
		if c.assert.IsUUIDEmpty(corpEmpls.ID.String()) {
			c.sugar.WithContext(ctx).Errorf("error user not found")
			return errors.New("error user not found")
		}
	}

	corpSubs := entity.CorporateSubscriptionEntity{CorporateID: corpEmpls.CorporateID}
	if c.corpSubsSvc.GetCorporateSubscriptionService(ctx, &corpSubs); err != nil {
		c.sugar.WithContext(ctx).Errorf("error get corporate subs:%v", err)
		return errors.New("internal server error to get corporate detail")
	}

	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func(k string, model entity.CorporateEmployeeEntity) {
		if !c.assert.IsUUIDEmpty(model.ID.String()) {
			if err := c.redisSvc.Set(ctxBack, k, model, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
				c.sugar.WithContext(ctxBack).Errorf("error set user information to redis:%v", err)
			}
			c.sugar.WithContext(ctxBack).Infof("success set redis value public user:%v", k)
		}
	}(key, corpEmpls)

	corporateEmployeeInfo.EmployeeInformation = corpEmpls
	corporateEmployeeInfo.SubscriptionDetail = corpSubs

	return err
}

func (c *CorporateEmployeeServiceImpl) GetAllCorporateEmployeeHashService(ctx context.Context) (result map[string]bool, err error) {

	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err := errors.New("error to get corporate id context")

		return nil, err
	}
	corporateID := tempID.(uuid.UUID)
	key := fmt.Sprintf("CORPORATE_EMPLOYEE_PAGE_%v_HASH", corporateID)
	// try to get data from redis first
	hashEmployee := map[string]bool{}
	if err := c.redisSvc.Get(ctx, key, &hashEmployee); err != nil {
		c.sugar.WithContext(ctx).Infof("getting from database for employee corporate:%v", corporateID)
		var corpEmpls []entity.CorporateEmployeeEntity
		if err := c.corpEmplRepo.GetCorporateEmployeeIDRepository(ctx, corporateID, &corpEmpls, &cm.PaginationModel{PageSize: -1, Page: -1}); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching data from database:%v", err)

			return nil, err
		}
		// mutex
		wg, mx := c.commonSvc.GenerateWaitGroupAndMutex()
		for _, val := range corpEmpls {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, v entity.CorporateEmployeeEntity) {
				defer c.commonSvc.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				hashEmployee[v.Email] = true
			}(&wg, val)
		}
	}
	// wg.Wait()
	if err := ctx.Err(); err != nil {
		c.sugar.WithContext(ctx).Errorf("error in context when setting redis: %v", err)
	}

	if err := c.redisSvc.Set(context.Background(), key, corporateID, time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when setting cache in redis:%v", err)
	}

	return hashEmployee, nil
}

func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeStatusService(ctx context.Context, corporateID uuid.UUID, corporateStatus *model.CorporateStatus) (err error) {

	// check redis value first
	var amount int64
	if err = c.redisSvc.Get(ctx, keyAmount+corporateID.String(), &amount); err != nil {
		//fetching from database
		if err = c.corpEmplRepo.GetEmployeeSum(ctx, corporateID, &amount); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching data for employee count:%v", err)
			return
		}
	}
	// set redis value
	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func() {
		if err := c.redisSvc.Set(ctxBack, keyAmount+corporateID.String(), amount,
			time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctxBack).Errorf("error when set redis of amount corporate id :%v", corporateID)
		}
	}()

	// get corporate subscription
	corpSubs := entity.CorporateSubscriptionEntity{CorporateID: corporateID}
	if err = c.corpSubsSvc.GetCorporateSubscriptionService(ctx, &corpSubs); err != nil {
		return
	}

	var deletedAmount int64
	if err = c.redisSvc.Get(ctx, keyAmount+corporateID.String()+"_DELETED", &deletedAmount); err != nil {
		//fetching from database
		if err = c.corpEmplRepo.GetEmployeeDeleted(ctx, corporateID, &deletedAmount); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching data for employee count:%v", err)
			return
		}
	}
	// set redis value
	ctxBack = c.commonSvc.ContextBackground(ctx)
	go func() {
		if err := c.redisSvc.Set(ctxBack, keyAmount+corporateID.String()+"_DELETED", &deletedAmount,
			time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctxBack).Errorf("error when set redis of amount corporate id :%v", corporateID)
		}
	}()

	// count limit for corporate
	lim := math.Ceil(0.25*float64(corpSubs.EmployeeCapacity)) - float64(deletedAmount)
	// create model
	corporateStatus.EmployeeDetail = model.EmployeeDetail{
		EmployeeAmount: amount,
		EmployeeLimit:  int64(lim),
	}
	corporateStatus.SubscriptionDetail = corpSubs
	return err
}

// upsert corporate user and doing validation for public user
func (c *CorporateEmployeeServiceImpl) UpsertCorporateEmployeeService(ctx context.Context, corpEmpls *[]entity.CorporateEmployeeEntity) (errorRecord []model.CorporateEmployeeError, err error) {

	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err := errors.New("error to get corporate id context")

		return nil, err
	}
	corporateID := tempID.(uuid.UUID)
	key := fmt.Sprintf("CORPORATE_EMPLOYEE_PAGE_%v_*", corporateID)

	var userName string
	if ctx.Value("user_name") != nil {
		un := ctx.Value("user_name").(*string)
		userName = *un
	}

	var status model.CorporateStatus
	if err = c.GetCorporateEmployeeStatusService(ctx, corporateID, &status); err != nil {
		return nil, err
	}
	amount := status.EmployeeDetail.EmployeeAmount
	subscriptionID := status.SubscriptionDetail.ID

	// getting public user data based on email records
	publicUserChannel := make(chan map[string]bool)
	publicUserIDChannel := make(chan map[string]uuid.UUID)
	publicUserErrChannel := make(chan error)
	go func(pubChan chan map[string]bool, pubUserID chan map[string]uuid.UUID, pubErrChan chan error) {
		pubArr, pubIDArr, er := c.generateUserFormatFromEmail(ctx, corpEmpls)
		if er != nil {
			pubErrChan <- er
			return
		}
		pubErrChan <- nil
		pubChan <- pubArr
		pubUserID <- pubIDArr
	}(publicUserChannel, publicUserIDChannel, publicUserErrChannel)

	// remove redis data
	errChan := make(chan error)
	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func() {
		c.sugar.WithContext(ctxBack).Info("deleting cache for corporate employee")
		keys := []string{key,
			keyAmount + "*",
			keyAmount + "_*",
			"CORPORATE_PUBLIC_EMPLOYEE_*"}
		for _, k := range keys {
			c.redisSvc.Delete(ctxBack, k)
		}
		errChan <- nil
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	wg2, mx2 := c.commonSvc.GenerateWaitGroupAndMutex()
	// get department and role in hash first
	eHashChan := make(chan map[string]bool)
	dHashChan := make(chan map[string]bool)
	erHashChan := make(chan error, 2)

	// defer function
	wg2.Add(1)
	go func(wg_ *sync.WaitGroup, eh chan map[string]bool, er chan error) {
		defer wg_.Done()
		// mx2.Lock()
		temp, err := c.GetAllCorporateEmployeeHashService(ctx)
		if err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching employee HASH: %v", err)
			er <- err
			return
		}
		if temp == nil {
			temp = map[string]bool{}
		}

		er <- nil
		eh <- temp
	}(&wg2, eHashChan, erHashChan)

	// getting department hash
	wg2.Add(1)
	go func(wg_ *sync.WaitGroup, eh chan map[string]bool, er chan error) {
		// defer function
		defer wg_.Done()

		// create new transaction for query department
		ctx_ := context.WithValue(ctx, cm.TRANSACTION_KEY, nil)
		// ctx_ = context.WithValue(ctx_, string(cm.TRANSACTION_KEY), nil)

		departmentHash, err := c.deptSvc.GetDepartmentMap(ctx_, corporateID)
		if err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching department HASH: %v", err)
			er <- err
			return
		}
		if departmentHash == nil {
			departmentHash = map[string]bool{}
		}
		er <- nil
		eh <- departmentHash
		c.sugar.WithContext(ctx).Infof("success")
		// mx2.Unlock()
	}(&wg2, dHashChan, erHashChan)

	// getting role hash
	roleHash, err := c.roleService.GetCorporateRolesHashService(ctx)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error when fetching role HASH: %v", err)
		return nil, err
	}
	if roleHash == nil {
		roleHash = map[string]bool{}
	}

	// callback all channel
	for i := 0; i < 2; i++ {
		if err := <-erHashChan; err != nil {
			c.sugar.WithContext(ctx).Errorf("error in callback channel:%v", err)
			return nil, err
		}
	}
	// callback hashmap
	employeeHash := <-eHashChan
	departmentHash := <-dHashChan

	// error checking from generating user public hash
	if err = <-publicUserErrChannel; err != nil {
		c.sugar.WithContext(ctx).Errorf("error generating public user hash")
		err = errors.New("error fetching data from payload")
		return nil, err
	}
	var publicUserMap map[string]bool
	var publicUserIDMap map[string]uuid.UUID
	if publicUserMap = <-publicUserChannel; publicUserMap == nil {
		publicUserMap = map[string]bool{}
	}
	if publicUserIDMap = <-publicUserIDChannel; publicUserIDMap == nil {
		publicUserIDMap = map[string]uuid.UUID{}
	}

	//init variables
	var correctFormat []entity.CorporateEmployeeEntity
	var insertedDepartment []ca.CorporateDepartmentEntity

	var emails []string
	emailNumber := make(map[string]int)
	emailDuplicated := make(map[string]bool)

	for idx, val := range *corpEmpls {
		wg2.Add(1)
		go func(wg2_ *sync.WaitGroup, v entity.CorporateEmployeeEntity, i int) {
			defer c.commonSvc.EndWaitGroupAndMutex(wg2_, mx2)

			mx2.Lock()
			// check quota
			if c.assert.IsUUIDEmpty(v.ID.String()) {
				if employeeHash[v.Email] {
					c.sugar.WithContext(ctx).Errorf("error email already appended")
					errorRecord = append(errorRecord, model.CorporateEmployeeError{
						Email:        v.Email,
						Number:       v.Number,
						ErrorMessage: "email is duplicated",
					})
					emailDuplicated[v.Email] = true
					return
				}
				if amount >= int64(status.SubscriptionDetail.EmployeeCapacity) {
					c.sugar.WithContext(ctx).Errorf("error adding payload, exceed capacity")
					return
				}
				amount++
			}

			// remove white space from department
			v.Department = c.removeWhiteSpaces(v.Department)

			// validate payload
			employeeError := c.validateEmployeePayload(ctx, employeeHash, roleHash, &v)
			if employeeError != nil {
				errorRecord = append(errorRecord, *employeeError)
				return
			}

			// check department
			originDept := v.Department
			codeDept := strings.ReplaceAll(strings.ToUpper(v.Department), " ", "_")
			if !departmentHash[codeDept] {
				insertedDepartment = append(insertedDepartment, ca.CorporateDepartmentEntity{
					Code:        codeDept,
					Name:        originDept,
					CorporateID: corporateID,
				})
			}
			v.Department = codeDept
			v.CreatedAt = time.Now().Add(time.Duration(i * int(time.Minute)))

			if c.assert.IsUUIDEmpty(v.ID.String()) {
				// append email
				emails = append(emails, v.Email)
				emailNumber[v.Email] = v.Number
				v.CreatedBy = userName
			} else {
				v.UpdatedBy = userName
			}

			// append the user data
			if publicUserMap[v.Email] &&
				(strings.EqualFold(v.RecordFlag, "UNREGISTERED") || v.RecordFlag == "") {
				v.PublicUserID = publicUserIDMap[v.Email]
				v.RecordFlag = "ACTIVE"
			}

			// update value corporate id and append value
			v.CorporateID = corporateID
			v.CorporateSubscriptionID = subscriptionID
			v.Corporate = status.SubscriptionDetail.CorporateDetail.Name
			correctFormat = append(correctFormat, v)

			// update hash value
			departmentHash[v.Department] = true
			departmentHash[codeDept] = true
			employeeHash[v.Email] = true
		}(&wg2, val, idx)
	}
	wg2.Wait()

	// upsert department
	if len(insertedDepartment) > 0 {
		if err := c.deptSvc.UpsertDepartmentBatch(ctx, insertedDepartment); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when upserting department:%v", err)
			return nil, err
		}
		go func() {
			keys := []string{
				"DEPARTMENT_ENTITY_*",
				"CORPORATE_DEPARTMENT_MAP_*",
				"DEPARTMENT_ENTITY"}
			for _, v := range keys {
				c.redisSvc.Delete(context.Background(), v)
			}
		}()
	}

	// upsert employee repo
	if len(correctFormat) > 0 {
		// check email and take out
		var empls []entity.CorporateEmployeeEntity
		if err = c.corpEmplRepo.GetUsersByEmail(ctx, emails, &empls); err != nil {
			c.sugar.WithContext(ctx).Errorf("error fetching employee by email")
			return nil, err
		}
		emplHash := make(map[string]bool)
		for _, v := range empls {
			if !emplHash[v.Email] {
				emplHash[v.Email] = true
			}
		}

		//check email
		var final []entity.CorporateEmployeeEntity
		for _, v := range correctFormat {
			if emplHash[v.Email] {
				errorRecord = append(errorRecord, model.CorporateEmployeeError{
					Email:        v.Email,
					ErrorMessage: "email is duplicated in our system",
					Number:       emailNumber[v.Email],
				})
				continue
			}
			if emailDuplicated[v.Email] {
				errorRecord = append(errorRecord, model.CorporateEmployeeError{
					Email:        v.Email,
					ErrorMessage: "email is duplicated",
					Number:       emailNumber[v.Email],
				})
				continue
			}
			final = append(final, v)
		}
		correctFormat = final
		if len(correctFormat) <= 0 {
			err = errors.New("BAD_REQUEST")
			return errorRecord, err
		}

		// generating first tasks
		if err = c.generateFirstTasks(ctx, correctFormat, status.SubscriptionDetail); err != nil {
			return nil, err
		}
		// activating employee subs
		c.sugar.WithContext(ctx).Info("activating employee valid employee subscription")
		if err = c.userSubsSvc.ActiveUserSubsService(ctx, &correctFormat, &status.SubscriptionDetail); err != nil {
			c.sugar.WithContext(ctx).Errorf("error activating employee subscription :( :%v", err)
		}
		// sending email reminder
		emailContext := c.commonSvc.ContextBackground(ctx)
		emailContext = context.WithValue(emailContext, common.PUBLIC_KEY, ctx.Value(common.PUBLIC_KEY))
		c.sugar.WithContext(ctx).Infof("process sending email registration with public key:%v", emailContext.Value(common.PUBLIC_KEY))
		go c.validateRegirstrationEmail(emailContext, correctFormat)
		// upserting empl data
		if err = c.corpEmplRepo.UpsertCorporateEmployeeRepository(ctx, &correctFormat); err != nil {
			return nil, err
		}
	}

	return errorRecord, err
}

// delete function for deleting condition for corporate employee
func (c *CorporateEmployeeServiceImpl) DeleteCorporateEmployeeService(ctx context.Context, corpEmpls *entity.CorporateEmployeeEntity) (err error) {

	// get corporate_code from context
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err := errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)
	key := fmt.Sprintf("CORPORATE_EMPLOYEE_PAGE_%v_*", corporateID)

	var status model.CorporateStatus
	if err = c.GetCorporateEmployeeStatusService(ctx, corporateID, &status); err != nil {
		return err
	}

	if status.EmployeeDetail.EmployeeLimit <= 0 {
		err = errors.New("exceed limit for deleting employee")
		return err
	}
	ctxBack := c.commonSvc.ContextBackground(ctx)
	go func() {
		c.sugar.WithContext(ctxBack).Info("deleting cache for corporate employee")
		if err := c.redisSvc.Delete(ctxBack, key+"*"); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when deleting data from redis:%v", err)
		}
		if err := c.redisSvc.Delete(ctxBack, keyAmount+corporateID.String()+"*"); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when deleting data from redis:%v", err)
		}
	}()

	if err = c.corpEmplRepo.GetCorporateEmployeeByID(ctx, corpEmpls); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when get corporate employee by id:%v", err)
		return err
	}

	if corpEmpls.CorporateID != corporateID {
		c.sugar.WithContext(ctx).Errorf("error corporate id (%v) is unauthorized to delete this employee (%v)", corporateID, corpEmpls.ID)
		err = errors.New("corporate is unauthorized to delete employee")
		return err
	}
	corpEmpls.CorporateSubscriptionID = status.SubscriptionDetail.ID

	corpEmpls.CorporateID = corporateID
	if err := c.corpEmplRepo.DeleteEmployee(ctx, corpEmpls); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when deleting data from db:%v", err)
	}
	return err
}

// getting corporate employee entity from email
func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeByEmailService(ctx context.Context, corporateEmployee *entity.CorporateEmployeeEntity) (err error) {

	// searching email in database
	c.sugar.WithContext(ctx).Infof("searching corporate employee with registered email:%v", corporateEmployee.Email)
	if err = c.corpEmplRepo.GetCorporateEmployeeWithQuery(ctx, corporateEmployee); err != nil {
		c.sugar.WithContext(ctx).Errorf("cannot fetch corporate employee with error:%v", err.Error())
		return err
	}

	// check is corporate employee exist of not
	if c.assert.IsUUIDEmpty(corporateEmployee.ID.String()) {
		c.sugar.WithContext(ctx).Errorf("user is not part of corporate employee email:%v", corporateEmployee.Email)
		err = errors.New("NOT_EMPLOYEE")
		return err
	}

	c.sugar.WithContext(ctx).Infof("got employee with id:%v flag:%v", corporateEmployee.ID, corporateEmployee.RecordFlag)
	if !(strings.EqualFold(corporateEmployee.RecordFlag, "UNREGISTERED") || strings.EqualFold(corporateEmployee.RecordFlag, "ACTIVE")) {
		err = errors.New("user flag is not valid")
		c.sugar.WithContext(ctx).Errorf("cannot proceed user because flag:%v", corporateEmployee.RecordFlag)
	}
	return err
}

func (c *CorporateEmployeeServiceImpl) GetCorporateEmployeeDashboardService(ctx context.Context, query *model.CorpQuery, assessment *qm.CorporateAssessmentResponseModel) (errMsg cm.ErrorMessage) {
	ctx_, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	// get corp subs
	var corpStatus model.CorporateStatus
	if err := c.GetCorporateEmployeeStatusService(ctx_, query.CorporateId, &corpStatus); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when proceed corporate status service:%v", err.Error())
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
			Error:     err,
		}
	}
	var subs pb.Subscription
	if err := c.commonSvc.ObjectMapper(&corpStatus.SubscriptionDetail, &subs); err != nil {
		c.sugar.WithContext(ctx).Errorf("error mapping object:%v", err.Error())
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
			Error:     err,
		}
	}
	cid, err := uuid.Parse(subs.CorporateId)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error parsing corporate id:%v", err.Error())
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
			Error:     err,
		}
	}

	// get user by department
	publicId := strings.Split(query.Users, ";")
	if !c.assert.IsEmpty(query.Department) {
		if err = c.corpEmplRepo.GetUsersByDepartment(ctx_, cid, query.Department, &publicId); err != nil {
			c.sugar.WithContext(ctx).Errorf("error getting corporate employee:%v", err.Error())
			return cm.ErrorMessage{
				ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
				Error:     err,
			}
		}
	}

	if len(publicId) <= 0 {
		c.sugar.WithContext(ctx).Errorf("public id is empty")
		return cm.ErrorMessage{
			ErrorType: cm.ErrResponseType(cm.SUCCESS_STANDARD_SUCCESS_TYPE),
			Error:     nil,
		}
	}

	// check start end input
	start := time.UnixMilli(int64(query.StartPeriod))
	if start.IsZero() {
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_BAD_REQUEST_TYPE,
			Error:     errors.New("start period is not valid"),
		}
	}
	end := time.UnixMilli(int64(query.EndPeriod))
	if end.IsZero() {
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_BAD_REQUEST_TYPE,
			Error:     errors.New("end period is not valid"),
		}
	}
	// build grpc model
	queryPaylod := pb.CorpQuery{
		QuizType:           query.QuizType,
		StartPeriod:        uint64(start.UnixMilli()),
		EndPeriod:          uint64(end.UnixMilli()),
		CorporateId:        query.CorporateId.String(),
		Users:              publicId,
		SubscriptionDetail: &subs,
	}
	// fetch grpc
	ctx_ = c.commonSvc.InsertCorrelationIdFromGrpc(ctx_)
	result, err := c.quizAssessCln.GetClient().GetCorporateAssessmentSummary(ctx_, &queryPaylod)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error calling grpc to quiz assess:%v", err.Error())
		return cm.ErrorMessage{
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
			Error:     err,
		}
	}
	if err = c.commonSvc.ObjectMapper(&result, &assessment); err != nil {
		c.sugar.WithContext(ctx).Errorf("error mapping object result:%v", err.Error())
	}
	assessment.Meta.StartDate = start
	assessment.Meta.EndDate = end
	return errMsg
}

func (c *CorporateEmployeeServiceImpl) GetCorpEmployeeAssessmentStatus(ctx context.Context, assessStatus *model.AssessmentStatus) (err error) {
	// get corporate status first
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err := errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)

	var status model.CorporateStatus
	if err = c.GetCorporateEmployeeStatusService(ctx, corporateID, &status); err != nil {
		return err
	}

	// fetch total in database
	if err = c.corpEmplRepo.GetEmployeeAssessStatus(ctx, status.SubscriptionDetail.ID, assessStatus); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching from database")
		return err
	}
	// merge the data
	assessStatus.PreAssessment = fmt.Sprintf("%v/%v",
		assessStatus.PreAssessmentAmount,
		status.EmployeeDetail.EmployeeAmount)
	assessStatus.PostAssessment = fmt.Sprintf("%v/%v",
		assessStatus.PostAssessmentAmount,
		status.EmployeeDetail.EmployeeAmount)
	assessStatus.FollowupAssessment = fmt.Sprintf("%v/%v",
		assessStatus.FollowupAssessmentAmount,
		status.EmployeeDetail.EmployeeAmount)

	return err
}

func (c *CorporateEmployeeServiceImpl) GetWellBeingAssessmentStatus(ctx context.Context, wellBeingStatus *model.WellBeingStatus) (errMsg cm.ErrorMessage) {
	// get corporate status first
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err := errors.New("error to get corporate id context")
		return cm.ErrorMessage{
			Error:     err,
			ErrorType: cm.ERR_STANDARD_BAD_REQUEST_TYPE,
		}
	}
	corporateID := tempID.(uuid.UUID)

	var status model.CorporateStatus
	if err := c.GetCorporateEmployeeStatusService(ctx, corporateID, &status); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching GetCorporateEmployeeStatusService with error:%v", err)
		return cm.ErrorMessage{
			Error:     err,
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
		}
	}

	if wellBeingStatus.Department == "" {
		wellBeingStatus.Department = "ALL"
	}

	// fetch total in database
	c.sugar.WithContext(ctx).Infof("getting well being assessment for :%v with department %v subs id :%v",
		status.SubscriptionDetail.CorporateID, wellBeingStatus.Department, status.SubscriptionDetail.ID)
	if err := c.corpEmplRepo.GetEmployeeWellBeingAssessStatus(ctx, status.SubscriptionDetail.ID, wellBeingStatus); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching from database")
		return cm.ErrorMessage{
			Error:     err,
			ErrorType: cm.ERR_STANDARD_INTERNAL_TYPE,
		}
	}
	return errMsg
}

func (c *CorporateEmployeeServiceImpl) UpdateCorpEmployeeByCorpAndSubsId(ctx context.Context, subsEntity entity.CorporateSubscriptionEntity, status string) (err error) {
	c.sugar.WithContext(ctx).Infof("getting employee for %v and subs id %v", subsEntity.CorporateID, subsEntity.ID)
	if c.assert.IsUUIDEmpty(subsEntity.CorporateID.String()) || c.assert.IsUUIDEmpty(subsEntity.ID.String()) {
		err = errors.New("uuid is invalid")
		c.sugar.WithContext(ctx).Errorf("error checking uuid for %v and subs id %v", subsEntity.CorporateID, subsEntity.ID)
		return err
	}

	// get from repo
	var empls []entity.CorporateEmployeeEntity
	if err = c.corpEmplRepo.GetEmployeeByCorpAndSubsId(ctx, subsEntity.CorporateID, subsEntity.ID, &empls); err != nil {
		c.sugar.WithContext(ctx).Errorf("err fetching corporate employee:%v", err.Error())
		return err
	}
	// generate channel
	emplChan := c.GenerateEmployeeChanByCorpAndSubsId(ctx, empls)
	// update concurrent
	ctxBack := c.commonSvc.ContextBackground(ctx)
	switch status {
	case "ACTIVE":
		go c.userSubsSvc.ActivateConcurrentUserFromCorpEmpl(ctxBack, emplChan, &subsEntity)
	default:
		go c.userSubsSvc.DeactivateConcurrentUserFromCorpEmpl(ctxBack, emplChan, &subsEntity, status)
	}

	return err
}

func (c *CorporateEmployeeServiceImpl) GenerateEmployeeChanByCorpAndSubsId(ctx context.Context, employees []entity.CorporateEmployeeEntity) <-chan entity.CorporateEmployeeEntity {
	employeeChan := make(chan entity.CorporateEmployeeEntity)
	go func() {
		for _, v := range employees {
			employeeChan <- v
		}
		close(employeeChan)
	}()
	return employeeChan
}
