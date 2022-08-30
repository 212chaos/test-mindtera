package svc_feedback

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	"github.com/mindtera/go-common-module/common/logger"

	repo "github.com/mindtera/corporate-service/repository/feedback"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
	commonservice "github.com/mindtera/go-common-module/common/v2/service/common"
	redisservice "github.com/mindtera/go-common-module/common/v2/service/redis"
)

const (
	FeedbackCategoryKey      = "FEEDBACK_CATEGORY"
	FeedbackShownCategoryKey = "FEEDBACK_SHOWN_CATEGORY"
	FeedbackKey              = "FEEDBACK"

	// in hour
	FeedbackCategoryTtl      int = 108
	FeedbackShownCategoryTtl int = 108
	FeedbackTtl              int = 24
)

type FeedbackImpl struct {
	sugar         logger.CustomLogger
	commonService commonservice.CommonService
	redisService  redisservice.RedisSvc
	feedbackRepo  repo.FeedbackRepository
	maxGoRoutine  int
}

func NewFeedbackSvc(sugar logger.CustomLogger,
	commonService commonservice.CommonService,
	redisService redisservice.RedisSvc,
	feedbackRepo repo.FeedbackRepository) Feedback {

	return &FeedbackImpl{
		sugar:         sugar,
		commonService: commonService,
		redisService:  redisService,
		feedbackRepo:  feedbackRepo,
		maxGoRoutine:  5,
	}
}

func (f *FeedbackImpl) GetFeedbackByFilterSvc(ctx context.Context, feedbackFilter model.FeedbackFilter, feedbackPaging *mdl.PaginationResponseModel) (errMsg mdl.ErrorMessage) {
	// getting corporate id
	corpId := ctx.Value("corporate_id")
	if corpId == nil {
		f.sugar.WithContext(ctx).Errorf("error when getting corporate_id")
		return mdl.ErrorMessage{
			Error:     errors.New("unauthorized"),
			ErrorType: mdl.ERR_UNAUTHORIZED_TYPE}
	}
	var corpUUID uuid.UUID
	f.commonService.ObjectMapper(&corpId, &corpUUID)
	// updateing filter value
	feedbackFilter.CorporateId = corpUUID
	feedbackFilter.TransformIntToTime()

	// validate category code
	resChan := make(chan error)
	ctxBack := f.commonService.ContextBackground(ctx)
	go f.validateCategory(ctxBack, resChan, feedbackFilter.Category, feedbackFilter.ShownCategory)
	for i := 0; i < 2; i++ {
		if err := <-resChan; err != nil {
			return mdl.ErrorMessage{
				Error:     err,
				ErrorType: mdl.ERR_STANDARD_BAD_REQUEST_TYPE,
			}
		}
	}

	// check filter and make key
	feedbackKey := fmt.Sprintf("%v:%v:%v_%v_%v_%v:%v_%v",
		FeedbackKey,
		corpId,
		feedbackFilter.StartDateInt,
		feedbackFilter.EndDateInt,
		feedbackFilter.Category,
		feedbackFilter.ShownCategory,
		feedbackPaging.MetaData.Page,
		feedbackPaging.MetaData.PageSize)

	var feedback []entity.Feedback
	f.sugar.WithContext(ctx).Infof("performing fetch from redis for:%v", feedbackKey)
	if err := f.redisService.Get(ctx, feedbackKey, feedbackPaging); err != nil {
		f.sugar.WithContext(ctx).Infof("performing fetch from database for:%v", feedbackFilter)
		if err := f.feedbackRepo.GetFeedbackByFilter(ctx, feedbackFilter, &feedbackPaging.MetaData, &feedback); err != nil {
			f.sugar.WithContext(ctx).Errorf("error perform fetch feedback err:%v", err)
			errMsg = mdl.ErrorMessage{
				Error:     err,
				ErrorType: mdl.ERR_STANDARD_INTERNAL_TYPE,
			}
			return errMsg
		}
		// make channel
		maxGo := f.maxGoRoutine
		if len(feedback) < maxGo {
			maxGo = len(feedback)
		}
		feedbackChan := f.generateFeedbackChannel(feedback)
		// process channel
		var filterChannel ([]<-chan entity.Feedback)
		for i := 0; i < maxGo; i++ {
			filterChannel = append(filterChannel, f.filterFeedbackValue(feedbackChan))
		}
		resultChan := f.feedbackFanIn(filterChannel...)
		// get value back from channel
		finalFeedback := []entity.Feedback{}
		for i := range resultChan {
			finalFeedback = append(finalFeedback, i)
		}

		sort.Slice(finalFeedback, func(i, j int) bool {
			return finalFeedback[i].CreatedAt.After(finalFeedback[j].CreatedAt)
		})

		feedbackPaging.RawData = finalFeedback
	}

	// set to redis
	if len(feedback) > 0 {
		ctxBackground := f.commonService.ContextBackground(ctx)
		go f.setRedisValue(ctxBackground, feedbackKey, feedback, FeedbackTtl)
	}
	return errMsg
}

