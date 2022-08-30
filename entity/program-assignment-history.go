package entity

import (
	"time"

	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
	commonmodel "github.com/mindtera/go-common-module/fajar/model"
)

type EmployeeProgramAssignmentHistory struct {
	ID            uuid.UUID         `json:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" `
	ProgramID     uuid.UUID         `json:"prorgram_id"`
	CorporateID   uuid.UUID         `json:"corporate_id"`
	RelatedID     uuid.UUID         `json:"related_id"`
	AssigningType TypeAssigningMode `json:"assigning_type"`

	c.DefaultColumn
}

func (EmployeeProgramAssignmentHistory) TableName() string {
	return "employee_program_assignment_history"
}

type EmployeeProgramAssignmentHistoryRelatedData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EmployeeProgramAssignmentHistory2
// includes relatedID data based on the following:
// CORPORATE = All active employee in the company
// DEPARTMENT = All active employee in the department
// INDIVIDUAL = an employee in the department
type EmployeeProgramAssignmentHistory2 struct {
	ID            uuid.UUID                                   `json:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" `
	ProgramID     uuid.UUID                                   `json:"prorgram_id"`
	MProgram      commonmodel.MProgram                        `json:"m_program"`
	CorporateID   uuid.UUID                                   `json:"corporate_id"`
	RelatedID     uuid.UUID                                   `json:"related_id"`
	RelatedData   EmployeeProgramAssignmentHistoryRelatedData `json:"related_data"`
	AssigningType TypeAssigningMode                           `json:"assigning_type"`

	c.DefaultColumn
	CreatedAt time.Time `json:"created_at"`
}
