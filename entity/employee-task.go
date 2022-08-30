package entity

import (
	"time"

	"github.com/google/uuid"
	ce "github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
)

type TypeAssigningMode string

const (
	ASSIGN_MODE_INDIVIDUAL TypeAssigningMode = "INDIVIDUAL"
	ASSIGN_MODE_DEPARTMENT TypeAssigningMode = "DEPARTMENT"
	ASSIGN_MODE_CORPORATE  TypeAssigningMode = "CORPORATE"
)

type EMPLOYEE_TASK_TYPE string

const (
	TYPE_PROGRAM_ASSIGNMENT     EMPLOYEE_TASK_TYPE = "PROGRAM_ASSIGNMENT"
	TYPE_POST_ASSESSMENT        EMPLOYEE_TASK_TYPE = "POST_ASSESSMENT"
	TYPE_PRE_ASSESSMENT         EMPLOYEE_TASK_TYPE = "PRE_ASSESSMENT"
	TYPE_FOLLOWUP_ASSESSMENT    EMPLOYEE_TASK_TYPE = "FOLLOWUP_ASSESSMENT"
	TYPE_MINDTERA_PROGRAM_RECOM EMPLOYEE_TASK_TYPE = "PROGRAM_RECOMMENDATION"
	TYPE_ONBOARDING_QUIZ        EMPLOYEE_TASK_TYPE = "ONBOARDING"
	TYPE_WELL_BEING_ASSESSMENT  EMPLOYEE_TASK_TYPE = "WELL_BEING"
)

// GetTaskTypeConst()
// Get all employee tasks with ordering
func GetTaskTypeConstOrdering() []EMPLOYEE_TASK_TYPE {
	return []EMPLOYEE_TASK_TYPE{
		TYPE_ONBOARDING_QUIZ,
		TYPE_PRE_ASSESSMENT,
		TYPE_POST_ASSESSMENT,
		TYPE_FOLLOWUP_ASSESSMENT,
		TYPE_WELL_BEING_ASSESSMENT,
		TYPE_PROGRAM_ASSIGNMENT,
		TYPE_MINDTERA_PROGRAM_RECOM,
	}
}

func GetProgramRelatedTaskType() []EMPLOYEE_TASK_TYPE {
	return []EMPLOYEE_TASK_TYPE{
		TYPE_MINDTERA_PROGRAM_RECOM,
		TYPE_PROGRAM_ASSIGNMENT,
	}
}

type EMPLOYEE_TASK_STATUS string

const (
	STATUS_ACTIVE   EMPLOYEE_TASK_STATUS = "ACTIVE"
	STATUS_FINISHED EMPLOYEE_TASK_STATUS = "FINISHED"
	STATUS_TAKEN    EMPLOYEE_TASK_STATUS = "TAKEN"
)

func GetAllTaskStatus() []EMPLOYEE_TASK_STATUS {
	return []EMPLOYEE_TASK_STATUS{
		STATUS_ACTIVE,
		STATUS_FINISHED,
	}

}

type EmployeeTaskEntity struct {
	ID             uuid.UUID            `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID            `json:"user_id" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	TaskRelatedID  uuid.UUID            `json:"task_related_id"`
	SubscriptionID uuid.UUID            `json:"subscription_id" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	CorporateId    *uuid.UUID           `json:"corporate_id,omitempty" gorm:"-"`
	Email          string               `form:"email" json:"email"`
	Department     string               `json:"department" gorm:"->;embedded;migration"`
	EmployeeName   string               `json:"employee_name" gorm:"->;embedded;migration"`
	RecordFlag     string               `json:"record_flag,omitempty" gorm:"index;default:'ACTIVE'"`
	StartDate      time.Time            `json:"start_date" gorm:"-"`
	EndDate        time.Time            `json:"end_date" gorm:"-"`
	AssignDate     *time.Time           `json:"assign_date" gorm:"default:NOW();uniqueIndex:unique_corporate_subs_employee_task"`
	DueDate        *time.Time           `json:"due_date" gorm:"default:NOW()"`
	TaskType       EMPLOYEE_TASK_TYPE   `form:"task_type" json:"task_type" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	TaskStatus     EMPLOYEE_TASK_STATUS `json:"task_status"`
	c.DefaultColumn

	DepartmentDetail *ce.CorporateDepartmentEntity `json:"department_detail,omitempty" gorm:"references:Department;foreignKey:Code"`
}

func (e EmployeeTaskEntity) TableName() string {
	return "employee_tasks"
}

type EmployeeTaskEntityForConsumer struct {
	ID             uuid.UUID            `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID            `json:"user_id" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	TaskRelatedID  uuid.UUID            `json:"task_related_id"`
	SubscriptionID uuid.UUID            `json:"subscription_id" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	CorporateId    *uuid.UUID           `json:"corporate_id,omitempty" gorm:"-"`
	Email          string               `form:"email" json:"email"`
	Department     string               `json:"department" gorm:"->;embedded;migration"`
	EmployeeName   string               `json:"employee_name" gorm:"->;embedded;migration"`
	RecordFlag     string               `json:"record_flag,omitempty" gorm:"index;default:'ACTIVE'"`
	StartDate      time.Time            `json:"start_date" gorm:"-"`
	EndDate        time.Time            `json:"end_date" gorm:"-"`
	AssignDate     *time.Time           `json:"assign_date" gorm:"default:NOW();uniqueIndex:unique_corporate_subs_employee_task"`
	DueDate        *time.Time           `json:"-" gorm:"default:NOW()"`
	TaskType       EMPLOYEE_TASK_TYPE   `form:"task_type" json:"task_type" gorm:"type:uuid;uniqueIndex:unique_corporate_subs_employee_task"`
	TaskStatus     EMPLOYEE_TASK_STATUS `json:"task_status"`
	c.DefaultColumn
}

func (e EmployeeTaskEntityForConsumer) TableName() string {
	return "employee_tasks"
}
