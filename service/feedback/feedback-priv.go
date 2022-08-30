package svc_feedback

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
)

func (f *FeedbackImpl) generateFeedbackChannel(feedback []entity.Feedback) <-chan entity.Feedback {
	feedbackChan := make(chan entity.Feedback)
	go func() {
		for _, v := range feedback {
			feedbackChan <- v
		}
		close(feedbackChan)
	}()
	return feedbackChan
}

func (f *FeedbackImpl) filterFeedbackValue(feedbackChan <-chan entity.Feedback) <-chan entity.Feedback {
	filterChan := make(chan entity.Feedback)
	go func() {
		for v := range feedbackChan {
			// do process filtering value

			// remove user id
			v.UserId = uuid.Nil
			v.EmployeeDetail.PublicUserID = uuid.Nil

			// filter by shown category
			switch v.ShownCategoryCode {
			case "ANONYMOUS":
				v.EmployeeDetail = nil
			case "DEPARTMENT":
				v.EmployeeDetail.EmployeeName = ""
			}
			filterChan <- v
		}
		close(filterChan)
	}()
	return filterChan
}

func (f *FeedbackImpl) feedbackFanIn(fanInChan ...<-chan entity.Feedback) <-chan entity.Feedback {
	finalChan := make(chan entity.Feedback)
	go func() {
		var wg sync.WaitGroup
		for _, v := range fanInChan {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, ch <-chan entity.Feedback) {
				defer wg.Done()
				for i := range ch {
					finalChan <- i
				}
			}(&wg, v)
		}
		wg.Wait()
		close(finalChan)
	}()
	return finalChan
}

func (f *FeedbackImpl) validateCategory(ctx context.Context, resChan chan error, category, shownCategory string) {
	f.sugar.WithContext(ctx).Info("validate category")
	go func(cat string) {
		if cat == "" {
			resChan <- nil
			return
		}

		var category []entity.FeedbackCategory
		if errMsg := f.GetFeedbackCategorySvc(ctx, &category); errMsg.Error != nil {
			resChan <- errMsg.Error
		}
		for _, v := range category {
			if v.Code == cat {
				resChan <- nil
				return
			}
		}
		resChan <- errors.New("category is invalid")
	}(category)

	f.sugar.WithContext(ctx).Info("validate shown category")
	go func(cat string) {
		if cat == "" {
			resChan <- nil
			return
		}

		var category []entity.FeedbackShownCategory
		if errMsg := f.GetFeedbackShownCategorySvc(ctx, &category); errMsg.Error != nil {
			resChan <- errMsg.Error
		}
		for _, v := range category {
			if v.Code == cat {
				resChan <- nil
				return
			}
		}
		resChan <- errors.New("shown category is invalid")
	}(shownCategory)
}

func (f *FeedbackImpl) setRedisValue(ctx context.Context, key string, data any, ttl int) {
	f.sugar.WithContext(ctx).Infof("set value for key:%v", key)
	if err := f.redisService.Set(ctx, key, data, time.Duration(ttl*int(time.Hour))); err != nil {
		f.sugar.WithContext(ctx).Errorf("error set value for key:%v", err)
		return
	}
}

func (f *FeedbackImpl) deleteFeedbackRedisCache(ctx context.Context, key string) {
	f.sugar.WithContext(ctx).Infof("delete value for key:%v", key)
	if err := f.redisService.Delete(ctx, key); err != nil {
		f.sugar.WithContext(ctx).Errorf("error delete value for key:%v", err)
		return
	}
}
