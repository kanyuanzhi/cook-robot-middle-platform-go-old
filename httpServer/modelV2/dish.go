package modelV2

import (
	"github.com/google/uuid"
	"time"
)

type Dish struct {
	ID        uint32                   `json:"id" form:"id" gorm:"primaryKey"`
	CreatedAt time.Time                `json:"createdAt" form:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time                `json:"updatedAt" form:"updatedAt" gorm:"autoUpdateTime"`
	UUID      uuid.UUID                `json:"uuid" form:"uuid"`
	Name      string                   `json:"name" form:"name"`
	IsTaste   bool                     `json:"isTaste" form:"isTaste"`
	DishUUID  uuid.UUID                `json:"dishUuid" form:"dishUuid"`
	Steps     []map[string]interface{} `json:"steps" form:"steps" gorm:"json"`
	Cuisine   int32                    `json:"cuisine" form:"cuisine"`
	Image     []byte                   `json:"image" form:"image"`
}
