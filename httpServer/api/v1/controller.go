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
	"strconv"
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
	//logger.Log.Println(commandReq)
	var commandStruct command.Command
	if commandReq.CommandType == command.COMMAND_TYPE_MULTIPLE {
		// 多指令
		if commandReq.CommandName == command.COMMAND_NAME_COOK {
			var dbDish model.DBDish
			err := db.SQLiteDB.First(&dbDish, "uuid = ?", commandReq.CommandData).Error
			if err != nil {
				logger.Log.Println(err)
				model.NewFailResponse(ctx, err.Error())
				return
			}
			//logger.Log.Println(dbDish)
			var stepsJSON []map[string]interface{}
			err = json.Unmarshal([]byte(dbDish.Steps), &stepsJSON)
			if err != nil {
				model.NewFailResponse(ctx, err.Error())
				return
			}
			var instructions []instruction.Instructioner

			// 开始先启动转动
			instructions = append(instructions, instruction.NewRotateInstruction("转动自启动中", "start", 1, 350, 0))

			for _, step := range stepsJSON {
				instructionType := instruction.InstructionType(step["instructionType"].(string))
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
					instructionStruct = instruction.NewSeasoningInstruction(step["instructionName"].(string), pumpToWeightMap)
				}
				instructions = append(instructions, instructionStruct)
			}

			instructions = append(instructions, instruction.NewResetRTInstruction("停止中"))

			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_COOK,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				DishUuid:     commandReq.CommandData,
				Instructions: instructions,
			}

		} else if commandReq.CommandName == command.COMMAND_NAME_PREPARE {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPrepareInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_PREPARE,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.COMMAND_NAME_DISH_OUT {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDishOutInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_DISH_OUT,
				CommandType:  command.COMMAND_TYPE_MULTIPLE,
				Instructions: instructions,
			}
		}

	} else {
		// 单指令，立即执行
		if commandReq.CommandName == command.COMMAND_NAME_DOOR_UNLOCK {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewDoorUnlockInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_DOOR_UNLOCK,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.COMMAND_NAME_PAUSE_TO_ADD {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewPauseToAddInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_PAUSE_TO_ADD,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.COMMAND_NAME_RESUME {
			var instructions []instruction.Instructioner
			instructions = append(instructions, instruction.NewResumeInstruction())
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_RESUME,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else if commandReq.CommandName == command.COMMAND_NAME_HEAT {
			var instructions []instruction.Instructioner
			temperature, err := strconv.ParseFloat(commandReq.CommandData, 10)
			if err != nil {
				logger.Log.Println("无法将字符串转换为uint32")
				return
			}
			instructions = append(instructions, instruction.NewHeatInstruction(temperature, 0, 0, instruction.NO_JUDGE))
			commandStruct = command.Command{
				CommandName:  command.COMMAND_NAME_HEAT,
				CommandType:  command.COMMAND_TYPE_SINGLE,
				Instructions: instructions,
			}
		} else {
			logger.Log.Println("命令名称错误")
			return
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
