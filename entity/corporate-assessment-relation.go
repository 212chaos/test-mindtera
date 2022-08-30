package entity

import (
	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
	cq "github.com/mindtera/quiz-assessment-service/entity"
)

type CorporateAssessmentRelation struct {
	ID               uuid.UUID `gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	CorporateValueID uuid.UUID `json:"corporate_value_id" gorm:"type:uuid;uniqueIndex:corporate_assessment_relation"`
	AssessmentID     uuid.UUID `json:"assessment_id" gorm:"type:uuid;uniqueIndex:corporate_assessment_relation"`
	ItemCode         string    `json:"item_code"`
	ParentSMM        string    `json:"parent_smm"`
	c.DefaultColumn
	RecordFlag string `json:"record_flag" gorm:"column:record_flag;index;default:'ACTIVE'"`

	CorporateValueDetail      *CorporateValueEntity `json:"corporate_value_detail,omitempty" gorm:"references:CorporateValueID;foreignKey:ID"`
	CorporateAssessmentDetail *cq.QuizDetailEntity  `json:"corporate_assessment_detail,omitempty" gorm:"references:AssessmentID;foreignKey:ID"`
}
