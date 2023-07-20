package httpServer

import (
	"cook-robot-middle-platform-go/grpc"
	v1 "cook-robot-middle-platform-go/httpServer/api/v1"
	"cook-robot-middle-platform-go/httpServer/middleware"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router     *gin.Engine
	grpcClient *grpc.GRPCClient
}

func NewHTTPServer(grpcClient *grpc.GRPCClient) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Cors())
	//router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/controller/fetchStatus"))

	return &HTTPServer{
		router:     router,
		grpcClient: grpcClient,
	}
}

func (h *HTTPServer) Run() {

	dish := v1.NewDish()
	customDish := v1.NewCustomDish()
	cuisine := v1.NewCuisine()
	seasoning := v1.NewSeasoning()

	controller := v1.NewController(h.grpcClient)

	apiV1 := h.router.Group("/api/v1")
	{
		apiV1.POST("/dish", dish.Create)
		apiV1.PUT("/dish", dish.Update)
		apiV1.GET("/dish", dish.Get)
		apiV1.DELETE("/dish", dish.Delete)
		apiV1.GET("/dishes", dish.ListByCuisine)

		apiV1.GET("/customDishes", customDish.ListByDishUUID)

		apiV1.GET("/cuisines", cuisine.List)

		apiV1.GET("/seasonings", seasoning.List)

		apiV1.POST("/controller/execute", controller.Execute)
		apiV1.GET("/controller/fetchStatus", controller.FetchStatus)
		apiV1.GET("/controller/pause", controller.Pause)
		apiV1.GET("/controller/resume", controller.Resume)
	}

	err := h.router.Run(":8889")
	if err != nil {
		return
	}
}
