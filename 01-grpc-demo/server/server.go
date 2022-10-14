package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

//step : 3 - service implementation
type appService struct {
	proto.UnimplementedAppServiceServer
}

func (service appService) Add(ctx context.Context, req *proto.AddRequest) (res *proto.AddResponse, err error) {
	x := req.GetX()
	y := req.GetY()
	fmt.Printf("Add : Processing %d and %d\n", x, y)
	result := x + y
	res = &proto.AddResponse{
		Result: result,
	}
	return
}

func (service appService) GeneratePrimes(req *proto.PrimeRequest, serverStream proto.AppService_GeneratePrimesServer) error {
	start := req.GetStart()
	end := req.GetEnd()
	for no := start; no <= end; no++ {
		if isPrime(no) {
			res := &proto.PrimeResponse{
				PrimeNo: no,
			}
			fmt.Printf("GeneratePrimes : sending prime no %d\n", no)
			err := serverStream.Send(res)
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	return nil
}

func isPrime(no int32) bool {
	for i := int32(2); i <= (no / 2); i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	//step - 4 : hosting the service
	asi := &appService{}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAppServiceServer(grpcServer, asi)
	grpcServer.Serve(listener)
}
