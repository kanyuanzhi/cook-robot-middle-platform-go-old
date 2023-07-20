package v1

import (
	"context"
	"cook-robot-middle-platform-go/command"
	"cook-robot-middle-platform-go/db"
	"cook-robot-middle-platform-go/grpc"
	pb "cook-robot-middle-platform-go/grpc/commandRPC"
	"cook-robot-middle-platform-go/httpServer/model"
	"cook-robot-middle-platform-go/instruction"
	"cook-robot-middle-platform-go/logger"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"time"
)

type Controller struct {
	grpcClient *grpc.GRPCClient
}

func NewController(grpcClient *grpc.GRPCClient) *Controller {
	return &Controller{
		grpcClient: grpcClient,
	}
}

type CommandReq struct {
	CommandType string `json:"commandType" form:"commandType"`
	CommandName string `json:"commandName" form:"commandName"`
	CommandData string `json:"commandData" form:"commandData"`
}

func (c *Controller) Execute(ctx *gin.Context) {
	var commandReq CommandReq
	if err := ctx.BindJSON(&commandReq); err != nil {
		model.NewFailResponse(ctx, err.Error())
		return
	}
	logger.Log.Println(commandReq)
	var commandStruct command.Command
	if commandReq.CommandType == command.MULTIPLE {
		// 多指令
		if commandReq.CommandName == command.COOK {
			var dbDish model.DBDish
			err := db.SQLiteDB.First(&dbDish, "uuid = ?", commandReq.CommandData).Error
			if err != nil {
				logger.Log.Println(err)
				model.NewFailResponse(ctx, err.Error())
				return
			}
			logger.Log.Println(dbDish)
			var stepsJSON []map[string]interface{}
			err = json.Unmarshal([]byte(dbDish.Steps), &stepsJSON)
			if err != nil {
				model.NewFailResponse(ctx, err.Error())
				return
			}
			var instructions []instruction.Instructioner

			// 开始先启动转动
			instructions = append(instructions, instruction.NewRotateInstruction("start", 1, 350, 0))

			for _, step := range stepsJSON {
				instructionType := instruction.InstructionType(step["type"].(string))
				var instructionStruct instruction.Instructioner
				if instructionType != instruction.SEASONING {
					instructionStruct = instruction.InstructionTypeToStruct[instructionType]
					err := mapstructure.Decode(step, &instructionStruct)
					if err != nil {
						logger.Log.Println(err)
					}
				} else {
					pumpToWeightMap := map[string]uint32{}
					for _, seasoning := range step["seasonings"].([]interface{}) {
						pumpNumber := fmt.Sprintf("%.0f", seasoning.(map[string]interface{})["pumpNumber"].(float64))
						pumpToWeightMap[pumpNumber] = uint32(seasoning.(map[string]interface{})["weight"].(float64))
					}
					instructionStruct = instruction.NewSeasoningInstruction(pumpToWeightMap)
				}
				instructions = append(instructions, instructionStruct)
			}

			instructions = append(instructions, instruction.NewResetRTInstruction())

			commandStruct = command.Command{
				CommandName:  command.COOK,
				CommandType:  command.MULTIPLE,
				Instructions: instructions,
			}

		} else if commandReq.CommandName == command.RESET {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewResetRTInstruction())
			instructions = append(instructions, instruction.NewResetXYInstruction())
			commandStruct = command.Command{
				CommandName:  command.RESET,
				CommandType:  command.MULTIPLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.DISH_OUT {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDishOutInstruction())
			commandStruct = command.Command{
				CommandName:  command.DISH_OUT,
				CommandType:  command.MULTIPLE,
				Instructions: instructions,
			}
		}

	} else {
		// 单指令，立即执行
		if commandReq.CommandName == command.DOOR_UNLOCK {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDoorUnlockInstruction())
			commandStruct = command.Command{
				CommandName:  command.DOOR_UNLOCK,
				CommandType:  command.SINGLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.PAUSE_TO_ADD {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPauseToAddInstruction())
			commandStruct = command.Command{
				CommandName:  command.PAUSE_TO_ADD,
				CommandType:  command.SINGLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.RESUME {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewRestartInstruction())
			commandStruct = command.Command{
				CommandName:  command.RESUME,
				CommandType:  command.SINGLE,
				Instructions: instructions,
			}
		}
	}

	commandJSON, err := json.Marshal(commandStruct)

	req := &pb.CommandRequest{
		CommandJson: string(commandJSON),
	}

	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.grpcClient.Client.Execute(ctxGRPC, req)

	if err != nil {
		logger.Log.Println("gRPC调用失败: %v", err)
		return
	}

	if res.GetResult() == 0 {
		model.NewFailResponse(ctx, "机器占用中")
		return
	}

	model.NewSuccessResponse(ctx, nil)
}

func (c *Controller) FetchStatus(ctx *gin.Context) {
	model.NewSuccessResponse(ctx, c.grpcClient.ControllerStatus)
}

func (c *Controller) Pause(ctx *gin.Context) {
	req := &pb.PauseAndResumeRequest{Empty: true}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.grpcClient.Client.Pause(ctxGRPC, req)
	if err != nil {
		logger.Log.Println("gRPC调用失败: %v", err)
		return
	}
	model.NewSuccessResponse(ctx, res)
}

func (c *Controller) Resume(ctx *gin.Context) {
	req := &pb.PauseAndResumeRequest{Empty: true}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.grpcClient.Client.Resume(ctxGRPC, req)
	if err != nil {
		logger.Log.Println("gRPC调用失败: %v", err)
		return
	}
	model.NewSuccessResponse(ctx, res)
}
