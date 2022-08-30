package entity

import (
	"fmt"

	"github.com/google/uuid"
	cp "github.com/mindtera/dashboard-program/entity"
	c "github.com/mindtera/go-common-module/common/entity"
	"github.com/mindtera/go-common-module/common/v2/model"
)

type CorporateValueEntity struct {
	Hide             bool      `json:"hide,omitempty" gorm:"default:false"`
	SequenceNumber   int       `json:"sequence_number,omitempty" gorm:"index"`
	ID               uuid.UUID `gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	ProgramTaggingID uuid.UUID `json:"program_tagging_id,omitempty" gorm:"type:uuid;"`
	Code             string    `gorm:"unique;not null;index:corporate_value_code_index" json:"code,omitempty"`
	Name             string    `json:"name,omitempty"`
	EnglishName      string    `json:"english_name,omitempty"`
	RecordFlag       string    `json:"record_flag,omitempty" gorm:"column:record_flag;index:record_value_corporate_index;default:'ACTIVE'"`
	IconLink         string    `json:"icon_link"`
	IconLinkOutline  string    `json:"icon_link_outline"`

	LocalizationName model.Localization `json:"localization_name" gorm:"embedded;embeddedPrefix:localization_name_;->"`
	c.DefaultColumn

	ProgramTaggingDetail *cp.ProgramTaggingEntity `json:"program_tagging_detail,omitempty" gorm:"references:ID;foreignKey:ProgramTaggingID;"`
}

func (c *CorporateValueEntity) TransformNames() {
	c.LocalizationName.Indonesia = fmt.Sprintf("%v", c.Name)
	c.LocalizationName.English = fmt.Sprintf("%v", c.EnglishName)
	// c.RemoveNames()
}

func (c *CorporateValueEntity) RemoveNames() {
	c.Name = ""
	c.EnglishName = ""
}
