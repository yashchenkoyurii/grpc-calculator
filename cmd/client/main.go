package main

import (
	"context"
	"fmt"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial(":5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := calculatorpb.NewCalculatorClient(cc)

	fmt.Println(Sum(client))

	fmt.Println("============")

	Prime(client)

	ComputeAverage(client)
}

func Sum(client calculatorpb.CalculatorClient) int32 {
	res, err := client.Sum(context.Background(), &calculatorpb.SumRequest{
		A: 10,
		B: 14,
	})

	if err != nil {
		log.Fatal(err)
	}

	return res.GetC()
}

func Prime(client calculatorpb.CalculatorClient) {
	res, err := client.PrimeDecomposition(context.Background(), &calculatorpb.PrimeRequest{A: 1})
	if err != nil {
		log.Fatal(err)
	}

	for {
		stream, err := res.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(stream.GetB())
	}
}

func ComputeAverage(client calculatorpb.CalculatorClient) {
	requests := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{
			A: 150,
		},
		&calculatorpb.AverageRequest{
			A: 32,
		},
		&calculatorpb.AverageRequest{
			A: 11,
		},
		&calculatorpb.AverageRequest{
			A: 24,
		},
	}

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, request := range requests {
		stream.Send(request)
		fmt.Println(request.GetA())
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Result: ")
	fmt.Println(res.GetResult())
}
