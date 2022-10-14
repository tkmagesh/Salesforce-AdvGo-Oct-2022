package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"io"
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

	//request & response
	//doRequestResponse(ctx, service)

	//server streaming
	doServerStreaming(ctx, service)

}

func doRequestResponse(ctx context.Context, service proto.AppServiceClient) {
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

func doServerStreaming(ctx context.Context, service proto.AppServiceClient) {
	primeReq := &proto.PrimeRequest{
		Start: 3,
		End:   100,
	}
	clientStream, err := service.GeneratePrimes(ctx, primeReq)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		res, err := clientStream.Recv()
		if err == io.EOF {
			fmt.Println("All prime numbers generated have been received")
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Prime No : ", res.GetPrimeNo())
	}
}
