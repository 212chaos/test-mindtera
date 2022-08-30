package service_corporateemployee

import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	publicuser "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	mailjet "github.com/mindtera/go-common-module/common/client/mailjet-client"
	"gorm.io/gorm"
)

// conditional function to get corporate period type
func (c *CorporateEmployeeServiceImpl) getCorporatePeriodType(ctx context.Context, subsDetail entity.CorporateSubscriptionEntity) (subscriptionID uuid.UUID, currentCondition string, dueDate time.Time, err error) {
	subscriptionID = subsDetail.ID
	// get corporate period condition
	currentCondition, dueDate, err = c.util.GetCorporateSubsPeriod(ctx, subsDetail.StartPeriod, subsDetail.EndPeriod)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error to get corporate subs period:%v", err)
	}
	return subscriptionID, currentCondition, dueDate, err
}

// function to generate task template for employee to assign employee to task
func (c *CorporateEmployeeServiceImpl) generateTaskTemplateForEmployee(ctx context.Context, task string, dueDate time.Time, subsDetail entity.CorporateSubscriptionEntity) (template entity.EmployeeTaskEntity, err error) {
	// assigned time
	timeNow := time.Now()
	// generate detail from subscription detail
	subsID, currentCondition, dDate, err := c.getCorporatePeriodType(ctx, subsDetail)
	if err != nil {
		return template, err
	}

	// checking
	if strings.EqualFold(task, "CORPORATE_ASSESSMENT") {
		task = currentCondition
		// check current condition
		switch task {
		case string(entity.TYPE_PRE_ASSESSMENT):
			dDate = timeNow.AddDate(0, 2, 0)
		case string(entity.TYPE_POST_ASSESSMENT):
			dDate = timeNow.AddDate(0, 3, 0)
		}

		if strings.EqualFold(common.ENV, "DEVELOPMENT") || strings.EqualFold(common.ENV, "DEV") {
			switch task {
			case string(entity.TYPE_PRE_ASSESSMENT):
				dDate = timeNow.Add(10 * time.Minute)
			case string(entity.TYPE_POST_ASSESSMENT):
				dDate = timeNow.Add(15 * time.Minute)
			case string(entity.TYPE_FOLLOWUP_ASSESSMENT):
				dDate = timeNow.Add(20 * time.Minute)
			}
		}
		dueDate = dDate
	}

	if strings.EqualFold(currentCondition, "INVALID") {
		err = errors.New("subscription is not valid")
		return template, err
	}

	template = entity.EmployeeTaskEntity{
		AssignDate:     &timeNow,
		DueDate:        &dueDate,
		SubscriptionID: subsID,
		TaskType:       entity.EMPLOYEE_TASK_TYPE(task),
		TaskStatus:     entity.STATUS_ACTIVE,
	}
	return template, err
}

func (c *CorporateEmployeeServiceImpl) validateEmployeePayload(ctx context.Context, employeeHash, roleHash map[string]bool, v *entity.CorporateEmployeeEntity) (err *model.CorporateEmployeeError) {
	// check required payload
	if c.assert.IsEmpty(v.Email) {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "email is empty",
		}
		return err
	}
	// check email valid or not
	_, e := mail.ParseAddress(v.Email)
	if e != nil {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "email format is invalid",
		}
		return err
	}

	if c.assert.IsEmpty(v.Department) {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "department is empty",
		}
		return err
	}
	if c.assert.IsEmpty(v.Role) {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "role is empty",
		}
		return err
	}
	if c.assert.IsEmpty(v.EmployeeName) {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "employee name is empty",
		}
		return err
	}

	v.Role = strings.ReplaceAll(strings.ToUpper(v.Role), " ", "_")
	v.Email = strings.ToLower(v.Email)
	// check role
	if !roleHash[v.Role] {
		err = &model.CorporateEmployeeError{
			Email:        v.Email,
			Number:       v.Number,
			ErrorMessage: "role is invalid",
		}
		return err
	}

	// check email exist or not
	if c.assert.IsUUIDEmpty(v.ID.String()) {
		if employeeHash[v.Email] && c.assert.IsUUIDEmpty(v.ID.String()) {
			err = &model.CorporateEmployeeError{
				Email:        v.Email,
				Number:       v.Number,
				ErrorMessage: "email is duplicated",
			}
			return err
		}
	}

	if !c.assert.IsUUIDEmpty(v.ID.String()) {
		c.redisSvc.Delete(context.Background(), "CORPORATE_PUBLIC_EMPLOYEE_"+v.ID.String()+"*")
	}

	// update flag
	if v.RecordFlag == "ACTIVE" {
		if c.assert.IsUUIDEmpty(v.PublicUserID.String()) {
			v.RecordFlag = "UNREGISTERED"
		}
		c.sugar.WithContext(ctx).Infof("changed status to:%v", v.RecordFlag)
		v.DeletedAt = gorm.DeletedAt{
			Valid: false,
			Time:  time.Now()}
	}
	return nil
}

