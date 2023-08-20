package v2

import (
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/httpServer/modelV2"
	"github.com/gin-gonic/gin"
)

type Cuisine struct {
}

func NewCuisine() *Cuisine {
	return &Cuisine{}
}

func (cui *Cuisine) List(c *gin.Context) {
	var cuisines []modelV2.Cuisine

	db.Postgres.Order("sort_id").Find(&cuisines)

	modelV2.NewSuccessResponse(c, cuisines)
}
