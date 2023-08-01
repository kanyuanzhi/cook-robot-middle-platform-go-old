package httpServer

import (
	"cook-robot-middle-platform-go/config"
	"cook-robot-middle-platform-go/grpc"
	v1 "cook-robot-middle-platform-go/httpServer/api/v1"
	"cook-robot-middle-platform-go/httpServer/middleware"
	"cook-robot-middle-platform-go/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

type HTTPServer struct {
	host string
	port uint16

	router     *gin.Engine
	grpcClient *grpc.GRPCClient
}

func NewHTTPServer(host string, port uint16, grpcClient *grpc.GRPCClient) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.Cors())
	//router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/controller/fetchStatus"))

	return &HTTPServer{
		host:       host,
		port:       port,
		router:     router,
		grpcClient: grpcClient,
	}
}

func (h *HTTPServer) Run() {

	dish := v1.NewDish()
	//customDish := v1.NewCustomDish()
	cuisine := v1.NewCuisine()
	seasoning := v1.NewSeasoning()

	controller := v1.NewController(h.grpcClient)

	system := v1.NewSystem(h.grpcClient)

	apiV1 := h.router.Group("/api/v1")
	{
		apiV1.POST("/dish", dish.Create)
		apiV1.PUT("/dish", dish.Update)
		apiV1.GET("/dish", dish.Get)
		apiV1.DELETE("/dish", dish.Delete)
		apiV1.GET("/dishes", dish.ListByCuisine)
		apiV1.GET("/allDishes", dish.ListAll)
		apiV1.GET("/customDishes", dish.ListCustomDishes)
		apiV1.PUT("/customDishes", dish.UpdateCustomDishes)
		//apiV1.GET("/customDishes", customDish.ListByDishUUID)

		apiV1.GET("/cuisines", cuisine.List)

		apiV1.GET("/seasonings", seasoning.List)

		apiV1.POST("/controller/execute", controller.Execute)
		apiV1.GET("/controller/fetchStatus", controller.FetchStatus)
		apiV1.GET("/controller/pause", controller.Pause)
		apiV1.GET("/controller/resume", controller.Resume)

		apiV1.GET("/system/getQrCode", system.GetQrCode)
		apiV1.GET("/system/shutdown", system.Shutdown)
		apiV1.GET("/system/checkUpdatePermission", system.CheckUpdatePermission)
		apiV1.GET("/system/update", system.Update)
	}

	var err error
	if config.App.HTTP.UseSSL {
		logger.Log.Println("使用HTTPS")
		dir, _ := os.Getwd()
		cerFilePath := filepath.Join(dir, config.App.HTTP.SSLDir, config.App.HTTP.CerFile)
		keyFilePath := filepath.Join(dir, config.App.HTTP.SSLDir, config.App.HTTP.KeyFile)
		err = h.router.RunTLS(fmt.Sprintf("%s:%d", h.host, h.port), cerFilePath, keyFilePath)
	} else {
		logger.Log.Println("使用HTTP")
		err = h.router.Run(fmt.Sprintf("%s:%d", h.host, h.port))
	}
	if err != nil {
		logger.Log.Println(err)
		return
	}
}
