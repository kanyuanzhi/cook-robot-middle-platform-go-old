package model

import "github.com/google/uuid"

type DBSeasoning struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt int64     `gorm:"autoCreateTime"`
	UpdatedAt int64     `gorm:"autoCreateTime"`
	UUID      uuid.UUID `json:"uuid" form:"uuid"`
	Name      string    `json:"name" form:"name"`
	Pump      uint32    `json:"pump" form:"pump"`
	Ratio     uint32    `json:"ratio" form:"ratio"`
}

func (DBSeasoning) TableName() string {
	return "seasonings"
}
