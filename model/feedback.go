package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
	mdl "github.com/mindtera/go-common-module/common/v2/model"
)

type Feedbacks struct {
	FeedbackArr []entity.Feedback `json:"feedback_arr"`
}

type FeedbackFilter struct {
	CorporateId   uuid.UUID `json:"corporate_id"`
	Category      string    `json:"category" form:"category"`
	ShownCategory string    `json:"shown_category" form:"shown_category"`
	EndDateInt    int64     `form:"end_date"`
	StartDateInt  int64     `form:"start_date"`
	EndDate       time.Time `json:"end_time"`
	StartDate     time.Time `json:"start_time"`
	mdl.PaginationModel
}

type FeedbackQueryFilter struct {
	ShownCategoryCode string    `json:"shown_category_code"`
	CategoryCode      string    `json:"category_code"`
	CorporateId       uuid.UUID `json:"corporate_id"`
}

func (f *FeedbackFilter) TransformIntToTime() {
	f.EndDate = time.UnixMilli(f.EndDateInt)
	f.StartDate = time.UnixMilli(f.StartDateInt)
}