func (f *FeedbackImpl) UpsertFeedbackSvc(ctx context.Context, feedback *entity.Feedback) (errMsg mdl.ErrorMessage) {
	// validate category code
	resChan := make(chan error)
	ctxBack := f.commonService.ContextBackground(ctx)
	go f.validateCategory(ctxBack, resChan, feedback.CategoryCode, feedback.ShownCategoryCode)
	for i := 0; i < 2; i++ {
		if err := <-resChan; err != nil {
			return mdl.ErrorMessage{
				Error:     errors.New("BAD_REQUEST:" + err.Error()),
				ErrorType: mdl.ERR_STANDARD_BAD_REQUEST_TYPE,
			}
		}
	}

	f.sugar.WithContext(ctx).Infof("got feedback :%v", feedback.ID)
	// getting corporate id
	corpUUID := feedback.CorporateId

	// upsert feedback
	f.sugar.WithContext(ctx).Infof("perform upsert feedback for:%v", feedback.ID)
	if err := f.feedbackRepo.UpsertFeedback(ctx, feedback); err != nil {
		f.sugar.WithContext(ctx).Errorf("error perform upsert feedback err:%v", err)
		errMsg = mdl.ErrorMessage{
			Error:     err,
			ErrorType: mdl.ERR_STANDARD_INTERNAL_TYPE,
		}
		return errMsg
	}
	// delete cache
	ctxBackground := f.commonService.ContextBackground(ctx)
	go f.deleteFeedbackRedisCache(ctxBackground, fmt.Sprintf("%v:%v:*", FeedbackKey, corpUUID))
	return errMsg
}

func (f *FeedbackImpl) MarkFeedbackAsReadSvc(ctx context.Context, feedbackId uuid.UUID) (errMsg mdl.ErrorMessage) {
	f.sugar.WithContext(ctx).Infof("mark feedback as read with id %v", feedbackId)
	// getting corporate id
	corpId := ctx.Value("corporate_id")
	if corpId == nil {
		f.sugar.WithContext(ctx).Errorf("error when getting corporate_id")
		return mdl.ErrorMessage{
			Error:     errors.New("unauthorized"),
			ErrorType: mdl.ERR_UNAUTHORIZED_TYPE}
	}
	var corpUUID uuid.UUID
	f.commonService.ObjectMapper(&corpId, &corpUUID)

	// processing corp uuid
	if err := f.feedbackRepo.MarkFeedbackAsRead(ctx, feedbackId, corpUUID); err != nil {
		f.sugar.WithContext(ctx).Errorf("error when update as read")
		return mdl.ErrorMessage{
			Error:     errors.New("error when updating flag"),
			ErrorType: mdl.ERR_STANDARD_INTERNAL_TYPE}
	}

	// delete cache
	ctxBackground := f.commonService.ContextBackground(ctx)
	go f.deleteFeedbackRedisCache(ctxBackground, fmt.Sprintf("%v:%v:*", FeedbackKey, corpUUID))
	return errMsg
}

func (f *FeedbackImpl) GetFeedbackCategorySvc(ctx context.Context, feedbackCategory *[]entity.FeedbackCategory) (errMsg mdl.ErrorMessage) {
	f.sugar.WithContext(ctx).Infof("getting feedback category")
	if err := f.redisService.Get(ctx, FeedbackCategoryKey, feedbackCategory); err != nil {
		f.sugar.WithContext(ctx).Infof("error getting redis value from key :%v err:%v", FeedbackCategoryKey, err)
		if err = f.feedbackRepo.GetFeedbackCategory(ctx, feedbackCategory); err != nil {
			f.sugar.WithContext(ctx).Errorf("error getting data from databae err:%v", FeedbackCategoryKey, err)
			errMsg = mdl.ErrorMessage{
				Error:     errors.New("error fetching from database"),
				ErrorType: mdl.ERR_STANDARD_INTERNAL_TYPE,
			}
			return errMsg
		}
	}

	// set value in redis
	if len(*feedbackCategory) > 0 {
		ctxBackground := f.commonService.ContextBackground(ctx)
		go f.setRedisValue(ctxBackground, FeedbackCategoryKey, feedbackCategory, FeedbackCategoryTtl)
	}

	return errMsg
}

func (f *FeedbackImpl) GetFeedbackShownCategorySvc(ctx context.Context, feedbackCategory *[]entity.FeedbackShownCategory) (errMsg mdl.ErrorMessage) {
	f.sugar.WithContext(ctx).Infof("getting feedback shown cat")
	if err := f.redisService.Get(ctx, FeedbackShownCategoryKey, feedbackCategory); err != nil {
		f.sugar.WithContext(ctx).Infof("error getting redis feedback shown cat value from key :%v err:%v", FeedbackCategoryKey, err)
		if err = f.feedbackRepo.GetFeedbackShownCategory(ctx, feedbackCategory); err != nil {
			f.sugar.WithContext(ctx).Errorf("error getting feedback shown cat data from databae err:%v", FeedbackCategoryKey, err)
			errMsg = mdl.ErrorMessage{
				Error:     errors.New("error fetching from database"),
				ErrorType: mdl.ERR_STANDARD_INTERNAL_TYPE,
			}
			return errMsg
		}
	}

	// set value in redis
	if len(*feedbackCategory) > 0 {
		ctxBackground := f.commonService.ContextBackground(ctx)
		go f.setRedisValue(ctxBackground, FeedbackShownCategoryKey, feedbackCategory, FeedbackShownCategoryTtl)
	}
	return errMsg
}
