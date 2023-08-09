package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
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
		UUID:     uid,
		Name:     dish.Name,
		Steps:    string(stepsStr),
		IsTaste:  false,
		DishUUID: uuid.Nil,
		Cuisine:  dish.Cuisine,
		Image:    imageData,
	}
	dbDishesAndItsTastes := []model.DBDish{
		dbDish,
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味1）",
			Steps:    string(stepsStr),
			IsTaste:  true,
			DishUUID: uid,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味2）",
			Steps:    string(stepsStr),
			IsTaste:  true,
			DishUUID: uid,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
		{
			UUID:     uuid.New(),
			Name:     dish.Name + "（口味3）",
			Steps:    string(stepsStr),
			IsTaste:  true,
			DishUUID: uid,
			Cuisine:  dish.Cuisine,
			Image:    nil,
		},
	}
	err = db.SQLiteDB.Create(&dbDishesAndItsTastes).Error
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
	dish.IsTaste = dbDish.IsTaste
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
	err := db.SQLiteDB.Where("uuid = ?", dish.UUID).Or("dish_uuid = ?", dish.UUID).Delete(&dbDish).Error
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

	var dbCustomDishes []model.DBDish
	err = db.SQLiteDB.Where("dish_uuid = ?", dish.UUID).Find(&dbCustomDishes).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	for index, dbCustomDish := range dbCustomDishes {
		err = db.SQLiteDB.Model(&dbCustomDish).Updates(model.DBDish{Name: dish.Name + fmt.Sprintf("（口味%d）", index+1), Steps: string(stepsStr), Cuisine: dish.Cuisine}).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
	}

	//err = db.SQLiteDB.Model(&model.DBDish{}).Where("dish_uuid = ?", dish.UUID).Updates(model.DBDish{Name: dish.Name, Steps: string(stepsStr), Cuisine: dish.Cuisine}).Error
	//if err != nil {
	//	logger.Log.Println(err)
	//	model.NewFailResponse(ctx, err.Error())
	//	return
	//}
	model.NewSuccessResponse(ctx, nil)
}

type DishQueryReq struct {
	CuisineID string `json:"cuisineID" form:"cuisineID" `
	PageSize  int    `json:"pageSize" form:"pageSize"`
	PageIndex int    `json:"pageIndex" form:"pageIndex"`
}

func (d *Dish) ListAll(ctx *gin.Context) {
	var dishQueryReq DishQueryReq
	if err := ctx.BindQuery(&dishQueryReq); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dbDishes []model.DBDish
	var count int64
	err := db.SQLiteDB.Where("is_taste = ?", false).Limit(dishQueryReq.PageSize).Offset((dishQueryReq.PageIndex - 1) * dishQueryReq.PageSize).
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

func (d *Dish) ListByCuisine(ctx *gin.Context) {
	var dishQueryReq DishQueryReq
	if err := ctx.BindQuery(&dishQueryReq); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dbDishes []model.DBDish
	var count int64
	err := db.SQLiteDB.Where("cuisine = ? AND is_taste = ?", dishQueryReq.CuisineID, false).
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
		"count":  count,
	})
}

func (d *Dish) ListCustomDishes(ctx *gin.Context) {
	var customDish model.Dish
	if err := ctx.BindQuery(&customDish); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	var dbCustomDishes []model.DBDish
	err := db.SQLiteDB.Where("dish_uuid = ?", customDish.DishUUID).Find(&dbCustomDishes).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	if len(dbCustomDishes) == 0 { // 临时使用
		var dbDish model.DBDish
		err := db.SQLiteDB.First(&dbDish, "uuid = ?", customDish.DishUUID).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
		dbDishesAndItsTastes := []model.DBDish{
			{
				UUID:     uuid.New(),
				Name:     dbDish.Name + "（口味1）",
				Steps:    dbDish.Steps,
				IsTaste:  true,
				DishUUID: dbDish.UUID,
				Cuisine:  dbDish.Cuisine,
				Image:    nil,
			},
			{
				UUID:     uuid.New(),
				Name:     dbDish.Name + "（口味2）",
				Steps:    dbDish.Steps,
				IsTaste:  true,
				DishUUID: dbDish.UUID,
				Cuisine:  dbDish.Cuisine,
				Image:    nil,
			},
			{
				UUID:     uuid.New(),
				Name:     dbDish.Name + "（口味2）",
				Steps:    dbDish.Steps,
				IsTaste:  true,
				DishUUID: dbDish.UUID,
				Cuisine:  dbDish.Cuisine,
				Image:    nil,
			},
		}
		err = db.SQLiteDB.Create(&dbDishesAndItsTastes).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
		fmt.Println(customDish.DishUUID)
		err = db.SQLiteDB.Model(&model.DBDish{}).Where("uuid = ?", customDish.DishUUID).Update("dish_uuid", uuid.Nil).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
	}
	var customDishes []model.Dish

	for _, dbCustomDish := range dbCustomDishes {
		var stepsJSON []map[string]interface{}
		err = json.Unmarshal([]byte(dbCustomDish.Steps), &stepsJSON)
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
		customDish.Name = dbCustomDish.Name
		customDish.UUID = dbCustomDish.UUID.String()
		customDish.DishUUID = dbCustomDish.DishUUID.String()
		customDish.Steps = stepsJSON
		customDish.IsTaste = dbCustomDish.IsTaste
		customDish.Cuisine = dbCustomDish.Cuisine
		customDish.CreatedAt = dbCustomDish.CreatedAt
		customDish.UpdatedAt = dbCustomDish.UpdatedAt
		customDishes = append(customDishes, customDish)
	}

	model.NewSuccessResponse(ctx, customDishes)
}

type UpdateCustomDishesReq struct {
	UUIDToSteps map[string][]map[string]interface{} `json:"uuidToSteps" form:"uuidToSteps"`
}

func (d *Dish) UpdateCustomDishes(ctx *gin.Context) {
	var req UpdateCustomDishesReq
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	for uid := range req.UUIDToSteps {
		stepsStr, err := json.Marshal(req.UUIDToSteps[uid])
		if err != nil {
			logger.Log.Println("转换数组为字符串失败:", err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
		err = db.SQLiteDB.Model(&model.DBDish{}).Where("uuid = ?", uid).Update("steps", stepsStr).Error
		if err != nil {
			logger.Log.Println(err)
			model.NewFailResponse(ctx, err.Error())
			return
		}
	}
	model.NewSuccessResponse(ctx, nil)
}

func (d *Dish) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	logger.Log.Println(file.Filename)
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}
	uid := ctx.PostForm("uuid")

	err = db.SQLiteDB.Model(&model.DBDish{}).Where("uuid = ?", uid).Update("image", fileData).Error
	if err != nil {
		logger.Log.Println(err)
		model.NewFailResponse(ctx, err.Error())
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)
	model.NewSuccessResponse(ctx, encodedData)
}
