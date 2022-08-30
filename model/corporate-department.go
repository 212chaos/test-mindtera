package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
	"gorm.io/gorm"
)

type CorporateDepartmentGormable struct {
	ID          uuid.UUID `gorm:"column:corporate_department_entities.id"`
	CorporateID uuid.UUID `gorm:"column:corporate_department_entities.corporate_id"`
	Code        string    `gorm:"column:corporate_department_entities.code"`
	Name        string    `gorm:"column:corporate_department_entities.name"`
	CreatedAt   time.Time `gorm:"column:corporate_department_entities.created_at"`
	UpdatedAt   time.Time `gorm:"column:corporate_department_entities.updated_at"`
	DeletedAt   time.Time `gorm:"column:corporate_department_entities.deleted_at"`
	CreatedBy   string    `gorm:"column:corporate_department_entities.created_by"`
	UpdatedBy   string    `gorm:"column:corporate_department_entities.updated_by"`
	RecordFlag  string    `gorm:"column:corporate_department_entities.record_flag"`
}

func (d *CorporateDepartmentGormable) ToCorporateDepartmentEntity() entity.CorporateDepartmentEntity {
	return entity.CorporateDepartmentEntity{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		CorporateID: d.CorporateID,
		DefaultColumn: c.DefaultColumn{
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
			RecordFlag: d.RecordFlag,
			DeletedAt: gorm.DeletedAt{
				Time:  d.DeletedAt,
				Valid: true,
			},
		},
		RecordFlag: d.RecordFlag,
	}
}
