package repository_corporateprogram

import (
	"context"

	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
)

type CorporateProgram interface {
	GetCorporateEmployeeProgram(ctx context.Context, employee entity.EmployeeTaskEntity) (result []model.EmployeeProgram, err error)
}