// generate public user id queries from corporate employee correct format
func (c *CorporateEmployeeServiceImpl) generateUserFormatFromEmail(ctx context.Context, corporateEmployees *[]entity.CorporateEmployeeEntity) (map[string]bool, map[string]uuid.UUID, error) {
	// init variable
	userArrHash := map[string]bool{}
	userIDHash := map[string]uuid.UUID{}

	// build array of email
	var emails []string
	for _, val := range *corporateEmployees {
		emails = append(emails, val.Email)
	}
	// fetching user data from databae
	var userArr []publicuser.User
	if err := c.corpEmplRepo.GetUserByUserEmail(ctx, emails, &userArr); err != nil {
		c.sugar.WithContext(ctx).Errorf("error fetching public user data:%v", err)
		err = errors.New("error fetching data from payload")
		return userArrHash, userIDHash, err
	}

	for _, val := range userArr {
		userArrHash[val.Email] = true
		userIDHash[val.Email] = val.ID
	}

	return userArrHash, userIDHash, nil
}

func (c *CorporateEmployeeServiceImpl) removeWhiteSpaces(text string) string {
	arr := strings.Split(text, " ")

	// init all variables needed
	var tempArr []string

	// loop over the strings
	for _, v := range arr {
		if !c.assert.IsEmpty(v) && !(v == "\t") && !(v == "\r") {
			tempArr = append(tempArr, v)
		}
	}
	return strings.Join(tempArr, " ")
}

// method to sending registration reminder email
func (c *CorporateEmployeeServiceImpl) sendingRegistrationEmail(ctx context.Context,
	corporateEmployee entity.CorporateEmployeeEntity) {
	payload := struct {
		mailjet.MailjetBaseTemplate
		Username      string `json:"user_name"`
		Corporate     string `json:"corporate"`
		EmployeeEmail string `json:"employee_email"`
	}{}
	// declaring value
	payload.Subject = common.REGISTRATION_EMAIL_SUBJECT
	payload.TemplateID = common.REGISTRATION_EMAIL_TEMPLATE
	payload.MailToEmail = corporateEmployee.Email
	payload.MailToName = corporateEmployee.EmployeeName
	payload.Username = corporateEmployee.EmployeeName
	payload.EmployeeEmail = corporateEmployee.Email
	payload.Corporate = corporateEmployee.Corporate

	// creating hashmap
	var mp map[string]interface{}
	b, err := json.Marshal(&payload)
	if err != nil {
		c.sugar.WithContext(ctx).Errorf("error when marshaling payload:%v", err.Error())
		return
	}
	if err := json.Unmarshal(b, &mp); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when unmarshaling payload:%v", err.Error())
		return
	}
	// sending email
	c.sugar.WithContext(ctx).Infof("sending email reminder for:%v", payload.MailToEmail)
	if err := c.mailjet.SendEmailExternal(ctx, mp); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when sending payload:%v", err.Error())
		return
	}
}

func (c *CorporateEmployeeServiceImpl) validateRegirstrationEmail(ctx context.Context, corporateEmployees []entity.CorporateEmployeeEntity) {
	// check public key
	isPublic := ctx.Value(common.PUBLIC_KEY)
	if isPublic == nil {
		for _, v := range corporateEmployees {
			if (c.assert.IsUUIDEmpty(v.PublicUserID.String()) && strings.EqualFold(v.RecordFlag, "UNREGISTERED")) ||
				(!c.assert.IsUUIDEmpty(v.PublicUserID.String()) && strings.EqualFold(v.RecordFlag, "ACTIVE")) {
				go c.sendingRegistrationEmail(ctx, v)
			}
		}
	}
}

func (c *CorporateEmployeeServiceImpl) generateFirstTasks(ctx context.Context, correctFormat []entity.CorporateEmployeeEntity, subsDetail entity.CorporateSubscriptionEntity) (err error) {
	// validate correct format with id
	var generatedCorrectFormat []entity.CorporateEmployeeEntity
	publicKey := ctx.Value(common.PUBLIC_KEY)
	for _, v := range correctFormat {
		if c.assert.IsUUIDEmpty(v.ID.String()) || (publicKey != nil) {
			generatedCorrectFormat = append(generatedCorrectFormat, v)
		}
	}

	if len(generatedCorrectFormat) > 0 {
		// generate corporate assessment task
		c.sugar.WithContext(ctx).Info("generating corp assessment employee task template process")
		taskTemplate, err := c.generateTaskTemplateForEmployee(ctx, "CORPORATE_ASSESSMENT", time.Time{}, subsDetail)
		if err != nil {
			c.sugar.WithContext(ctx).Errorf("error generate template task for employee:%v", err)
			return err
		}
		c.sugar.WithContext(ctx).Info("upserting employee task table")
		if err = c.employeeTaskSvc.UpsertEmployeeTaskByEmployeeEntityService(ctx, &generatedCorrectFormat, taskTemplate); err != nil {
			c.sugar.WithContext(ctx).Errorf("error upserting corp assessment employee task :( :%v", err)
			return err
		}

		// generate well being task
		c.sugar.WithContext(ctx).Info("generating well being employee task template process")
		_, endMonth := c.commonSvc.GenerateFirstEndMonth(time.Now())
		taskTemplate, err = c.generateTaskTemplateForEmployee(ctx, "WELL_BEING", endMonth, subsDetail)
		if err != nil {
			c.sugar.WithContext(ctx).Errorf("error generate template task for employee:%v", err)
			return err
		}
		c.sugar.WithContext(ctx).Info("upserting employee task table")
		if err = c.employeeTaskSvc.UpsertEmployeeTaskByEmployeeEntityService(ctx, &generatedCorrectFormat, taskTemplate); err != nil {
			c.sugar.WithContext(ctx).Errorf("error upserting well being employee task :( :%v", err)
			return err
		}
	}
	return err
}
