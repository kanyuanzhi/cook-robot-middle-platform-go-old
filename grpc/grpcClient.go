package grpc

import (
	"context"
	pb "cook-robot-middle-platform-go/grpc/commandRPC" // 替换为你的实际包路径
	"cook-robot-middle-platform-go/logger"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ControllerStatus struct {
	CurrentCommandName string `json:"currentCommandName"`
	IsPausing          bool   `json:"isPausing"`
}

type GRPCClient struct {
	Client pb.CommandServiceClient

	ControllerStatus ControllerStatus
}

func NewGRPCClient() *GRPCClient {
	return &GRPCClient{}
}

func (g *GRPCClient) Run() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Println(err)
		return
	}
	//defer conn.Close()
	g.Client = pb.NewCommandServiceClient(conn)

	ticker := time.NewTicker(250 * time.Millisecond)
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
		logger.Log.Printf("gRPC调用失败: %v", err)
	}

	var controllerStatus ControllerStatus
	err = json.Unmarshal([]byte(res.GetStatusJson()), &controllerStatus)
	if err != nil {
		logger.Log.Printf("无法解析命令JSON：%v", err)
	}
	g.ControllerStatus = controllerStatus
	//logger.Log.Println(controllerStatus)
}
