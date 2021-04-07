package main

import (
	"context"
	"fmt"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
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
	res, err := client.PrimeDecomposition(context.Background(), &calculatorpb.PrimeRequest{A: 120})
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
