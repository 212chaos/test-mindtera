package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"

	commonentity "github.com/mindtera/go-common-module/common/entity"
)

type CorporateServiceBookingsGormable struct {
	ID          uuid.UUID `gorm:"column:corporate_service_bookings.id"`
	ServiceID   string    `gorm:"column:corporate_service_bookings.service_id"`
	CorporateID string    `gorm:"column:corporate_service_bookings.corporate_id"`
	ScheduledAt time.Time `gorm:"column:corporate_service_bookings.scheduled_at"`
	Status      string    `gorm:"column:corporate_service_bookings.status"`

	CreatedAt  time.Time `gorm:"column:corporate_service_bookings.created_at"`
	UpdatedAt  time.Time `gorm:"column:corporate_service_bookings.updated_at"`
	CreatedBy  string    `gorm:"column:corporate_service_bookings.created_by"`
	UpdatedBy  string    `gorm:"column:corporate_service_bookings.updated_by"`
	RecordFlag string    `gorm:"column:corporate_service_bookings.record_flag"`
}

func (d CorporateServiceBookingsGormable) ToCorporateServiceBookings() entity.CorporateServiceBooking {
	return entity.CorporateServiceBooking{
		ID:          d.ID,
		ServiceID:   d.ServiceID,
		CorporateID: d.CorporateID,
		ScheduledAt: d.ScheduledAt,
		Status:      entity.CorporateServiceBookingStatus(d.Status),
		DefaultColumn: commonentity.DefaultColumn{
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
			RecordFlag: d.RecordFlag,
		},
		CreatedAt: d.CreatedAt,
	}
}

type CorporateServiceBookingsGormable2 struct {
	CorporateServiceBookingsGormable
	MCorporateServiceGormable
}

func (d CorporateServiceBookingsGormable2) ToCorporateServiceBookings() entity.CorporateServiceBooking {
	corporateServiceBooking := d.CorporateServiceBookingsGormable.ToCorporateServiceBookings()
	corporateServiceBooking.MCorporateService = &entity.MCorporateService{
		ID:              d.MCorporateServiceGormable.ID,
		Name:            d.MCorporateServiceGormable.Name,
		Description:     d.MCorporateServiceGormable.Description,
		LongDescription: d.MCorporateServiceGormable.LongDescription,
		LandingURL:      d.MCorporateServiceGormable.LandingURL,
		ImageURLSource:  d.MCorporateServiceGormable.ImageURLSource,
		ServiceType:     entity.CorporateServiceType(d.MCorporateServiceGormable.ServiceType),
		DefaultColumn: commonentity.DefaultColumn{
			CreatedAt:  d.MCorporateServiceGormable.CreatedAt,
			UpdatedAt:  d.MCorporateServiceGormable.UpdatedAt,
			CreatedBy:  d.MCorporateServiceGormable.CreatedBy,
			UpdatedBy:  d.MCorporateServiceGormable.UpdatedBy,
			RecordFlag: d.MCorporateServiceGormable.RecordFlag,
		},
	}
	return corporateServiceBooking
}

type CorporateServiceBookingJson struct {
	ServiceID   string `json:"service_id"`
	ScheduledAt string `json:"scheduled_at"`
}

func (d CorporateServiceBookingJson) GenerateCreateModel() (CorporateServiceBookingCreation, error) {
	scheduledAt, err := time.Parse(time.RFC3339, d.ScheduledAt)
	if err != nil {
		return CorporateServiceBookingCreation{}, err
	}

	return CorporateServiceBookingCreation{
		ScheduledAt: scheduledAt,
		ServiceID:   d.ServiceID,
	}, nil
}

type CorporateServiceBookingCreation struct {
	ServiceID   string
	ScheduledAt time.Time
}
