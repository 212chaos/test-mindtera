package model

import (
	"time"

	"github.com/google/uuid"
	consumermodel "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/entity"
	c "github.com/mindtera/go-common-module/common/entity"
	commonmodel "github.com/mindtera/go-common-module/fajar/model"
	"gorm.io/gorm"
)

type EmployeeProgramAssignmentHistoryGormable struct {
	ID            uuid.UUID `gorm:"column:employee_program_assignment_history.id"`
	ProgramID     uuid.UUID `gorm:"column:employee_program_assignment_history.program_id"`
	CorporateID   uuid.UUID `gorm:"column:employee_program_assignment_history.corporate_id"`
	RelatedID     uuid.UUID `gorm:"column:employee_program_assignment_history.related_id"`
	AssigningType string    `gorm:"column:employee_program_assignment_history.assigning_type"`
	CreatedAt     time.Time `gorm:"column:employee_program_assignment_history.created_at"`
	UpdatedAt     time.Time `gorm:"column:employee_program_assignment_history.updated_at"`
	DeletedAt     time.Time `gorm:"column:employee_program_assignment_history.deleted_at"`
	CreatedBy     string    `gorm:"column:employee_program_assignment_history.created_by"`
	UpdatedBy     string    `gorm:"column:employee_program_assignment_history.updated_by"`
	RecordFlag    string    `gorm:"column:employee_program_assignment_history.record_flag"`
}

// func (d *EmployeeProgramAssignmentHistoryGormable) ToEmployeeProgramAssignmentHistory2() entity.EmployeeProgramAssignmentHistory {
// 	return

// }

type EmployeeProgramAssignmentHistoryGormable2 struct {
	EmployeeProgramAssignmentHistoryGormable
	consumermodel.MProgramGormable3
}

func (d *EmployeeProgramAssignmentHistoryGormable2) ToEmployeeProgramAssignmentHistory2() entity.EmployeeProgramAssignmentHistory2 {
	output := entity.EmployeeProgramAssignmentHistory2{
		ID:            d.EmployeeProgramAssignmentHistoryGormable.ID,
		ProgramID:     d.EmployeeProgramAssignmentHistoryGormable.ProgramID,
		CorporateID:   d.EmployeeProgramAssignmentHistoryGormable.CorporateID,
		RelatedID:     d.EmployeeProgramAssignmentHistoryGormable.RelatedID,
		AssigningType: entity.TypeAssigningMode(d.EmployeeProgramAssignmentHistoryGormable.AssigningType),
		DefaultColumn: c.DefaultColumn{
			CreatedAt: d.EmployeeProgramAssignmentHistoryGormable.CreatedAt,
			UpdatedAt: d.EmployeeProgramAssignmentHistoryGormable.UpdatedAt,
			DeletedAt: gorm.DeletedAt{
				Time:  d.EmployeeProgramAssignmentHistoryGormable.DeletedAt,
				Valid: true,
			},
			CreatedBy:  d.EmployeeProgramAssignmentHistoryGormable.CreatedBy,
			UpdatedBy:  d.EmployeeProgramAssignmentHistoryGormable.UpdatedBy,
			RecordFlag: d.EmployeeProgramAssignmentHistoryGormable.RecordFlag,
		},
		CreatedAt: d.EmployeeProgramAssignmentHistoryGormable.CreatedAt,
		MProgram: commonmodel.MProgram{
			ID:          d.MProgramGormable3.ID,
			Name:        d.MProgramGormable3.Name,
			ContentType: commonmodel.MProgramContentType(d.MProgramGormable3.ContentType),
			Category:    commonmodel.MProgramCat(d.MProgramGormable3.Category),
			SubCategory: commonmodel.MProgramSubCat(d.MProgramGormable3.SubCategory),
			Description: d.MProgramGormable3.Description,
			// Objective: (*commonmodel.PGTextArr)(d.MProgramGormable3.Objective),
			ThumbnailSource:  d.MProgramGormable3.ThumbnailSource,
			CoachID:          d.MProgramGormable3.CoachID,
			IntroVideoID:     d.MProgramGormable3.IntroVideoID,
			Rating:           d.MProgramGormable3.Rating,
			Duration:         d.MProgramGormable3.Duration,
			SubscriptionType: commonmodel.TSubscription(d.MProgramGormable3.SubscriptionType),
			Status:           commonmodel.TProgramStatus(d.MProgramGormable3.Status),
			DefaultColumns: commonmodel.DefaultColumns{
				CreatedAt:  d.MProgramGormable3.CreatedAt,
				UpdatedAt:  d.MProgramGormable3.UpdatedAt,
				CreatedBy:  d.MProgramGormable3.CreatedBy,
				UpdatedBy:  d.MProgramGormable3.UpdatedBy,
				RecordFlag: commonmodel.TRecordFlag(d.MProgramGormable3.RecordFlag),
			},
		},
	}
	return output

}
