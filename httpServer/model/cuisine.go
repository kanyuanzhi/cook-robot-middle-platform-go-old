package model

import "github.com/google/uuid"

type Cuisine struct {
	ID        uint32 `json:"id" form:"id"`
	CreatedAt int64  `json:"createdAt" form:"createdAt"`
	UpdatedAt int64  `json:"updatedAt" form:"updatedAt"`
	UUID      string `json:"uuid" form:"uuid"`
	Name      string `json:"name" form:"name"`
	SortID    uint8  `json:"sortID" form:"sortID"`
}

type DBCuisine struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt int64     `gorm:"autoCreateTime"`
	UpdatedAt int64     `gorm:"autoCreateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	Name      string    `json:"name" form:"name"`
	SortID    uint8     `json:"sortID" form:"sortID"`
}

func (DBCuisine) TableName() string {
	return "cuisines"
}
