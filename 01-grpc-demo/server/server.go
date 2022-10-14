package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"log"
	"net"

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
