package service_corporatevaluerelation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"

	corpassessrelrepo "github.com/mindtera/corporate-service/repository/corporate-assessment-relation"
	corpemplrepo "github.com/mindtera/corporate-service/repository/corporate-employee"
	corpvalrelrepo "github.com/mindtera/corporate-service/repository/corporate-value-relation"
	corpvalsvc "github.com/mindtera/corporate-service/service/corporate-value"
	"github.com/mindtera/go-common-module/common/logger"
	pb "github.com/mindtera/go-common-module/common/pb"
	googlestorage "github.com/mindtera/go-common-module/common/v2/configuration/google-cloud"
	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
)

type CorporateValueRelationServiceImpl struct {
	sugar             logger.CustomLogger
	gStorage          googlestorage.GoogleStorageConfig
	corpValRelRepo    corpvalrelrepo.CorporateValueRelationRepo
	corpEmplRepo      corpemplrepo.CorporateEmployeeRepository
	corpAssessRelRepo corpassessrelrepo.CorporateAssessmentRelationRepository
	redisSvc          redissvc.RedisSvc
	assert            assert.Assert
	commonService     commonsvc.CommonService
	corpValSvc        corpvalsvc.CorporateValueService
}

func NewCorporateValueRelationService(
	sugar logger.CustomLogger,
	gStorage googlestorage.GoogleStorageConfig,
	corpValRelRepo corpvalrelrepo.CorporateValueRelationRepo,
	corpAssessRelRepo corpassessrelrepo.CorporateAssessmentRelationRepository,
	corpEmplRepo corpemplrepo.CorporateEmployeeRepository,
	redisSvc redissvc.RedisSvc,
	assert assert.Assert,
	commonService commonsvc.CommonService,
	corpValSvc corpvalsvc.CorporateValueService) CorporateValueRelationService {
	return &CorporateValueRelationServiceImpl{
		sugar:             sugar,
		gStorage:          gStorage,
		corpValRelRepo:    corpValRelRepo,
		corpAssessRelRepo: corpAssessRelRepo,
		corpEmplRepo:      corpEmplRepo,
		redisSvc:          redisSvc,
		assert:            assert,
		commonService:     commonService,
		corpValSvc:        corpValSvc,
	}
}

