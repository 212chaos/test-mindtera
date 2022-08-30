package entity

import (
	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateServiceType string

const (
	EXERCISE_TYPE CorporateServiceType = "EXERCISE"
	WEBINAR_TYPE  CorporateServiceType = "WEBINAR"
	COACHING_TYPE CorporateServiceType = "COACHING"
)

type MCorporateService struct {
	ID              uuid.UUID            `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	LongDescription string               `json:"long_description"`
	LandingURL      string               `json:"landing_url"`
	ImageURLSource  string               `json:"-"`
	ImageURL        string               `json:"image_url"`
	ServiceType     CorporateServiceType `json:"service_type"`

	c.DefaultColumn
}
