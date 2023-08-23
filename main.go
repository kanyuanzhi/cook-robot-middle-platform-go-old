package main

import (
	"cook-robot-middle-platform-go/config"
	"cook-robot-middle-platform-go/grpc"
	"cook-robot-middle-platform-go/httpServer"
	"time"
)

func main() {
	controllerGRPCClient := grpc.NewControllerGRPCClient(config.App.Controller.GRPCHost, config.App.Controller.GRPCPort)
	go controllerGRPCClient.Run()

	updaterGPRCClient := grpc.NewUpdaterGRPCClient(config.App.Updater.Host, config.App.Updater.GRPCPort)
	go updaterGPRCClient.Run()

	httpSever := httpServer.NewHTTPServer(config.App.HTTP.Host, config.App.HTTP.Port, controllerGRPCClient, updaterGPRCClient)
	go httpSever.Run()

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
