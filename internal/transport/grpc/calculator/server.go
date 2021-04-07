package calculator

import (
	"context"
	"github.com/yashchenkoyurii/sum-api/internal/service"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"time"
)

type Server struct {
	calculator service.ICalculator
}

func (s Server) PrimeDecomposition(
	request *calculatorpb.PrimeRequest,
	server calculatorpb.Calculator_PrimeDecompositionServer,
) error {
	prime := s.calculator.PrimeDecomposition(request.GetA())

	for p := range prime {
		server.Send(&calculatorpb.PrimeResponse{
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
