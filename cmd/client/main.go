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

	fmt.Println("SUM: ")
	Sum(client)
	fmt.Println("Prime: ")
	Prime(client)
	fmt.Println("Compute average: ")
	ComputeAverage(client)
	fmt.Println("FIND MAX:")
	FindMax(client)
}

func FindMax(client calculatorpb.CalculatorClient) {
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan struct{})

	numbers := []int32{1, 5, 3, 6, 2, 20, 1}

	go func() {
		for _, n := range numbers {
			fmt.Println("Request: ", n)
			stream.Send(&calculatorpb.MaximumRequest{
				Number: n,
			})
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				close(wait)
				return
			}

			if err != nil {
				log.Fatal(err)
				close(wait)
			}

			fmt.Printf("Maximum: %v \n", res.GetMaximum())
		}
	}()

	<-wait
}

func Sum(client calculatorpb.CalculatorClient) {
	res, err := client.Sum(context.Background(), &calculatorpb.SumRequest{
		A: 10,
		B: 14,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.GetC())
}

func Prime(client calculatorpb.CalculatorClient) {
	res, err := client.PrimeDecomposition(context.Background(), &calculatorpb.PrimeRequest{A: 123456})
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
