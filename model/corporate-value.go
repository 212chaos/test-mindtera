package model

import (
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
)

type CorporateValueModel struct {
	TotalAssessment       int                           `json:"total_assessment,omitempty"`
	ParentSMM             string                        `json:"parent_smm"`
	SMMDetail             entity.CorporateSMM           `json:"smm_detail"`
	CorporateValueRecords []entity.CorporateValueEntity `json:"corporate_value_records"`
}

type CorporateRelationAssessmentModel struct {
	ParentSMM        string                   `json:"parent_smm"`
	Name             string                   `json:"name,omitempty"`
	EnglishName      string                   `json:"english_name,omitempty"`
	Code             string                   `json:"code,omitempty"`
	AssessmentID     uuid.UUID                `json:"assessment_id" gorm:"type:uuid;uniqueIndex:corporate_assessment_relation"`
	LocalizationName commonmodel.Localization `json:"localization_name" gorm:"embedded;embeddedPrefix:localization_name_;->"`
	entity.CorporateValueRelation
}

type CorporateValueRelationResultModel struct {
	TotalAssessment int                             `json:"total_assessment,omitempty"`
	RawData         []entity.CorporateValueRelation `json:"raw_data"`
	Detail          []CorporateValueModel           `json:"detail,omitempty"`
}

type CorpValueId struct {
	TaskAssessId uuid.UUID `json:"task_assess_id" gorm:"->;embedded"`
	CorpSmm      string    `json:"corp_smm" gorm:"->;embedded"`
	CorpValueId  uuid.UUID `json:"corp_value_id" gorm:"->;embedded"`
}

type CorpValue struct {
	TaskAssessId uuid.UUID                   `json:"task_assess_id" gorm:"->;embedded"`
	CorpSmm      string                      `json:"corp_smm" gorm:"->;embedded"`
	CorpValueId  uuid.UUID                   `json:"corp_value_id" gorm:"->;embedded"`
	Detail       entity.CorporateValueEntity `json:"detail" gorm:"->;embedded;embeddedPrefix:detail_"`
}

type CorpQuery struct {
	CorporateId uuid.UUID `json:"corporate_id,omitempty" form:"corporate_id"`
	QuizType    string    `json:"quiz_type,omitempty" form:"quiz_type"`
	StartPeriod uint64    `json:"start_period,omitempty" form:"start_period"`
	EndPeriod   uint64    `json:"end_period,omitempty" form:"end_period"`
	EndTime     int64     `json:"end_time" form:"end_date"`
	StartTime   int64     `json:"start_time" form:"start_date"`
	Department  string    `json:"department,omitempty" form:"department"`
	Users       string    `json:"users,omitempty" form:"users"`
}

type CorpValueName struct {
	CorporateValueId uuid.UUID                `json:"corporate_value_id" gorm:"->;column:corporate_value_id"`
	LocalizationName commonmodel.Localization `json:"localization_name" gorm:"->;embedded;embeddedPrefix:localization_name_"`
}
