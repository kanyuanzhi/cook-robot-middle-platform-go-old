package modelV2

import (
	"github.com/google/uuid"
	"time"
)

type Cuisine struct {
	ID        uint32    `json:"id" form:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt" gorm:"autoUpdateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	Name      string    `json:"name" form:"name"`
	SortID    uint8     `json:"sortID" form:"sortID"`
}
