package v1

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/model"
	"github.com/gin-gonic/gin"
)

type Cuisine struct {
}

func NewCuisine() *Cuisine {
	return &Cuisine{}
}

func (cui *Cuisine) List(c *gin.Context) {
	var dbCuisines []model.DBCuisine

	db.SQLiteDB.Order("sort_id").Find(&dbCuisines)

	var cuisines []model.Cuisine
	for _, dbCuisine := range dbCuisines {
		cuisines = append(cuisines, model.Cuisine{
			ID:        dbCuisine.ID,
			CreatedAt: dbCuisine.CreatedAt,
			UpdatedAt: dbCuisine.UpdatedAt,
			UUID:      dbCuisine.UUID.String(),
			Name:      dbCuisine.Name,
			SortID:    dbCuisine.SortID,
		})
	}

	model.NewSuccessResponse(c, cuisines)
}
