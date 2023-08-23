package grpc

import (
	"context"
	pb "cook-robot-middle-platform-go/grpc/updater"
	"cook-robot-middle-platform-go/info"
	"cook-robot-middle-platform-go/logger"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type UpdaterGRPCClient struct {
	host string
	port uint16

	Client pb.UpdateClient
}

func NewUpdaterGRPCClient(host string, port uint16) *UpdaterGRPCClient {
	return &UpdaterGRPCClient{
		host: host,
		port: port,
	}
}

func (u *UpdaterGRPCClient) Run() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", u.host, u.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Println(err)
		return
	}
	//defer conn.Close()
	u.Client = pb.NewUpdateClient(conn)
	logger.Log.Printf("updaterGRPC客户端启动，目标地址：%s，端口：%d", u.host, u.port)
}

func (u *UpdaterGRPCClient) Check() (*pb.CheckResponse, error) {
	req := &pb.CheckRequest{
		Version:      info.Software.Version,
		MachineModel: info.Software.MachineModel,
	}
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := u.Client.Check(ctxGRPC, req)
	if err != nil {
		logger.Log.Printf("gRPC调用失败: %v", err)
		return nil, err
	}
	return res, nil
}
