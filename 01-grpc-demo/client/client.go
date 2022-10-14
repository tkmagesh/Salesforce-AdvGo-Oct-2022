package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"io"
	"log"
	"time"

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
	//doServerStreaming(ctx, service)

	//client streaming
	//doClientStreaming(ctx, service)

	//Bidirectional streaming
	doBiDiStreaming(ctx, service)

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

func doClientStreaming(ctx context.Context, service proto.AppServiceClient) {
	nos := []int32{9, 1, 8, 4, 7, 2, 6, 5, 3}
	clientStream, err := service.CalculateAverage(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, no := range nos {
		fmt.Printf("Sending no : %d\n", no)
		averageRequest := &proto.AverageRequest{
			No: no,
		}
		if err := clientStream.Send(averageRequest); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("All the inputs are sent!")
	res, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Average : %v\n", res.GetResult())
}

func doBiDiStreaming(ctx context.Context, service proto.AppServiceClient) {
	personNames := []proto.PersonName{
		{FirstName: "Magesh", LastName: "Kuppan"},
		{FirstName: "Suresh", LastName: "Kannan"},
		{FirstName: "Rajesh", LastName: "Pandit"},
		{FirstName: "Ganesh", LastName: "Easwaran"},
		{FirstName: "Ramesh", LastName: "Jayaraman"},
	}
	clientStream, err := service.Greet(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//handle responses
	done := func() <-chan struct{} {
		doneCh := make(chan struct{})
		go func() {
			for {
				res, err := clientStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalln(err)
				}
				msg := res.GetGreetMessage()
				fmt.Println(msg)
			}
			close(doneCh)
		}()
		return doneCh
	}()

	for _, personName := range personNames {
		fmt.Printf("Sending Person %v\n", personName)
		time.Sleep(500 * time.Millisecond)
		req := &proto.GreetRequest{
			Person: &personName,
		}
		err := clientStream.Send(req)
		if err != nil {
			log.Fatalln(err)
		}
	}
	clientStream.CloseSend()
	<-done
}
