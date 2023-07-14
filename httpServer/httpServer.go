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
	router := gin.Default()
	router.Use(middleware.Cors())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &HTTPServer{
		router:     router,
		grpcClient: grpcClient,
	}
}

func (h *HTTPServer) Run() {
	gin.SetMode(gin.ReleaseMode)

	dish := v1.NewDish()

	apiV1 := h.router.Group("/api/v1")
	{
		apiV1.POST("/dish", dish.Create)
		apiV1.GET("/dish", dish.Get)
	}
	//
	//	//var instructions []instruction.Instructioner
	//	//instructions = append(instructions, instruction.NewResetXYInstruction())
	//	//
	//	//for _, step := range dish.Steps {
	//	//	instructionType := instruction.InstructionType(step.(map[string]interface{})["type"].(string))
	//	//	logger.Log.Println(instructionType)
	//	//	var instructionStruct instruction.Instructioner
	//	//	if instructionType != instruction.SEASONING {
	//	//		instructionStruct = instruction.InstructionTypeToStruct[instructionType]
	//	//		err := mapstructure.Decode(step, &instructionStruct)
	//	//		if err != nil {
	//	//			logger.Log.Println(err)
	//	//		}
	//	//	} else {
	//	//		pumpToWeightMap := map[string]uint32{}
	//	//		for _, seasoning := range step.(map[string]interface{})["seasonings"].([]interface{}) {
	//	//			//logger.Log.Println(seasoning)
	//	//			pumpNumber := fmt.Sprintf("%.0f", seasoning.(map[string]interface{})["pump"].(float64))
	//	//			pumpToWeightMap[pumpNumber] = uint32(seasoning.(map[string]interface{})["weight"].(float64))
	//	//		}
	//	//		instructionStruct = instruction.NewSeasoningInstruction(pumpToWeightMap)
	//	//	}
	//	//	instructions = append(instructions, instructionStruct)
	//	//	logger.Log.Println(instructionStruct)
	//	//}
	//	//instructions = append(instructions, instruction.NewResetRTInstruction())
	//	//
	//	//command := command.Command{CommandType: command.Multiple, Instructions: instructions}
	//	//commandJSON, err := json.Marshal(command)
	//	//if err != nil {
	//	//	logger.Log.Printf("无法将动作转换为JSON: %v", err)
	//	//}
	//	//
	//	//req := &pb.CommandRequest{
	//	//	CommandJson: string(commandJSON),
	//	//}
	//	//
	//	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//	//defer cancel()
	//	//res, err := h.grpcClient.Client.Execute(ctx, req)
	//	//if err != nil {
	//	//	logger.Log.Printf("gRPC调用失败: %v", err)
	//	//}
	//	//
	//	//logger.Log.Printf("结果X：%d\n", res.GetResult())
	//	c.JSON(200, gin.H{
	//		"message": "Hello, World!",
	//	})
	//})

	err := h.router.Run(":8889")
	if err != nil {
		return
	}
}
