package main

import (
	"cook-robot-middle-platform-go/config"
	"cook-robot-middle-platform-go/grpc"
	"cook-robot-middle-platform-go/httpServer"
	"fmt"
	"github.com/sstallion/go-hid"
	"log"
	"time"
)

func main() {
	grpcClient := grpc.NewGRPCClient(config.App.GRPC.TargetHost, config.App.GRPC.TargetPort)
	go grpcClient.Run()

	httpSever := httpServer.NewHTTPServer(config.App.HTTP.Host, config.App.HTTP.Port, grpcClient)
	go httpSever.Run()

	//seasonings := []*model.DBSeasoning{&model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "食用油",
	//	Pump: 1,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "酱油",
	//	Pump: 2,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "老抽",
	//	Pump: 3,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "醋",
	//	Pump: 4,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "料酒",
	//	Pump: 5,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "待定（液体）",
	//	Pump: 6,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "水",
	//	Pump: 7,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "盐",
	//	Pump: 9,
	//}, &model.DBSeasoning{
	//	UUID: uuid.New(),
	//	Name: "待定（固体）",
	//	Pump: 10,
	//}}
	//
	//db.SQLiteDB.Create(seasonings)

	//cuisines := []*model.DBCuisine{&model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "其他",
	//	Index: 0,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "川菜",
	//	Index: 1,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "湘菜",
	//	Index: 2,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "粤菜",
	//	Index: 3,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "闽菜",
	//	Index: 4,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "浙菜",
	//	Index: 5,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "鲁菜",
	//	Index: 6,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "徽菜",
	//	Index: 7,
	//}, &model.DBCuisine{
	//	UUID:  uuid.New(),
	//	Name:  "苏菜",
	//	Index: 8,
	//}}
	//db.SQLiteDB.Create(cuisines)

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

	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		fmt.Printf("%s: ID %04x:%04x %s %s\n",
			info.Path,
			info.VendorID,
			info.ProductID,
			info.MfrStr,
			info.ProductStr)
		return nil
	})

	device, err := hid.OpenFirst(0xac90, 0x3002)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(device)
	for {
		buffer := make([]byte, 1000)
		_, err = device.Read(buffer)
		if err != nil {
			fmt.Println(123, err)
		}
		fmt.Println(123213, string(buffer))
	}

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
