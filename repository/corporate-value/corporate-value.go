package corporatevaluerepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
)

type CorporateValueRepository interface {
	GetCorporateValue(ctx context.Context, corpVal *[]entity.CorporateValueEntity) (err error)
	GetCorporateValueByID(ctx context.Context, cids []uuid.UUID, corpVal *[]entity.CorporateValueEntity) (err error)
	GetCorporateValueByCode(ctx context.Context, corporateValue *entity.CorporateValueEntity) (err error)
}
