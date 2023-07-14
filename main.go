package main

import (
	"cook-robot-middle-platform-go/grpc"
	"cook-robot-middle-platform-go/httpServer"
	"time"
)

func main() {
	grpcClient := grpc.NewGRPCClient()
	go grpcClient.Run()

	httpSever := httpServer.NewHTTPServer(grpcClient)
	go httpSever.Run()

	time.Sleep(30 * time.Millisecond)

	//var instructions []instruction.Instructioner
	//
	//pumpToWeightMap := map[string]uint32{}
	//pumpToWeightMap["1"] = 111
	//pumpToWeightMap["2"] = 111
	//ins1 := instruction.NewSeasoningInstruction(pumpToWeightMap)
	//instructions = append(instructions, ins1)
	//instructions = append(instructions, instruction.NewHeatInstruction(20.4, 19.2, 20, 1))
	////instructions = []instruction.Instructioner{*ins1}
	////instructions := []instruction.Instructioner{ins1, ins2, ins3}
	//
	//logger.Log.Println(instructions)
	//
	//command := command.Command{CommandType: command.Single, Instructions: instructions}
	//commandJSON, err := json.Marshal(command)
	//if err != nil {
	//	log.Fatalf("无法将动作转换为JSON: %v", err)
	//}
	//
	//req := &pb.CommandRequest{
	//	CommandJson: string(commandJSON),
	//}
	//
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//res, err := grpcClient.Client.Execute(ctx, req)
	//if err != nil {
	//	log.Fatalf("gRPC调用失败: %v", err)
	//}
	//
	//fmt.Printf("结果X：%d\n", res.GetResult())

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
