package service_corporatevalue

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
)

type CorporateValueService interface {
	GetCorporateValueService(ctx context.Context, corporateValues *[]entity.CorporateValueEntity) (err error)
	GetCorporateValueServiceByID(ctx context.Context, cids []uuid.UUID, corporateValues *[]entity.CorporateValueEntity) (err error)
	GetCorporateValueBySMMService(ctx context.Context, corporateValues *[]model.CorporateValueModel) (err error)
	GetCorporateSMMHash(ctx context.Context) (map[string]entity.CorporateSMM, error)
	GetValueIcons(englishName string) (iconUrl, iconOutlineUrl string)
}
