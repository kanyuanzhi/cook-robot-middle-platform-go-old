package model

import "github.com/google/uuid"

type Dish struct {
	CreatedAt int64                    `json:"createdAt" form:"createdAt"`
	UpdatedAt int64                    `json:"updatedAt" form:"updatedAt"`
	UUID      string                   `json:"uuid" form:"uuid"`
	Name      string                   `json:"name" form:"name"`
	IsTaste   bool                     `json:"isTaste" form:"isTaste"`
	DishUUID  string                   `json:"dishUuid" form:"dishUuid"`
	Steps     []map[string]interface{} `json:"steps" form:"steps"`
	Cuisine   int32                    `json:"cuisine" form:"cuisine"`
	Image     string                   `json:"image" form:"image"`
}

type DBDish struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt int64     `gorm:"autoCreateTime"`
	UpdatedAt int64     `gorm:"autoCreateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	IsTaste   bool      `json:"isTaste" form:"isTaste"`
	DishUUID  uuid.UUID `json:"dishUuid" form:"dishUuid"`
	Name      string    `json:"name" form:"name"`
	Steps     string    `json:"steps" form:"steps"`
	Cuisine   int32     `json:"cuisine" form:"cuisine"`
	Image     []byte    `json:"image" form:"image"`
}

func (DBDish) TableName() string {
	return "dishes"
}
