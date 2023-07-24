package v1

import (
	"context"
	"cook-robot-middle-platform-go/grpc"
	pb "cook-robot-middle-platform-go/grpc/commandRPC"
	"cook-robot-middle-platform-go/logger"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type System struct {
	grpcClient *grpc.GRPCClient
}

func NewSystem(grpcClient *grpc.GRPCClient) *System {
	return &System{
		grpcClient: grpcClient,
	}
}

func (s *System) Shutdown(ctx *gin.Context) {
	req := &pb.ShutdownRequest{
		Empty: true,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, _ := s.grpcClient.Client.Shutdown(ctxGRPC, req)
	logger.Log.Printf("controller关闭成功%d", res)
	os.Exit(1)
}

func (s *System) GetQrCode(ctx *gin.Context) {

}
