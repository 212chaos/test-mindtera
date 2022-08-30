package repo_corporateassessmentrelation

import (
	"context"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
)

type CorporateAssessmentRelationRepository interface {
	GetCorporateAssessmentRelation(ctx context.Context, corporateAssessment *[]entity.CorporateAssessmentRelation) (err error)
	GetCorporateAssessmentRelationByValueID(ctx context.Context, corporateValueID []uuid.UUID, corporateAssessment *[]entity.CorporateAssessmentRelation) (err error)
}
