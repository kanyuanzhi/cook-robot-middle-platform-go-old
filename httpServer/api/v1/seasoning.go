package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/logger"
	"github.com/gin-gonic/gin"
)

type Seasoning struct {
}

func NewSeasoning() *Seasoning {
	return &Seasoning{}
}

func (s *Seasoning) ListName(ctx *gin.Context) {
	var dbSeasonings []model.DBSeasoning

	db.SQLiteDB.Order("pump").Find(&dbSeasonings)

	seasonings := map[uint32]string{}

	for _, dbSeasoning := range dbSeasonings {
		seasonings[dbSeasoning.Pump] = dbSeasoning.Name
	}

	model.NewSuccessResponse(ctx, seasonings)
}

func (s *Seasoning) ListNameAndRatio(ctx *gin.Context) {
	var dbSeasonings []model.DBSeasoning
	db.SQLiteDB.Select("uuid", "name", "pump", "ratio").Order("pump").Find(&dbSeasonings)

	var seasonings []map[string]interface{}

	for _, dbSeasoning := range dbSeasonings {
		seasonings = append(seasonings, map[string]interface{}{
			"uuid":       dbSeasoning.UUID,
			"name":       dbSeasoning.Name,
			"pumpNumber": dbSeasoning.Pump,
			"ratio":      dbSeasoning.Ratio,
		})
	}

	model.NewSuccessResponse(ctx, seasonings)
}

type UpdateRatioReq struct {
	UUIDToRatio map[string]uint32 `json:"uuidToRatio" form:"uuidToRatio"`
}

func (s *Seasoning) UpdateRatio(ctx *gin.Context) {
	var req UpdateRatioReq
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	for uid := range req.UUIDToRatio {
		err := db.SQLiteDB.Model(&model.DBSeasoning{}).Where("uuid = ?", uid).Update("ratio", req.UUIDToRatio[uid]).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
	}
	model.NewSuccessResponse(ctx, nil)
}
