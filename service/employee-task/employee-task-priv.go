package service_employeetask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	mailjetclient "github.com/mindtera/go-common-module/common/client/mailjet-client"
)

// private: employee task validation service
func (e *EmployeeTaskServiceImpl) validateEmployeeTaskPayload(ctx context.Context, employeeTask entity.EmployeeTaskEntity) (err error) {
	if e.assert.IsUUIDEmpty(employeeTask.UserID.String()) {
		e.sugar.WithContext(ctx).Errorf("error in payload uuid  is not valid :%v", employeeTask)
		err := errors.New("user uuid  is not valid")
		return err
	}

	if e.assert.IsEmpty(string(employeeTask.TaskType)) {
		e.sugar.WithContext(ctx).Errorf("error in payload task type is not valid:%v", employeeTask)
		err := errors.New("user task type is not valid")
		return err
	}
	return err
}

// private: get subscription ID from context
func (e *EmployeeTaskServiceImpl) getSubscriptionIDFromContext(ctx context.Context) (subsID uuid.UUID, err error) {
	// do query logic for subscription id
	// preparing corporate ID
	corporateID := ctx.Value("corporate_id")
	if corporateID == nil {
		err = errors.New("context is invalid")
		return subsID, err
	}
	corporateSubscription := entity.CorporateSubscriptionEntity{
		CorporateID: corporateID.(uuid.UUID)}
	if err = e.corporateSubsSvc.GetCorporateSubscriptionService(ctx, &corporateSubscription); err != nil {
		return subsID, err
	}
	subsID = corporateSubscription.ID
	return subsID, err
}

// sending email for assessment reminder
func (e *EmployeeTaskServiceImpl) sendingAssessmentEmail(ctx context.Context,
	task entity.EmployeeTaskEntity) {
	payload := struct {
		mailjetclient.MailjetBaseTemplate
		Username string `json:"user_name"`
	}{}
	// declaring value
	payload.Subject = common.ASSESSMENT_REMINDER_SUBJECT
	payload.TemplateID = common.ASSESSMENT_REMINDER_TEMPLATE_ID
	payload.MailToEmail = task.Email
	payload.MailToName = task.EmployeeName
	payload.Username = task.EmployeeName

	// creating hashmap
	var mp map[string]interface{}
	b, err := json.Marshal(&payload)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error when marshaling payload:%v", err.Error())
		return
	}
	if err := json.Unmarshal(b, &mp); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when unmarshaling payload:%v", err.Error())
		return
	}
	// sending email
	if err := e.mailjet.SendEmailExternal(ctx, mp); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when sending payload:%v", err.Error())
		return
	}
}

// loop over the payload
func (e *EmployeeTaskServiceImpl) validateTaskReminderEmail(ctx context.Context, tasks []entity.EmployeeTaskEntity) {
	for _, v := range tasks {
		e.sendingAssessmentEmail(ctx, v)
	}
}

func (e *EmployeeTaskServiceImpl) getGooglePublicObject(bucketName, object string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, object)
}

func (e *EmployeeTaskServiceImpl) sendBookingInfoEmail(ctx context.Context, serviceID string, scheduledAt time.Time) (err error) {
	tempID := ctx.Value("corporate_id")
	if tempID == nil {
		err = errors.New("error to get corporate id context")

		return err
	}
	corporateID := tempID.(uuid.UUID)
	type payload struct {
		mailjetclient.MailjetBaseTemplate
		Company     string `json:"company"`
		ProposeTime string `json:"propose_time"`
		ServiceType string `json:"service_type"`
		ServiceName string `json:"service_name"`
	}
	corpService, err := e.repository.GetCorporateServiceByID(ctx, serviceID)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error sendBookingInfoEmail:%v", err.Error())
		return errors.New("internal server error sendBookingInfoEmail")
	}
	corporateFetch, err := e.repository.GetCorporateDetails(ctx, corporateID.String())
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error sendBookingInfoEmail:%v", err.Error())
		return errors.New("internal server error sendBookingInfoEmail")
	}

	pl := payload{
		MailjetBaseTemplate: mailjetclient.MailjetBaseTemplate{
			Subject:     fmt.Sprintf("Mindtera Service - Booking from %s", corporateFetch.Name),
			MailToEmail: common.MAILJET_MAIL_TO_CORPORATE_SERVICE_BOOKING,
			MailToName:  "corporate",
			TemplateID:  common.ASSESSMENT_REMINDER_TEMPLATE_ID,
		},
		Company:     corporateFetch.Code,
		ProposeTime: scheduledAt.Format(time.RFC1123),
		ServiceType: string(corpService.ServiceType),
		ServiceName: corporateFetch.Name,
	}

	var mp map[string]interface{}
	b, err := json.Marshal(&pl)
	if err != nil {
		e.sugar.WithContext(ctx).Errorf("error when marshaling payload:%v", err.Error())
		return
	}
	if err = json.Unmarshal(b, &mp); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when unmarshaling payload:%v", err.Error())
		return
	}
	// sending email
	if err = e.mailjet.SendEmailExternal(ctx, mp); err != nil {
		e.sugar.WithContext(ctx).Errorf("error when sending payload:%v", err.Error())
		return
	}
	return nil
}
