package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Dish struct {
}

func NewDish() *Dish {
	return &Dish{}
}

func (d *Dish) Create(ctx *gin.Context) {
	var dish model.Dish
	if err := ctx.BindJSON(&dish); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	stepsStr, err := json.Marshal(dish.Steps)
	if err != nil {
		logger.Log.Println("转换数组为字符串失败:", err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	uid := uuid.New()
	dbDish := model.DBDish{
		UUID:    uid,
		Name:    dish.Name,
		Steps:   string(stepsStr),
		Cuisine: dish.Cuisine,
	}
	err = db.SQLiteDB.Create(&dbDish).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	model.NewSuccessResponse(ctx, uid)
}

func (d *Dish) Get(ctx *gin.Context) {
	var dish model.Dish
	if err := ctx.BindQuery(&dish); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}

	var dbDish model.DBDish
	err := db.SQLiteDB.First(&dbDish, "uuid = ?", dish.UUID).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	var stepsJSON []map[string]interface{}
	err = json.Unmarshal([]byte(dbDish.Steps), &stepsJSON)
	if err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}

	dish.Name = dbDish.Name
	dish.UUID = dbDish.UUID.String()
	dish.Steps = stepsJSON
	dish.Cuisine = dbDish.Cuisine
	dish.CreatedAt = dbDish.CreatedAt
	dish.UpdatedAt = dbDish.UpdatedAt

	model.NewSuccessResponse(ctx, dish)
}

type DishQueryReq struct {
	CuisineID string `json:"cuisineID" form:"cuisineID" `
	PageSize  int    `json:"pageSize" form:"pageSize"`
	PageIndex int    `json:"pageIndex" form:"pageIndex"`
}

func (d *Dish) ListByCuisine(ctx *gin.Context) {
	var dishQueryReq DishQueryReq
	if err := ctx.BindQuery(&dishQueryReq); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dbDishes []model.DBDish
	var count int64
	err := db.SQLiteDB.Where("cuisine = ?", dishQueryReq.CuisineID).
		Limit(dishQueryReq.PageSize).Offset((dishQueryReq.PageIndex - 1) * dishQueryReq.PageSize).
		Find(&dbDishes).Count(&count).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dishes []model.Dish
	for _, dbDish := range dbDishes {
		dishes = append(dishes, model.Dish{
			CreatedAt: dbDish.CreatedAt,
			UpdatedAt: dbDish.UpdatedAt,
			UUID:      dbDish.UUID.String(),
			Name:      dbDish.Name,
		})
	}
	model.NewSuccessResponse(ctx, map[string]interface{}{
		"dishes": dishes,
		"count":  count})
}