func (c *CorporateValueRelationServiceImpl) GetCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID, corporateValueRelationModel *model.CorporateValueRelationResultModel) (err error) {
	// get corporate raw relation
	rawChan := make(chan []entity.CorporateValueRelation)
	rawErrChan := make(chan error)
	go func(id uuid.UUID) {
		var raw []entity.CorporateValueRelation
		ctx_ := context.WithValue(ctx, commonmodel.TRANSACTION_KEY, nil)
		ctx_ = context.WithValue(ctx_, string(commonmodel.TRANSACTION_KEY), nil)

		rawErrChan <- c.GetCorporateValueRawByCorporateID(ctx_, id, &raw)

		// assign channel value
		rawChan <- raw
	}(corporateID)

	// get parent detail
	chanSMM := make(chan map[string]entity.CorporateSMM)
	chanErr := make(chan error)
	go func() {
		ctx_ := context.WithValue(ctx, commonmodel.TRANSACTION_KEY, nil)
		ctx_ = context.WithValue(ctx_, string(commonmodel.TRANSACTION_KEY), nil)

		hm, er := c.corpValSvc.GetCorporateSMMHash(ctx_)
		chanErr <- er
		chanSMM <- hm
	}()

	// get corporate id and set key
	key := "CORPORATE_VALUE_RELATION_" + corporateID.String()

	// check redis first
	if err = c.redisSvc.Get(ctx, key, corporateValueRelationModel); err != nil {
		// get from database
		c.sugar.WithContext(ctx).Infof("fetching value relation from database:%v", corporateID)
		var assessment []model.CorporateRelationAssessmentModel
		if err = c.corpValRelRepo.GetCorporateValueRelationWithAssessmentByCorporateID(ctx, corporateID, &assessment); err != nil {
			return err
		}

		// check error from channel
		if err = <-chanErr; err != nil {
			c.sugar.WithContext(ctx).Errorf("error in channel smm:%v", err.Error())
			return err
		}
		smmDetail := <-chanSMM

		// check raw channel
		if err = <-rawErrChan; err != nil {
			c.sugar.WithContext(ctx).Errorf("error in raw channel:%v", err.Error())
			return err
		}

		// count total assessment
		records := map[string]bool{}
		totalParent := map[string]model.CorporateValueModel{}

		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for _, k := range assessment {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, v model.CorporateRelationAssessmentModel) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				// check for total assessment
				if !records[v.AssessmentID.String()] {
					records[v.AssessmentID.String()] = true
					corporateValueRelationModel.TotalAssessment++
				}

				// check nil reference
				if !records[v.ParentSMM] {
					records[v.ParentSMM] = true
					totalParent[v.ParentSMM] = model.CorporateValueModel{
						ParentSMM: v.ParentSMM}
				}

				// check total assessment for parent
				if !records[v.ParentSMM+v.AssessmentID.String()] {
					records[v.ParentSMM+v.AssessmentID.String()] = true
					// buffer the parent
					tempParent := totalParent[v.ParentSMM]
					tempParent.TotalAssessment++
					totalParent[v.ParentSMM] = tempParent
				}

				// check records for each parent
				if !records[v.ParentSMM+v.CorporateValueID.String()] {
					records[v.ParentSMM+v.CorporateValueID.String()] = true
					// buffer the parent
					tempParent := totalParent[v.ParentSMM]

					// generate icon url
					iconUrl, iconOutlineUrl := c.corpValSvc.GetValueIcons(v.Code)

					// append record
					tempParentRecord := tempParent.CorporateValueRecords
					tempParentRecord = append(tempParentRecord, entity.CorporateValueEntity{
						ID: v.CorporateValueID,
						// Name:        v.Name,
						// EnglishName: v.EnglishName,
						Code:            v.Code,
						IconLink:        iconUrl,
						IconLinkOutline: iconOutlineUrl,
						LocalizationName: commonmodel.Localization{
							Indonesia: v.Name,
							English:   v.EnglishName,
						}})
					// back to hash map
					tempParent.CorporateValueRecords = tempParentRecord
					totalParent[v.ParentSMM] = tempParent
				}

				// check detail
				if !records["DETAIL"+v.ParentSMM] {
					records["DETAIL"+v.ParentSMM] = true
					tempParent := totalParent[v.ParentSMM]
					tempParent.SMMDetail = smmDetail[v.ParentSMM]
					totalParent[v.ParentSMM] = tempParent
				}
			}(&wg, k)
		}
		wg.Wait()

		// finalize result
		var tempCorporateValueModel []model.CorporateValueModel
		for _, v := range totalParent {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, k model.CorporateValueModel) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				tempCorporateValueModel = append(tempCorporateValueModel, k)
			}(&wg, v)
		}
		wg.Wait()
		// append parent raw and detail
		corporateValueRelationModel.RawData = <-rawChan
		corporateValueRelationModel.Detail = tempCorporateValueModel
	}
	// set redis cache
	go func() {
		if err := c.redisSvc.Set(context.Background(), key, corporateValueRelationModel, time.Duration(common.COMMON_TTL*int(time.Minute))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when set data to redis:%v", key)
		}
	}()

	return err
}

// get corporate value relation by raw only
func (c *CorporateValueRelationServiceImpl) GetCorporateValueByCorporateIDRawOnly(ctx context.Context, corporateID uuid.UUID, corporateValueRelationModel *model.CorporateValueRelationResultModel) (err error) {
	// get corporate raw relation
	rawChan := make(chan []entity.CorporateValueRelation)
	rawErrChan := make(chan error)
	go func(id uuid.UUID) {
		var raw []entity.CorporateValueRelation
		rawErrChan <- c.GetCorporateValueRawByCorporateID(ctx, id, &raw)
		// assign channel value
		rawChan <- raw
	}(corporateID)

	// check raw channel
	if err = <-rawErrChan; err != nil {
		c.sugar.WithContext(ctx).Errorf("error in raw channel:%v", err.Error())
		return err
	}
	corporateValueRelationModel.RawData = <-rawChan
	return err
}

