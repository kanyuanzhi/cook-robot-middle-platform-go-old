package v2

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/modelV2"
	"cook-robot-middle-platform-go/logger"
	"github.com/gin-gonic/gin"
)

type Seasoning struct {
}

func NewSeasoning() *Seasoning {
	return &Seasoning{}
}

func (s *Seasoning) ListName(ctx *gin.Context) {
	var seasonings []modelV2.Seasoning

	db.Postgres.Order("pump").Find(&seasonings)

	seasoningMap := map[uint32]string{}

	for _, seasoning := range seasonings {
		seasoningMap[seasoning.Pump] = seasoning.Name
	}

	modelV2.NewSuccessResponse(ctx, seasoningMap)
}

func (s *Seasoning) ListNameAndRatio(ctx *gin.Context) {
	var seasonings []modelV2.Seasoning
	db.Postgres.Select("uuid", "name", "pump", "ratio").Order("pump").Find(&seasonings)

	modelV2.NewSuccessResponse(ctx, seasonings)
}

type UpdateRatioReq struct {
	UUIDToRatio map[string]uint32 `json:"uuidToRatio" form:"uuidToRatio"`
}

func (s *Seasoning) UpdateRatio(ctx *gin.Context) {
	var req UpdateRatioReq
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	for uid := range req.UUIDToRatio {
		err := db.Postgres.Model(&modelV2.Seasoning{}).Where("uuid = ?", uid).Update("ratio", req.UUIDToRatio[uid]).Error
		if err != nil {
			logger.Log.Println(err)
			modelV2.NewFailResponse(ctx, err.Error())
			return
		}
	}
	modelV2.NewSuccessResponse(ctx, nil)
}
