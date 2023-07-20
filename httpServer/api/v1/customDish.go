package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type CustomDish struct {
}

func NewCustomDish() *CustomDish {
	return &CustomDish{}
}

func (c *CustomDish) ListByDishUUID(ctx *gin.Context) {
	var customDish model.CustomDish
	if err := ctx.BindQuery(&customDish); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}

	var dbCustomDishes []model.DBCustomDish
	err := db.SQLiteDB.Where("dish_uuid = ?", customDish.DishUUID).Find(&dbCustomDishes).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	var customDishes []model.CustomDish
	var stepsJSON []map[string]interface{}

	for _, dbCustomDish := range dbCustomDishes {
		err = json.Unmarshal([]byte(dbCustomDish.Steps), &stepsJSON)
		if err != nil {
			model.NewFailResponse(ctx, err.Error())
			return
		}
		customDish.UUID = dbCustomDish.UUID.String()
		customDish.DishUUID = dbCustomDish.DishUUID.String()
		customDish.Steps = stepsJSON
		customDish.CreatedAt = dbCustomDish.CreatedAt
		customDish.UpdatedAt = dbCustomDish.UpdatedAt
		customDishes = append(customDishes, customDish)
	}

	model.NewSuccessResponse(ctx, customDishes)
}
