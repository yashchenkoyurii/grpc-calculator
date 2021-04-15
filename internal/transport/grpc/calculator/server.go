package calculator

import (
	"context"
	"github.com/yashchenkoyurii/sum-api/internal/service"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

type Server struct {
	calculator service.ICalculator
}

func (s Server) SquareRoot(
	ctx context.Context,
	request *calculatorpb.SquareRootRequest,
) (
	*calculatorpb.SquareRootResponse,
	error,
) {
	root, err := s.calculator.SquareRoot(request.GetNumber())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &calculatorpb.SquareRootResponse{
		SquareRoot: root,
	}, nil
}

func (s Server) FindMaximum(stream calculatorpb.Calculator_FindMaximumServer) error {
	in := make(chan int32)
	out := make(chan int32)

	go s.calculator.FindMax(in, out)
	go func() {
		for max := range out {
			stream.Send(&calculatorpb.MaximumResponse{
				Maximum: max,
			})
		}
	}()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			close(in)
			return nil
		}

		if err != nil {
			return err
		}

		in <- req.GetNumber()
	}
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
