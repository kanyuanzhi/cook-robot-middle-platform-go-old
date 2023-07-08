package grpc

import (
	pb "cook-robot-middle-platform-go/grpc/commandRPC" // 替换为你的实际包路径
	"cook-robot-middle-platform-go/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Client pb.CommandServiceClient
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
	defer conn.Close()
	g.Client = pb.NewCommandServiceClient(conn)
}
