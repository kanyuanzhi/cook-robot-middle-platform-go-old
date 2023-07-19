package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"github.com/gin-gonic/gin"
)

type Seasoning struct {
}

func NewSeasoning() *Seasoning {
	return &Seasoning{}
}

func (s *Seasoning) List(ctx *gin.Context) {
	var dbSeasonings []model.DBSeasoning

	db.SQLiteDB.Order("pump").Find(&dbSeasonings)

	seasonings := map[uint32]string{}

	for _, dbSeasoning := range dbSeasonings {
		seasonings[dbSeasoning.Pump] = dbSeasoning.Name
	}

	model.NewSuccessResponse(ctx, seasonings)
}
