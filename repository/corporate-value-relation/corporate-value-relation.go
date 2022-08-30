package corporatevaluerelationrepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
)

type CorporateValueRelationRepo interface {
	GetCorporateValueRelationByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValRel *[]entity.CorporateValueRelation) (err error)
	GetCorporateValueRelationWithAssessmentByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValRel *[]model.CorporateRelationAssessmentModel) (err error)
	GetCorpValRelNameByCorporateID(ctx context.Context, corporateID uuid.UUID, corpValName *[]model.CorpValueName) (err error)
	GetCorpValRelWithAssessByCorporateID(ctx context.Context, corpId uuid.UUID, corpValRel *[]model.CorpValueId) (err error)
	GetCorpValRelDetWithAssessByCorporateID(ctx context.Context, corpId uuid.UUID, corpValRel *[]model.CorpValue) (err error)
	UpsertCorporateValueRelationByCorporateID(ctx context.Context, corpValRel *[]entity.CorporateValueRelation) (err error)
	DeleteCorporateValueRelationByCorporateID(ctx context.Context, corporateID uuid.UUID) (err error)
}
