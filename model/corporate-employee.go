package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"gorm.io/gorm"

	ce "github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateEmployeeEntityGormable struct {
	ID                      uuid.UUID `gorm:"column:corporate_employee_entities.id" json:"id"`
	PublicUserID            uuid.UUID `gorm:"column:corporate_employee_entities.public_user_id" json:"public_user_id"`
	CorporateID             uuid.UUID `gorm:"column:corporate_employee_entities.corporate_id" json:"corporate_id"`
	Department              string    `gorm:"column:corporate_employee_entities.department" json:"department"`
	Role                    string    `gorm:"column:corporate_employee_entities.role" json:"role"`
	Email                   string    `gorm:"column:corporate_employee_entities.email" json:"email"`
	EmployeeName            string    `gorm:"column:corporate_employee_entities.employee_name" json:"employee_name"`
	CorporateSubscriptionID uuid.UUID `gorm:"column:corporate_employee_entities.corporate_subscription_id" json:"corporate_subscription_id"`
	CreatedAt               time.Time `gorm:"column:corporate_employee_entities.created_at" json:"created_at"`
	UpdatedAt               time.Time `gorm:"column:corporate_employee_entities.updated_at" json:"updated_at"`
	CreatedBy               string    `gorm:"column:corporate_employee_entities.created_by" json:"created_by"`
	UpdatedBy               string    `gorm:"column:corporate_employee_entities.updated_by" json:"updated_by"`
	RecordFlag              string    `gorm:"column:corporate_employee_entities.record_flag" json:"record_flag"`
}

func (d *CorporateEmployeeEntityGormable) ToCorporateEmployeeEntity() entity.CorporateEmployeeEntity {
	return entity.CorporateEmployeeEntity{
		ID:                      d.ID,
		PublicUserID:            d.PublicUserID,
		CorporateID:             d.CorporateID,
		EmployeeName:            d.EmployeeName,
		Department:              d.Department,
		Email:                   d.Email,
		Role:                    d.Role,
		CorporateSubscriptionID: d.CorporateSubscriptionID,
		DefaultColumn: c.DefaultColumn{
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
			RecordFlag: d.RecordFlag,
		},
	}
}

type CorporateEmployeeEntityGormable2 struct {
	CorporateEmployeeEntityGormable
	CorporateDepartmentGormable
}

func (d *CorporateEmployeeEntityGormable2) ToCorporateEmployeeEntity() entity.CorporateEmployeeEntityWithDepartment {
	employeeTemp := d.CorporateEmployeeEntityGormable.ToCorporateEmployeeEntity()
	return entity.CorporateEmployeeEntityWithDepartment{
		CorporateEmployeeEntity: employeeTemp,
		CorporateDepartment: &ce.CorporateDepartmentEntity{
			ID:              d.CorporateDepartmentGormable.ID,
			Code:            d.CorporateDepartmentGormable.Code,
			Name:            d.CorporateDepartmentGormable.Name,
			CorporateID:     d.CorporateDepartmentGormable.CorporateID,
			CorporateDetail: &ce.CorporateEntity{},
			DefaultColumn: c.DefaultColumn{
				CreatedAt: d.CorporateDepartmentGormable.CreatedAt,
				UpdatedAt: d.CorporateDepartmentGormable.UpdatedAt,
				DeletedAt: gorm.DeletedAt{
					Time:  d.CorporateDepartmentGormable.DeletedAt,
					Valid: true,
				},
				CreatedBy:  d.CorporateDepartmentGormable.CreatedBy,
				UpdatedBy:  d.CorporateDepartmentGormable.UpdatedBy,
				RecordFlag: d.CorporateDepartmentGormable.RecordFlag,
			},
		},
	}

}
