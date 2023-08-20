package grpc

import (
	"context"
	pb "cook-robot-middle-platform-go/grpc/commandRPC" // 替换为你的实际包路径
	"cook-robot-middle-platform-go/logger"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type InstructionInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Index        int    `json:"index"`
	ActionNumber int    `json:"actionNumber"`
}

type ControllerStatus struct {
	CurrentCommandName              string           `json:"currentCommandName"`
	CurrentDishUuid                 string           `json:"currentDishUuid"`
	CurrentInstructionName          string           `json:"currentInstructionName"`
	CurrentInstructionInfo          *InstructionInfo `json:"currentInstructionInfo"`
	IsPausing                       bool             `json:"isPausing"`
	IsRunning                       bool             `json:"isRunning"`
	IsCooking                       bool             `json:"isCooking"`
	IsPausingWithMovingFinished     bool             `json:"isPausingWithMovingFinished"`
	IsPausingWithMovingBackFinished bool             `json:"isPausingWithMovingBackFinished"`
	IsPausePermitted                bool             `json:"isPausePermitted"`
	BottomTemperature               uint32           `json:"bottomTemperature"`
	InfraredTemperature             uint32           `json:"infraredTemperature"`
	Pump1LiquidWarning              uint32           `json:"pump1LiquidWarning"`
	Pump2LiquidWarning              uint32           `json:"pump2LiquidWarning"`
	Pump3LiquidWarning              uint32           `json:"pump3LiquidWarning"`
	Pump4LiquidWarning              uint32           `json:"pump4LiquidWarning"`
	Pump5LiquidWarning              uint32           `json:"pump5LiquidWarning"`
	Pump6LiquidWarning              uint32           `json:"pump6LiquidWarning"`
	CookingTime                     int64            `json:"cookingTime"`
	CurrentHeatingTemperature       uint32           `json:"currentHeatingTemperature"`
}

type GRPCClient struct {
	targetHost string
	targetPort uint16

	Client pb.CommandServiceClient

	ControllerStatus ControllerStatus
}

func NewGRPCClient(targetHost string, targetPort uint16) *GRPCClient {
	return &GRPCClient{
		targetHost: targetHost,
		targetPort: targetPort,
	}
}

func (g *GRPCClient) Run() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", g.targetHost, g.targetPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Println(err)
		return
	}
	//defer conn.Close()
	g.Client = pb.NewCommandServiceClient(conn)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go g.FetchStatus()
		}
	}
}

func (g *GRPCClient) FetchStatus() {
	req := &pb.FetchRequest{
		Empty: true,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := g.Client.FetchStatus(ctxGRPC, req)
	if err != nil {
		//logger.Log.Printf("gRPC调用失败: %v", err)
		return
	}

	var controllerStatus ControllerStatus
	err = json.Unmarshal([]byte(res.GetStatusJson()), &controllerStatus)
	if err != nil {
		logger.Log.Printf("无法解析命令JSON：%v", err)
		return
	}
	g.ControllerStatus = controllerStatus
	//logger.Log.Println(controllerStatus)
}
