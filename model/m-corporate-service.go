package model

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mindtera/corporate-service/entity"
	commonentity "github.com/mindtera/go-common-module/common/entity"
	commonmodel "github.com/mindtera/go-common-module/common/model"
)

type MCorporateServiceGormable struct {
	ID              uuid.UUID `gorm:"column:m_corporate_services.id"`
	Name            string    `gorm:"column:m_corporate_services.name"`
	Description     string    `gorm:"column:m_corporate_services.description"`
	LongDescription string    `gorm:"column:m_corporate_services.long_description"`
	LandingURL      string    `gorm:"column:m_corporate_services.landing_url"`
	ImageURLSource  string    `gorm:"column:m_corporate_services.image_url"`
	ServiceType     string    `gorm:"column:m_corporate_services.service_type"`

	CreatedAt  time.Time `gorm:"column:m_corporate_services.created_at"`
	UpdatedAt  time.Time `gorm:"column:m_corporate_services.updated_at"`
	CreatedBy  string    `gorm:"column:m_corporate_services.created_by"`
	UpdatedBy  string    `gorm:"column:m_corporate_services.updated_by"`
	RecordFlag string    `gorm:"column:m_corporate_services.record_flag"`
	DeletedAt  time.Time `gorm:"column:m_corporate_services.deleted_at"`
}

func (d *MCorporateServiceGormable) ToMCorporateService() entity.MCorporateService {
	return entity.MCorporateService{
		ID:              d.ID,
		Name:            d.Name,
		Description:     d.Description,
		LongDescription: d.LongDescription,
		LandingURL:      d.LandingURL,
		ImageURLSource:  d.ImageURLSource,
		ServiceType:     entity.CorporateServiceType(d.ServiceType),
		DefaultColumn: commonentity.DefaultColumn{
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
			// DeletedAt: d.DeletedAt,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
			RecordFlag: d.RecordFlag,
		},
	}
}

type MCorporateServiceRequest struct {
	commonmodel.PaginationModel
	IncludeWebinar  bool `form:"include_webinar"`
	IncludeExercise bool `form:"include_exercise"`
	IncludeCoaching bool `form:"include_coaching"`
}

// search program request payload
type SearchCorporateServiceRequest struct {
	Query           string `form:"query"`
	IncludeWebinar  bool   `form:"include_webinar"`
	IncludeExercise bool   `form:"include_exercise"`
	IncludeCoaching bool   `form:"include_coaching"`
}

func (d *SearchCorporateServiceRequest) ValidateRequest() bool {

	if d.Query == "" {
		return false
	}

	d.Query = strings.ToUpper(d.Query)

	return true
}
