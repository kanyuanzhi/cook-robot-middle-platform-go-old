package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/logger"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
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
	imagePath := "./assets/test.png"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		logger.Log.Println("无法读取图片文件:", err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	uid := uuid.New()
	dbDish := model.DBDish{
		UUID:    uid,
		Name:    dish.Name,
		Steps:   string(stepsStr),
		Cuisine: dish.Cuisine,
		Image:   imageData,
	}
	err = db.SQLiteDB.Create(&dbDish).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	dbCustomDishes := []model.DBCustomDish{
		{
			UUID:     uuid.New(),
			DishUUID: uid,
			Steps:    string(stepsStr),
		}, {
			UUID:     uuid.New(),
			DishUUID: uid,
			Steps:    string(stepsStr),
		}, {
			UUID:     uuid.New(),
			DishUUID: uid,
			Steps:    string(stepsStr),
		},
	}
	err = db.SQLiteDB.Create(&dbCustomDishes).Error
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
	dish.Image = base64.StdEncoding.EncodeToString(dbDish.Image)
	dish.CreatedAt = dbDish.CreatedAt
	dish.UpdatedAt = dbDish.UpdatedAt

	model.NewSuccessResponse(ctx, dish)
}

func (d *Dish) Delete(ctx *gin.Context) {
	var dish model.Dish
	if err := ctx.BindQuery(&dish); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}

	var dbDish model.DBDish
	err := db.SQLiteDB.Where("uuid = ?", dish.UUID).Delete(&dbDish).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dbCustomDish model.DBCustomDish

	err = db.SQLiteDB.Where("dish_uuid = ?", dish.UUID).Delete(&dbCustomDish).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	model.NewSuccessResponse(ctx, dish)
}

func (d *Dish) Update(ctx *gin.Context) {
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
	dbDish := model.DBDish{
		Name:    dish.Name,
		Steps:   string(stepsStr),
		Cuisine: dish.Cuisine,
	}
	err = db.SQLiteDB.Model(&dbDish).Where("uuid = ?", dish.UUID).Updates(dbDish).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	err = db.SQLiteDB.Model(&model.DBCustomDish{}).Where("dish_uuid = ?", dish.UUID).Update("steps", dbDish.Steps).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	model.NewSuccessResponse(ctx, nil)
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
			Image:     base64.StdEncoding.EncodeToString(dbDish.Image),
		})
	}
	model.NewSuccessResponse(ctx, map[string]interface{}{
		"dishes": dishes,
		"count":  count})
}