func (c *CorporateValueRelationServiceImpl) GetCorporateValueIds(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValueId) (err error) {
	// validate public user request
	if err = c.validatePublicUser(ctx, corpInfo); err != nil {
		return err
	}

	corporateID := uuid.MustParse(corpInfo.Id)
	key := fmt.Sprintf("CORPORATE_VALUE_IDS_%v_%v", corpInfo.Id, corpInfo.SubsId)
	if err = c.redisSvc.Get(ctx, key, &values); err != nil {
		c.sugar.WithContext(ctx).Infof("searching in database for corp:%v with subs:%v", corpInfo.Id, corpInfo.SubsId)
		if err = c.corpValRelRepo.GetCorpValRelWithAssessByCorporateID(ctx, corporateID, values); err != nil {
			c.sugar.WithContext(ctx).Errorf("error fetching data from database:%v", err.Error())
			return err
		}
	}
	if len(*values) <= 0 {
		err = errors.New("corporate value relation is empty")
		return err
	}
	go c.redisSvc.Set(context.Background(), key, &values, time.Duration(common.COMMON_TTL*int(time.Second)))
	return err
}

func (c *CorporateValueRelationServiceImpl) GetCorporateValueDetail(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValue) (err error) {
	// validate public user request
	// if err = c.validatePublicUser(ctx, corpInfo); err != nil {
	// 	return err
	// }
	corporateID := uuid.MustParse(corpInfo.Id)
	key := fmt.Sprintf("CORPORATE_VALUE_DETAIL_GRPC_%v_%v", corpInfo.Id, corpInfo.SubsId)
	if err = c.redisSvc.Get(ctx, key, &values); err != nil {
		c.sugar.WithContext(ctx).Infof("searching in database for corp:%v with subs:%v", corpInfo.Id, corpInfo.SubsId)
		if err = c.corpValRelRepo.GetCorpValRelDetWithAssessByCorporateID(ctx, corporateID, values); err != nil {
			c.sugar.WithContext(ctx).Errorf("error fetching data from database:%v", err.Error())
			return err
		}
		if len(*values) <= 0 {
			err = errors.New("corporate value relation is empty")
			return err
		}
		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for idx, val := range *values {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, i int, v model.CorpValue) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				// loop over to generate localization and url
				// generate url data
				iconLink, iconOutlineLink := c.corpValSvc.GetValueIcons(v.Detail.Code)

				// populate localization data
				(*values)[i].Detail.TransformNames()
				(*values)[i].Detail.IconLink = iconLink
				(*values)[i].Detail.IconLinkOutline = iconOutlineLink
			}(&wg, idx, val)
		}
		wg.Wait()
		go c.redisSvc.Set(context.Background(), key, &values, time.Duration(common.COMMON_TTL*int(time.Second)))
	}
	return err
}

func (c *CorporateValueRelationServiceImpl) GetCorporateValueRawByCorporateID(ctx context.Context, corporateID uuid.UUID, corporateValueRelation *[]entity.CorporateValueRelation) (err error) {
	// get corporate id and set key
	key := "CORPORATE_VALUE_RELATION_RAW_" + corporateID.String()
	// check redis first
	if err = c.redisSvc.Get(ctx, key, corporateValueRelation); err != nil {
		// get from database
		c.sugar.WithContext(ctx).Infof("fetching value relation from database:%v", corporateID)
		if err = c.corpValRelRepo.GetCorporateValueRelationByCorporateID(ctx, corporateID, corporateValueRelation); err != nil {
			return err
		}

		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for i, val := range *corporateValueRelation {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, idx int, v entity.CorporateValueRelation) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()

				// generate url data
				iconLink, iconOutlineLink := c.corpValSvc.GetValueIcons(v.CorporateValueDetail.Code)

				// populate localization data
				(*corporateValueRelation)[idx].CorporateValueDetail.TransformNames()
				(*corporateValueRelation)[idx].CorporateValueDetail.IconLink = iconLink
				(*corporateValueRelation)[idx].CorporateValueDetail.IconLinkOutline = iconOutlineLink
			}(&wg, i, val)
		}
		wg.Wait()
		// set redis cache
		go c.redisSvc.Set(context.Background(), key, corporateValueRelation, time.Duration(common.COMMON_TTL*int(time.Minute)))
	}
	return err
}

