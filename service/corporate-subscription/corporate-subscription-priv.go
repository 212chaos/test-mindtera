package service_corporatesubscription

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	schedulerentity "github.com/mindtera/scheduler-service/entity"
)

func (c *CorporateSubscriptionServiceImpl) generateSchedulerInitialPayload(corporateSubscription entity.CorporateSubscriptionEntity) ([]schedulerentity.SchedulerPayload, error) {
	// make initial payload
	defaultCorporateConditions := map[string]string{
		"POST_ASSESSMENT":     common.CORPORATE_SERVICE_CALLBACK_URL + common.API_URL + "/scheduler/assessment-callback",
		"FOLLOWUP_ASSESSMENT": common.CORPORATE_SERVICE_CALLBACK_URL + common.API_URL + "/scheduler/assessment-callback",
		"END_PERIOD":          common.CORPORATE_SERVICE_CALLBACK_URL + common.API_URL + "/scheduler/corporate-end-period",
	}

	// create template for payload
	payloadTemplate := schedulerentity.SchedulerPayload{
		Authorization: common.STATIC_AUTH,
		Issuer:        common.SERVICE_NAME,
	}

	// create template message
	messageTemplate := model.CorporateSubscriptionScheduler{
		CorporateID:    corporateSubscription.CorporateID,
		SubscriptionID: corporateSubscription.ID,
	}

	// loop over default conditions
	var schedulerEntity []schedulerentity.SchedulerPayload
	for k, v := range defaultCorporateConditions {
		// generate message payload
		message, err := c.genereateMessageWithCorporatePeriod(messageTemplate, k)
		if err != nil {
			return nil, err
		}
		payloadTemplate.Message = message
		payloadTemplate.CallbackUrl = v

		// switch period based on time
		switch k {
		case "POST_ASSESSMENT":
			payloadTemplate.TimeMs = int(corporateSubscription.StartPeriod.AddDate(0, 2, 0).UnixMilli())
			if strings.EqualFold(common.ENV, "DEVELOPMENT") {
				payloadTemplate.TimeMs = int(time.Now().Add(time.Minute * 10).UnixMilli())
			}
		case "FOLLOWUP_ASSESSMENT":
			payloadTemplate.TimeMs = int(corporateSubscription.StartPeriod.AddDate(0, 5, 0).UnixMilli())
			if strings.EqualFold(common.ENV, "DEVELOPMENT") {
				payloadTemplate.TimeMs = int(time.Now().Add(time.Minute * 15).UnixMilli())
			}
		default:
			payloadTemplate.TimeMs = int(corporateSubscription.EndPeriod.UnixMilli())
		}
		// append scheduler payload
		schedulerEntity = append(schedulerEntity, payloadTemplate)
	}

	return schedulerEntity, nil
}

func (c *CorporateSubscriptionServiceImpl) genereateMessageWithCorporatePeriod(messageTemplate model.CorporateSubscriptionScheduler, period string) (string, error) {
	// generate message
	messageTemplate.CorporatePeriod = period
	b, err := json.Marshal(messageTemplate)
	if err != nil {
		return "", err
	}
	message := string(b)
	return message, nil
}
