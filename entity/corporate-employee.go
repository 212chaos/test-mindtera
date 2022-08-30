package entity

import (
	"time"

	"github.com/google/uuid"
	ce "github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
	fmodel "github.com/mindtera/go-common-module/fajar/model"
)

type CorporateEmployeeEntity struct {
	Number                  int       `json:"no,omitempty" gorm:"-"`
	ID                      uuid.UUID `json:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" `
	EmployeeName            string    `json:"employee_name,omitempty"`
	PublicUserID            uuid.UUID `json:"public_user_id" gorm:"type:uuid;index;"`
	CorporateID             uuid.UUID `json:"corporate_id" gorm:"type:uuid;index;"`
	CorporateSubscriptionID uuid.UUID `json:"corporate_subscription_id" gorm:"type:uuid;"`
	Department              string    `json:"department" binding:"required"`
	Role                    string    `json:"role" binding:"required"`
	Email                   string    `json:"email" binding:"required,email"`
	Corporate               string    `json:"corporate" gorm:"-"`
	RecordFlag              string    `json:"record_flag,omitempty" gorm:"column:record_flag;index:record_value_corporate_index;default:'UNREGISTERED'"`
	c.DefaultColumn

	CorporateDetail  *ce.CorporateEntity           `json:"corporate_detail,omitempty" gorm:"references:CorporateID;foreignKey:ID"`
	RoleDetail       *ce.CorporateRoleEntity       `json:"role_detail,omitempty" gorm:"references:Role;foreignKey:Code"`
	DepartmentDetail *ce.CorporateDepartmentEntity `json:"department_detail,omitempty" gorm:"references:Department;foreignKey:Code"`
}

type CorporateEmployeeJson struct {
	Number                  int       `json:"corporate_employee_entities.number"`
	ID                      uuid.UUID `json:"corporate_employee_entities.id"`
	EmployeeName            string    `json:"corporate_employee_entities.employee_name"`
	PublicUserID            uuid.UUID `json:"corporate_employee_entities.public_user_id"`
	CorporateID             uuid.UUID `json:"corporate_employee_entities.corporate_id"`
	CorporateSubscriptionID uuid.UUID `json:"corporate_employee_entities.coporate_subscription_id"`
	Department              string    `json:"corporate_employee_entities.department"`
	Role                    string    `json:"corporate_employee_entities.role"`
	Email                   string    `json:"corporate_employee_entities.email"`
	Corporate               string    `json:"corporate_employee_entities.corporate"`
	RecordFlag              string    `json:"corporate_employee_entities.record_flag"`
	CreatedAt               time.Time `json:"corporate_employee_entities.created_at"`
	UpdatedAt               time.Time `json:"corporate_employee_entities.updated_at"`
	CreatedBy               string    `json:"corporate_employee_entities.created_by"`
	UpdatedBy               string    `json:"corporate_employee_entities.updated_by"`
}

// CorporateEmployeeEntityWithDepartment
// contains:
// CorporateEntity, CorporateDepartmentEntity
type CorporateEmployeeEntityWithDepartment struct {
	CorporateEmployeeEntity
	CorporateDepartment *ce.CorporateDepartmentEntity
}

type CorporateEmployeeEntityWithUser struct {
	CorporateEmployeeEntity
	User *fmodel.User `json:"user"`
}

type CorporateEmployeeEntityMap struct {
	ID                      uuid.UUID `gorm:"column:corporate_employee_entities.id" json:"corporate_employee_entities.id"`
	PublicUserID            uuid.UUID `gorm:"column:corporate_employee_entities.public_user_id" json:"corporate_employee_entities.public_user_id"`
	CorporateID             uuid.UUID `gorm:"column:corporate_employee_entities.corporate_id" json:"corporate_employee_entities.corporate_id"`
	Department              string    `gorm:"column:corporate_employee_entities.department" json:"corporate_employee_entities.department"`
	Role                    string    `gorm:"column:corporate_employee_entities.role" json:"corporate_employee_entities.role"`
	Email                   string    `gorm:"column:corporate_employee_entities.email" json:"corporate_employee_entities.email"`
	EmployeeName            string    `gorm:"column:corporate_employee_entities.employee_name" json:"corporate_employee_entities.employee_name"`
	CorporateSubscriptionID uuid.UUID `gorm:"column:corporate_employee_entities.corporate_subscription_id" json:"corporate_employee_entities.corporate_subscription_id"`
	CreatedAt               time.Time `gorm:"column:corporate_employee_entities.created_at" json:"corporate_employee_entities.created_at"`
	UpdatedAt               time.Time `gorm:"column:corporate_employee_entities.updated_at" json:"corporate_employee_entities.updated_at"`
	CreatedBy               string    `gorm:"column:corporate_employee_entities.created_by" json:"corporate_employee_entities.created_by"`
	UpdatedBy               string    `gorm:"column:corporate_employee_entities.updated_by" json:"corporate_employee_entities.updated_by"`
	RecordFlag              string    `gorm:"column:corporate_employee_entities.record_flag" json:"corporate_employee_entities.record_flag"`
}

type CorporateEmployeeWithUserMap struct {
	CorporateEmployeeEntityMap
	fmodel.UserMap
}

type CorporateEmployeeJsonWithUser struct {
	CorporateEmployeeJson
	User fmodel.UserJson
}

func (d *CorporateEmployeeJsonWithUser) GetEntity() CorporateEmployeeEntityWithUser {
	return CorporateEmployeeEntityWithUser{
		CorporateEmployeeEntity: CorporateEmployeeEntity{
			ID:                      d.CorporateEmployeeJson.ID,
			EmployeeName:            d.CorporateEmployeeJson.EmployeeName,
			PublicUserID:            d.CorporateEmployeeJson.PublicUserID,
			CorporateID:             d.CorporateEmployeeJson.CorporateID,
			CorporateSubscriptionID: d.CorporateEmployeeJson.CorporateSubscriptionID,
			Department:              d.CorporateEmployeeJson.Department,
			Role:                    d.CorporateEmployeeJson.Role,
			Email:                   d.CorporateEmployeeJson.Email,
			DefaultColumn: c.DefaultColumn{
				CreatedAt:  d.CorporateEmployeeJson.CreatedAt,
				UpdatedAt:  d.CorporateEmployeeJson.UpdatedAt,
				CreatedBy:  d.CorporateEmployeeJson.CreatedBy,
				UpdatedBy:  d.CorporateEmployeeJson.UpdatedBy,
				RecordFlag: d.CorporateEmployeeJson.RecordFlag,
			},
			RecordFlag: d.CorporateEmployeeJson.RecordFlag,
		},
		User: &fmodel.User{
			ID:              d.User.ID,
			DOB:             &d.User.DOB,
			Nickname:        d.User.Nickname,
			AccountType:     d.User.AccountType,
			PersonaType:     d.User.PersonaType,
			Gender:          d.User.Gender,
			Name:            d.User.Name,
			Email:           d.User.Email,
			PhoneNumber:     d.User.PhoneNumber,
			EmailVerifiedAt: &d.User.EmailVerifiedAt,
			LastLoginMethod: d.User.LastLoginMethod,
			DefaultColumns: fmodel.DefaultColumns{
				CreatedAt:  &d.User.CreatedAt,
				UpdatedAt:  &d.User.UpdatedAt,
				CreatedBy:  d.User.CreatedBy,
				UpdatedBy:  d.User.UpdatedBy,
				RecordFlag: fmodel.TRecordFlag(d.User.RecordFlag),
			},
		},
	}
}