func (c *CorporateValueRelationServiceImpl) CalculateCorporateValueByCorporateRelation(ctx context.Context, corporateValueID []uuid.UUID, corporateValueRelationModel *model.CorporateValueRelationResultModel) (err error) {
	wg, mx := c.commonService.GenerateWaitGroupAndMutex()
	// get raw detail and populate
	rawChan := make(chan []entity.CorporateValueRelation)
	rawChanErr := make(chan error)
	go func() {
		var cids []uuid.UUID
		for _, v := range corporateValueRelationModel.RawData {
			cids = append(cids, v.CorporateValueID)
		}
		// query corporate value detail
		var corporateValues []entity.CorporateValueEntity
		ctx_ := context.Background()
		if err = c.corpValSvc.GetCorporateValueServiceByID(ctx_, cids, &corporateValues); err != nil {
			rawChanErr <- err
			return
		}

		// make hash map from corporate values
		cvMap := map[uuid.UUID]entity.CorporateValueEntity{}
		for _, v := range corporateValues {
			// generate icon url
			iconUrl, iconOutlineUrl := c.corpValSvc.GetValueIcons(v.Code)
			// populate localization
			v.TransformNames()
			v.IconLink = iconUrl
			v.IconLinkOutline = iconOutlineUrl
			cvMap[v.ID] = v
		}

		// loop over the details
		var rawData []entity.CorporateValueRelation
		for _, v := range corporateValueRelationModel.RawData {
			detail := cvMap[v.CorporateValueID]
			v.CorporateValueDetail = &detail
			rawData = append(rawData, v)
		}
		rawChanErr <- nil
		rawChan <- rawData
	}()

	// get parent detail
	chanSMM := make(chan map[string]entity.CorporateSMM)
	chanErr := make(chan error)
	go func() {
		ctx_ := context.Background()
		hm, er := c.corpValSvc.GetCorporateSMMHash(ctx_)
		chanErr <- er
		chanSMM <- hm
	}()
	// check error from channel
	if err = <-chanErr; err != nil {
		c.sugar.WithContext(ctx).Errorf("error in channel smm:%v", err.Error())
		return err
	}
	smmDetail := <-chanSMM

	// fetching corporate assessment from database
	var corporateAssessments []entity.CorporateAssessmentRelation
	if err = c.corpAssessRelRepo.
		GetCorporateAssessmentRelationByValueID(ctx, corporateValueID, &corporateAssessments); err != nil {
		return err
	}

	// count total assessment
	records := map[string]bool{}
	totalParent := map[string]model.CorporateValueModel{}

	// process corporate assessment
	for _, val := range corporateAssessments {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, v entity.CorporateAssessmentRelation) {
			defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
			// lock mutex first
			mx.Lock()

			// check for total assessment
			if !records[v.AssessmentID.String()] {
				records[v.AssessmentID.String()] = true
				corporateValueRelationModel.TotalAssessment++
			}

			// check nil reference
			if !records[v.ParentSMM] {
				records[v.ParentSMM] = true
				totalParent[v.ParentSMM] = model.CorporateValueModel{
					ParentSMM: v.ParentSMM}
			}

			// check total assessment for parent
			if !records[v.ParentSMM+v.AssessmentID.String()] {
				records[v.ParentSMM+v.AssessmentID.String()] = true
				// buffer the parent
				tempParent := totalParent[v.ParentSMM]
				tempParent.TotalAssessment++
				totalParent[v.ParentSMM] = tempParent
			}

			// check records for each parent
			if !records[v.ParentSMM+v.CorporateValueID.String()] {
				records[v.ParentSMM+v.CorporateValueID.String()] = true
				// buffer the parent
				tempParent := totalParent[v.ParentSMM]
				// generate icon url
				iconUrl, iconOutlineUrl := c.corpValSvc.GetValueIcons(v.CorporateValueDetail.Code)
				// append record
				tempParentRecord := tempParent.CorporateValueRecords
				tempParentRecord = append(tempParentRecord, entity.CorporateValueEntity{
					ID:               v.CorporateValueDetail.ID,
					ProgramTaggingID: v.CorporateValueDetail.ProgramTaggingID,
					IconLink:         iconUrl,
					IconLinkOutline:  iconOutlineUrl,
					LocalizationName: commonmodel.Localization{
						Indonesia: v.CorporateValueDetail.Name,
						English:   v.CorporateValueDetail.EnglishName,
					}})
				// back to hash map
				tempParent.CorporateValueRecords = tempParentRecord
				totalParent[v.ParentSMM] = tempParent
			}

			// check detail
			if !records["DETAIL"+v.ParentSMM] {
				records["DETAIL"+v.ParentSMM] = true
				tempParent := totalParent[v.ParentSMM]
				tempParent.SMMDetail = smmDetail[v.ParentSMM]
				totalParent[v.ParentSMM] = tempParent
			}
		}(&wg, val)
	}
	wg.Wait()

	// finalize result
	var tempCorporateValueModel []model.CorporateValueModel
	for _, v := range totalParent {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, k model.CorporateValueModel) {
			defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
			mx.Lock()
			tempCorporateValueModel = append(tempCorporateValueModel, k)
		}(&wg, v)
	}
	wg.Wait()

	if err = <-rawChanErr; err != nil {
		c.sugar.WithContext(ctx).Errorf("error in callback channel:%v", err)
		return err
	}

	// make final model
	corporateValueRelationModel.Detail = tempCorporateValueModel
	corporateValueRelationModel.RawData = <-rawChan
	return err
}

