package service_corporatevaluerelation

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"

	pb "github.com/mindtera/go-common-module/common/pb"
)

type CorporateValueRelationService interface {
	GetCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID, corporateValueRelation *model.CorporateValueRelationResultModel) (err error)
	GetCorporateValueByCorporateIDRawOnly(ctx context.Context, corporateID uuid.UUID, corporateValueRelation *model.CorporateValueRelationResultModel) (err error)
	GetCorporateValueIds(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValueId) (err error)
	GetCorporateValueDetail(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValue) (err error)
	GetCorporateValueName(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValueName) (err error)
	UpsertCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID, corporateValue *[]entity.CorporateValueEntity) (err error)
	CalculateCorporateValueByCorporateRelation(ctx context.Context, corporateValueID []uuid.UUID, corporateValueRelationModel *model.CorporateValueRelationResultModel) (err error)
	DeleteCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID) (err error)
}
