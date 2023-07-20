package model

import "github.com/google/uuid"

type CustomDish struct {
	CreatedAt int64                    `json:"createdAt" form:"createdAt"`
	UpdatedAt int64                    `json:"updatedAt" form:"updatedAt"`
	UUID      string                   `json:"uuid" form:"uuid"`
	DishUUID  string                   `json:"dishUuid" form:"dishUuid"`
	Steps     []map[string]interface{} `json:"steps" form:"steps"`
}

type DBCustomDish struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt int64     `gorm:"autoCreateTime"`
	UpdatedAt int64     `gorm:"autoCreateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	DishUUID  uuid.UUID `json:"dishUuid" form:"dishUuid"`
	Steps     string    `json:"steps" form:"steps"`
}

func (DBCustomDish) TableName() string {
	return "custom_dishes"
}