func (c *CorporateValueRelationServiceImpl) UpsertCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID, corporateValue *[]entity.CorporateValueEntity) (err error) {
	// deleting cache
	go c.deleteAllCache()

	// add created by and updated by
	un := c.commonService.GetUserNameFromContext(ctx)

	// loop the value first and generate relation
	var corporateValueRelation []entity.CorporateValueRelation
	for _, v := range *corporateValue {
		// initiate value
		c := entity.CorporateValueRelation{
			CorporateID:      corporateID,
			CorporateValueID: v.ID,
		}
		c.CreatedBy = un
		// append value
		corporateValueRelation = append(corporateValueRelation, c)
	}

	// set to database
	if err := c.corpValRelRepo.UpsertCorporateValueRelationByCorporateID(ctx, &corporateValueRelation); err != nil {
		return err
	}

	return err
}

func (c *CorporateValueRelationServiceImpl) DeleteCorporateValueByCorporateID(ctx context.Context, corporateID uuid.UUID) (err error) {

	// deleting cache
	go c.deleteAllCache()
	// delete corporate value relation
	if err = c.corpValRelRepo.DeleteCorporateValueRelationByCorporateID(ctx, corporateID); err != nil {
		return err
	}

	return err
}

func (c *CorporateValueRelationServiceImpl) GetCorporateValueName(ctx context.Context, corpInfo *pb.CorpInfo, values *[]model.CorpValueName) (err error) {

	// validate public user request
	if err = c.validatePublicUser(ctx, corpInfo); err != nil {
		return err
	}

	corporateID := uuid.MustParse(corpInfo.Id)
	key := fmt.Sprintf("CORPORATE_VALUE_NAMES_%v_%v", corporateID, corpInfo.SubsId)
	if err = c.redisSvc.Get(ctx, key, &values); err != nil {
		c.sugar.WithContext(ctx).Infof("searching in database corp values name for corp:%v with subs:%v", corpInfo.Id, corpInfo.SubsId)
		if err = c.corpValRelRepo.GetCorpValRelNameByCorporateID(ctx, corporateID, values); err != nil {
			c.sugar.WithContext(ctx).Errorf("error fetching data from database:%v", err.Error())
			return err
		}
	}

	if len(*values) <= 0 {
		err = errors.New("corporate value relation is empty")
		return err
	}
	go c.redisSvc.Set(context.Background(), key, &values, time.Duration(common.COMMON_TTL*int(time.Second)))

	return err
}
