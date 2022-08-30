package entity

import (
	"fmt"

	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
	"github.com/mindtera/go-common-module/common/v2/model"
)

type CorporateSMM struct {
	ID                      uuid.UUID          `gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Code                    string             `json:"code,omitempty" gorm:"uniqueIndex:smm_code_unique;"`
	Name                    string             `json:"name,omitempty"`
	EnglishName             string             `json:"english_name,omitempty"`
	SequenceNumber          int8               `json:"sequence_number"`
	Description             string             `json:"description,omitempty"`
	DescriptionEnglish      string             `json:"description_english,omitempty"`
	Icon                    string             `json:"icon,omitempty"`
	RecordFlag              string             `json:"record_flag,omitempty" gorm:"column:record_flag;index;default:'ACTIVE'"`
	LocalizationName        model.Localization `json:"localization_name" gorm:"embedded;embeddedPrefix:localization_name_;->"`
	LocalizationDescription model.Localization `json:"localization_description" gorm:"embedded;embeddedPrefix:localization_description_;->"`
	c.DefaultColumn
}

func (c *CorporateSMM) TransformNames() {
	c.LocalizationName.Indonesia = fmt.Sprintf("%v", c.Name)
	c.LocalizationName.English = fmt.Sprintf("%v", c.EnglishName)
	// c.RemoveNames()
}

func (c *CorporateSMM) TransformDescription() {
	c.LocalizationDescription.Indonesia = fmt.Sprintf("%v", c.Description)
	c.LocalizationDescription.English = fmt.Sprintf("%v", c.DescriptionEnglish)
	// c.RemoveDescription()
}

func (c *CorporateSMM) RemoveNames() {
	c.Name = ""
	c.EnglishName = ""
}

func (c *CorporateSMM) RemoveDescription() {
	c.Description = ""
	c.DescriptionEnglish = ""
}
