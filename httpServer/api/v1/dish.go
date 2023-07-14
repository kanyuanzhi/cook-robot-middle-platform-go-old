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

func (d *Dish) Create(c *gin.Context) {
	var dish model.Dish
	if err := c.BindJSON(&dish); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	stepsStr, err := json.Marshal(dish.Steps)
	if err != nil {
		logger.Log.Println("转换数组为字符串失败:", err)
	}
	uid := uuid.New()
	dbDish := model.DBDish{
		UUID:  uid,
		Name:  dish.Name,
		Steps: string(stepsStr),
	}
	result := db.SQLiteDB.Create(&dbDish)
	if result.Error != nil {
		logger.Log.Println(result.Error)
	}
	model.NewSuccessResponse(c, uid)
}

func (d *Dish) Get(c *gin.Context) {
	var dish model.Dish
	if err := c.BindQuery(&dish); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var dbDish model.DBDish
	db.SQLiteDB.First(&dbDish, "uuid = ?", dish.UUID)

	var stepsJSON []map[string]interface{}
	err := json.Unmarshal([]byte(dbDish.Steps), &stepsJSON)
	if err != nil {
		logger.Log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	dish.Name = dbDish.Name
	dish.Steps = stepsJSON
	dish.CreatedAt = dbDish.CreatedAt
	dish.UpdatedAt = dbDish.UpdatedAt

	model.NewSuccessResponse(c, dish)
}
