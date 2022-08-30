package entity

import (
	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
	"github.com/mindtera/go-common-module/common/v2/model"
)

type Feedback struct {
	ID                uuid.UUID                `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserId            uuid.UUID                `json:"user_id" gorm:"type:uuid"`
	CorporateId       uuid.UUID                `json:"corporate_id" gorm:"index;type:uuid"`
	Body              string                   `json:"body"`
	Sender            string                   `json:"sender"`
	CategoryCode      string                   `json:"category_code"`
	ShownCategoryCode string                   `json:"shown_category_code"`
	ReadFlag          bool                     `json:"read_flag" gorm:"index;default:false"`
	RecordFlag        string                   `json:"record_flag,omitempty" gorm:"index;default:'ACTIVE'"`
	EmployeeDetail    *CorporateEmployeeEntity `json:"employee_detail,omitempty" gorm:"references:UserId;foreignKey:PublicUserID"`
	c.DefaultColumn
}

func (Feedback) TableName() string {
	return "corporate_feedbacks"
}

type FeedbackCategory struct {
	ID               uuid.UUID          `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Code             string             `json:"code"`
	Name             string             `json:"name,omitempty"`
	EnglishName      string             `json:"english_name,omitempty"`
	LocalizationName model.Localization `json:"localization_name" gorm:"embedded;embeddedPrefix:localization_name_;->"`
	c.DefaultColumn
}

func (FeedbackCategory) TableName() string {
	return "corporate_feedback_categories"
}

type FeedbackShownCategory struct {
	ID                      uuid.UUID          `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Code                    string             `json:"code"`
	Name                    string             `json:"name,omitempty"`
	EnglishName             string             `json:"english_name,omitempty"`
	Description             string             `json:"description,omitempty"`
	EnglishDescription      string             `json:"english_description,omitempty"`
	LocalizationName        model.Localization `json:"localization_name" gorm:"embedded;embeddedPrefix:localization_name_;->"`
	LocalizationDescription model.Localization `json:"localization_description" gorm:"embedded;embeddedPrefix:localization_description_;->"`
	c.DefaultColumn
}

func (FeedbackShownCategory) TableName() string {
	return "corporate_feedback_shown_categories"
}
