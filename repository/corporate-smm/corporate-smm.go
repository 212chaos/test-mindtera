package repo_corporatesmm

import (
	"context"

	"github.com/mindtera/corporate-service/entity"
)

type CorporateSMMRepository interface {
	GetCorporateSMM(ctx context.Context, corporateSMM *[]entity.CorporateSMM) (err error)
	UpsertCorporateSMM(ctx context.Context, corporateSMM *entity.CorporateSMM) (err error)
}
