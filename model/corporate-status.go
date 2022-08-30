package model

import (
	"time"

	"github.com/mindtera/corporate-service/entity"
)

type EmployeeDetail struct {
	EmployeeAmount int64 `json:"employee_amount"`
	EmployeeLimit  int64 `json:"employee_limit"`
}

type CorporateStatus struct {
	EmployeeDetail     EmployeeDetail                     `json:"employee_detail"`
	SubscriptionDetail entity.CorporateSubscriptionEntity `json:"subscription_detail"`
}

type CorporateInformation struct {
	EmployeeInformation entity.CorporateEmployeeEntity     `json:"employee_information"`
	SubscriptionDetail  entity.CorporateSubscriptionEntity `json:"subscription_detail"`
}

type AssessmentStatus struct {
	PreAssessment            string `json:"pre_assessment" gorm:"-"`
	PostAssessment           string `json:"post_assessment" gorm:"-"`
	FollowupAssessment       string `json:"followup_assessment" gorm:"-"`
	PreAssessmentAmount      int    `json:"-" gorm:"pre_assessment_amount"`
	PostAssessmentAmount     int    `json:"-" gorm:"post_assessment_amount"`
	FollowupAssessmentAmount int    `json:"-" gorm:"followup_assessment_amount"`
}

type WellBeingStatus struct {
	Status     string    `json:"status"`
	Department string    `json:"department_code"`
	EndDate    time.Time `json:"end_time"`
	StartDate  time.Time `json:"start_time"`
}
