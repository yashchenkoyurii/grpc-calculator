package calculator

import (
	"context"
	"github.com/yashchenkoyurii/sum-api/internal/service"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"
)

type Server struct {
	calculator service.ICalculator
}

func (s Server) ComputeAverage(stream calculatorpb.Calculator_ComputeAverageServer) error {
	numbers := make(chan int32, 100)
	result := make(chan float32)
	go s.calculator.ComputeAverage(numbers, result)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			close(numbers)
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		numbers <- req.GetA()
	}

	stream.SendAndClose(&calculatorpb.AverageResponse{
		Result: <-result,
	})

	return nil
}

func (s Server) PrimeDecomposition(
	request *calculatorpb.PrimeRequest,
	stream calculatorpb.Calculator_PrimeDecompositionServer,
) error {
	prime := s.calculator.PrimeDecomposition(request.GetA())

	for p := range prime {
		stream.Send(&calculatorpb.PrimeResponse{
			B: p,
		})

		time.Sleep(time.Second)
	}

	return nil
}

func (s Server) Sum(
	ctx context.Context,
	request *calculatorpb.SumRequest,
) (*calculatorpb.SumResponse, error) {
	return &calculatorpb.SumResponse{C: s.calculator.Sum(request.GetA(), request.GetB())}, nil
}

func NewServer(calculator service.ICalculator) *Server {
	return &Server{calculator: calculator}
}
