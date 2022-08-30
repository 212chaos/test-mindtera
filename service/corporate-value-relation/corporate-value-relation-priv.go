package service_corporatevaluerelation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	pb "github.com/mindtera/go-common-module/common/pb"
)

func (c *CorporateValueRelationServiceImpl) validatePublicUser(ctx context.Context, corpInfo *pb.CorpInfo) (err error) {
	empl := entity.CorporateEmployeeEntity{
		CorporateSubscriptionID: uuid.MustParse(corpInfo.SubsId),
		PublicUserID:            uuid.MustParse(corpInfo.EmployeeId),
		CorporateID:             uuid.MustParse(corpInfo.Id),
	}
	c.sugar.WithContext(ctx).Infof("checking user with user id:%v", empl.PublicUserID)
	if err = c.corpEmplRepo.GetCorporateEmplByPublicAndSubsID(ctx, &empl); err != nil {
		return err
	}
	// check uuid and email
	if c.assert.IsEmpty(empl.Email) || c.assert.IsUUIDEmpty(empl.ID.String()) {
		err = errors.New("employee not found")
	}
	return err
}

func (c *CorporateValueRelationServiceImpl) deleteAllCache() {
	keys := []string{
		"CORPORATE_VALUE_DETAIL_GRPC_*",
		"CORPORATE_VALUE_RELATION_*",
		"CORPORATE_VALUE_RELATION_RAW_*",
		"CORPORATE_VALUE_IDS_*",
		"CORPORATE_VALUE_NAMES_*"}
	for _, key := range keys {
		c.redisSvc.Delete(context.Background(), key)
	}
}
