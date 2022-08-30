package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"

	c "github.com/mindtera/go-common-module/common/entity"
	defaultcommonmodelv2 "github.com/mindtera/go-common-module/common/v2/model"
	commonmodel "github.com/mindtera/go-common-module/fajar/model"
)

type EmployeeTaskRequest struct {
	EmployeeTasks []entity.EmployeeTaskEntity `json:"employee_tasks" binding:"required"`
}

type EmployeeTaskQuery struct {
	Types     []string  `json:"types"`
	Status    string    `json:"status"`
	UserId    uuid.UUID `json:"public_id"`
	StartDate time.Time `json:"start_date_go"`
	EndDate   time.Time `json:"end_date_go"`
}

type EmployeeProgram struct {
	ProgramId uuid.UUID `json:"program_id,omitempty" gorm:"->;column:task_related_id"`
	Amount    int32     `json:"amount,omitempty" gorm:"->;column:amount"`
}

type EmployeeTaskEntityGormable struct {
	ID             uuid.UUID                   `gorm:"column:employee_tasks.id"`
	UserID         uuid.UUID                   `gorm:"column:employee_tasks.user_id"`
	AssignDate     *time.Time                  `gorm:"column:employee_tasks.assign_date"`
	DueDate        *time.Time                  `gorm:"column:employee_tasks.due_date"`
	TaskType       entity.EMPLOYEE_TASK_TYPE   `gorm:"column:employee_tasks.task_type"`
	TaskRelatedID  uuid.UUID                   `gorm:"column:employee_tasks.task_related_id"`
	TaskStatus     entity.EMPLOYEE_TASK_STATUS `gorm:"column:employee_tasks.task_status"`
	SubscriptionID uuid.UUID                   `gorm:"column:employee_tasks.subscription_id"`
	CreatedAt      *time.Time                  `gorm:"column:employee_tasks.created_at"`
	UpdatedAt      *time.Time                  `gorm:"column:employee_tasks.updated_at"`
	CreatedBy      string                      `gorm:"column:employee_tasks.created_by"`
	UpdatedBy      string                      `gorm:"column:employee_tasks.updated_by"`
	RecordFlag     string                      `gorm:"column:employee_tasks.record_flag"`
	*CorporateEmployeeEntityGormable
}

func (d EmployeeTaskEntityGormable) ToEmployeeTaskEntityForConsumer() entity.EmployeeTaskEntityForConsumer {
	var createdAt, updatedAt time.Time
	if d.CreatedAt != nil {
		createdAt = *d.CreatedAt
	}
	if d.UpdatedAt != nil {
		updatedAt = *d.UpdatedAt
	}
	return entity.EmployeeTaskEntityForConsumer{
		ID:             d.ID,
		UserID:         d.UserID,
		AssignDate:     d.AssignDate,
		DueDate:        d.DueDate,
		TaskType:       d.TaskType,
		TaskRelatedID:  d.TaskRelatedID,
		TaskStatus:     d.TaskStatus,
		SubscriptionID: d.SubscriptionID,
		DefaultColumn: c.DefaultColumn{
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
			RecordFlag: d.RecordFlag,
		},
		Email:        d.CorporateEmployeeEntityGormable.Email,
		EmployeeName: d.CorporateEmployeeEntityGormable.EmployeeName,
		Department:   d.CorporateEmployeeEntityGormable.Department,
	}
}

// employee task request payload
type ProgramAssignmentRequest struct {
	AssigningMode             string   `json:"assigning_mode"`
	ID                        string   `json:"id"`
	AssignDate                string   `json:"assign_date"`
	ProgramID                 []string `json:"program_id"`
	ConfirmIgnoreTakenProgram bool     `json:"confirm_ignore_taken_program"`
}

func (d ProgramAssignmentRequest) GenerateCreateModel() (ProgramAssignmentCreation, error) {
	var output ProgramAssignmentCreation
	switch d.AssigningMode {
	case string(entity.ASSIGN_MODE_CORPORATE):
		output.AssigningMode = entity.ASSIGN_MODE_CORPORATE
	case string(entity.ASSIGN_MODE_DEPARTMENT):
		output.AssigningMode = entity.ASSIGN_MODE_DEPARTMENT
	case string(entity.ASSIGN_MODE_INDIVIDUAL):
		output.AssigningMode = entity.ASSIGN_MODE_INDIVIDUAL
	}

	assignDate, err := time.Parse(commonmodel.DATE_YYYYMMDD, d.AssignDate)
	if err != nil {
		return ProgramAssignmentCreation{}, err
	}

	output.AssignDate = assignDate
	output.RelatedID = d.ID
	output.ProgramID = d.ProgramID
	output.ConfirmIgnoreTakenProgram = d.ConfirmIgnoreTakenProgram
	return output, nil
}

type ProgramAssignmentCreation struct {
	AssigningMode entity.TypeAssigningMode
	// RelatedID
	// related id or department_code
	RelatedID                 string
	AssignDate                time.Time
	ProgramID                 []string
	ConfirmIgnoreTakenProgram bool
}

// search program request payload
type SearchProgramRequest struct {
	Query string `form:"query"`
}

func (d SearchProgramRequest) GenerateModel() (string, bool) {
	if d.Query == "" {
		return "", true
	}
	return d.Query, true
}

type EmployeeTaskWithPriority struct {
	entity.EmployeeTaskEntityForConsumer
	Priority    int         `json:"priority"`
	TaskDetails interface{} `json:"task_details"`
}

func (d *EmployeeTaskWithPriority) FillUpPriority() {
	taskTypeOrdering := entity.GetTaskTypeConstOrdering()
	type tempOrdering struct {
		TaskType entity.EMPLOYEE_TASK_TYPE
		Priority int
	}
	taskPriorityData := make([]tempOrdering, len(taskTypeOrdering))
	priority := 1
	for i, v := range taskTypeOrdering {
		taskPriorityData[i] = tempOrdering{
			TaskType: v,
			Priority: priority,
		}
		if d.TaskType == v {
			d.Priority = priority
			break
		}
		priority++
	}

}

type AssignedProgramRequestForm struct {
	defaultcommonmodelv2.PaginationModel
	Query string `form:"query"`
}

func (d *AssignedProgramRequestForm) ValidateRequest() {
	d.ValidatePaging()
	d.Query = strings.ToUpper(d.Query)
}
