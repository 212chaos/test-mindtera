package entity

import (
	"time"

	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateServiceBookingStatus string

const (
	WAITING_CONFIRM_STATUS CorporateServiceBookingStatus = "WAITING_CONFIRMATION"
	CONFIRMED_STATUS       CorporateServiceBookingStatus = "CONFIRMED"
	FINISHED_STATUS        CorporateServiceBookingStatus = "FINISHED"
	CANCELED               CorporateServiceBookingStatus = "CANCELED"
)

type CorporateServiceBooking struct {
	ID                uuid.UUID                     `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	ServiceID         string                        `json:"service_id"`
	CorporateID       string                        `json:"corporate_id"`
	ScheduledAt       time.Time                     `json:"scheduled_at"`
	Status            CorporateServiceBookingStatus `json:"status"`
	MCorporateService *MCorporateService            `json:"m_corporate_service"`
	c.DefaultColumn
	CreatedAt time.Time `json:"created_at"`
}

func (CorporateServiceBooking) TableName() string {
	return "corporate_service_bookings"
}
