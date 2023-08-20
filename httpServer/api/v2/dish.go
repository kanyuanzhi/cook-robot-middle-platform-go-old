package v2

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/modelV2"
	"cook-robot-middle-platform-go/logger"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	var dish modelV2.Dish
	if err := ctx.BindJSON(&dish); err != nil {
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	imagePath := "./assets/test.png"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		logger.Log.Println("无法读取图片文件:", err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	uid := uuid.New()
	dish.UUID = uid
	dish.IsTaste = false
	dish.DishUUID = uuid.Nil
	dish.Image = imageData

	dbDishesAndItsTastes := []modelV2.Dish{
		dish,
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味1）",
			Steps:    dish.Steps,
			IsTaste:  true,
			DishUUID: dish.UUID,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味2）",
			Steps:    dish.Steps,
			IsTaste:  true,
			DishUUID: dish.UUID,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味3）",
			Steps:    dish.Steps,
			IsTaste:  true,
			DishUUID: dish.UUID,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
	}
	err = db.Postgres.Create(&dbDishesAndItsTastes).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	modelV2.NewSuccessResponse(ctx, uid)
}

func (d *Dish) Get(ctx *gin.Context) {
	type Req struct {
		UUID string `json:"uuid" form:"uuid"`
	}
	var req Req
	if err := ctx.BindQuery(&req); err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	var dish modelV2.Dish
	//if err := ctx.BindQuery(&dish); err != nil {
	//	logger.Log.Println(err)
	//	modelV2.NewFailResponse(ctx, err.Error())
	//	return
	//}

	err := db.Postgres.First(&dish, "uuid = ?", req.UUID).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	modelV2.NewSuccessResponse(ctx, dish)
}

func (d *Dish) Delete(ctx *gin.Context) {
	type Req struct {
		UUID string `json:"uuid" form:"uuid"`
	}

	var req Req
	if err := ctx.BindQuery(&req); err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	var dish modelV2.Dish
	err := db.Postgres.Where("uuid = ?", req.UUID).Or("dish_uuid = ?", req.UUID).Delete(&dish).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	modelV2.NewSuccessResponse(ctx, dish)
}

func (d *Dish) Update(ctx *gin.Context) {
	var dish modelV2.Dish
	if err := ctx.BindJSON(&dish); err != nil {
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	err := db.Postgres.Model(&dish).Where("uuid = ?", dish.UUID).Updates(dish).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	var customDishes []modelV2.Dish
	err = db.Postgres.Where("dish_uuid = ?", dish.UUID).Find(&customDishes).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	for index, dbCustomDish := range customDishes {
		err = db.Postgres.Model(&dbCustomDish).Updates(modelV2.Dish{Name: dish.Name + fmt.Sprintf("（口味%d）", index+1), Steps: dish.Steps, Cuisine: dish.Cuisine}).Error
		if err != nil {
			logger.Log.Println(err)
			modelV2.NewFailResponse(ctx, err.Error())
			return
		}
	}

	modelV2.NewSuccessResponse(ctx, nil)
}

type DishQueryReq struct {
	CuisineID string `json:"cuisineID" form:"cuisineID" `
	PageSize  int    `json:"pageSize" form:"pageSize"`
	PageIndex int    `json:"pageIndex" form:"pageIndex"`
}

func (d *Dish) ListAll(ctx *gin.Context) {
	var dishQueryReq DishQueryReq
	if err := ctx.BindQuery(&dishQueryReq); err != nil {
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	var dishes []modelV2.Dish
	var count int64
	err := db.Postgres.Where("is_taste = ?", false).Limit(dishQueryReq.PageSize).Offset((dishQueryReq.PageIndex - 1) * dishQueryReq.PageSize).
		Find(&dishes).Count(&count).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	modelV2.NewSuccessResponse(ctx, map[string]interface{}{
		"dishes": dishes,
		"count":  count})
}

func (d *Dish) ListByCuisine(ctx *gin.Context) {
	var dishQueryReq DishQueryReq
	if err := ctx.BindQuery(&dishQueryReq); err != nil {
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	var dishes []modelV2.Dish
	var count int64
	err := db.Postgres.Where("cuisine = ? AND is_taste = ?", dishQueryReq.CuisineID, false).
		Limit(dishQueryReq.PageSize).Offset((dishQueryReq.PageIndex - 1) * dishQueryReq.PageSize).
		Find(&dishes).Count(&count).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	modelV2.NewSuccessResponse(ctx, map[string]interface{}{
		"dishes": dishes,
		"count":  count,
	})
}

func (d *Dish) ListCustomDishes(ctx *gin.Context) {
	var customDish modelV2.Dish
	if err := ctx.BindQuery(&customDish); err != nil {
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	var customDishes []modelV2.Dish
	err := db.Postgres.Where("dish_uuid = ?", customDish.DishUUID).Find(&customDishes).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	modelV2.NewSuccessResponse(ctx, customDishes)
}

type UpdateCustomDishesReq struct {
	UUIDToSteps map[string][]map[string]interface{} `json:"uuidToSteps" form:"uuidToSteps"`
}

func (d *Dish) UpdateCustomDishes(ctx *gin.Context) {
	var req UpdateCustomDishesReq
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	for uid := range req.UUIDToSteps {
		stepsStr, err := json.Marshal(req.UUIDToSteps[uid])
		if err != nil {
			logger.Log.Println("转换数组为字符串失败:", err)
			modelV2.NewFailResponse(ctx, err.Error())
			return
		}
		err = db.Postgres.Model(&modelV2.Dish{}).Where("uuid = ?", uid).Update("steps", stepsStr).Error
		if err != nil {
			logger.Log.Println(err)
			modelV2.NewFailResponse(ctx, err.Error())
			return
		}
	}
	modelV2.NewSuccessResponse(ctx, nil)
}

func (d *Dish) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	logger.Log.Println(file.Filename)
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}
	uid := ctx.PostForm("uuid")

	err = db.Postgres.Model(&modelV2.Dish{}).Where("uuid = ?", uid).Update("image", fileData).Error
	if err != nil {
		logger.Log.Println(err)
		modelV2.NewFailResponse(ctx, err.Error())
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)
	modelV2.NewSuccessResponse(ctx, encodedData)
}
