package modelV2

import (
	"github.com/google/uuid"
	"time"
)

type Seasoning struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	Name      string    `json:"name" form:"name"`
	Pump      uint32    `json:"pump" form:"pump"`
	Ratio     uint32    `json:"ratio" form:"ratio"`
}

func (Seasoning) TableName() string {
	return "seasonings"
}
