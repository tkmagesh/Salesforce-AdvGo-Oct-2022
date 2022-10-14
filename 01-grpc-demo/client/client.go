package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial("localhost:50051", options)
	if err != nil {
		log.Fatalln(err)
	}
	service := proto.NewAppServiceClient(clientConn)
	ctx := context.Background()
	req := &proto.AddRequest{
		X: 100,
		Y: 200,
	}
	response, err := service.Add(ctx, req)
	if err != nil {
		log.Fatalln()
	}
	fmt.Println(response.GetResult())

}
