package models

import (
	"gorm.io/datatypes"
)

type DynamicQuestion struct {
	ID       uint           `gorm:"primaryKey"`
	Question string         `json:"question"`
	Options  datatypes.JSON `gorm:"type:jsonb" json:"options"`
	Answers  datatypes.JSON `gorm:"type:jsonb" json:"answers"`
}

func (DynamicQuestion) TableName(category string) string {
	return category // returns the dynamic table name
}
