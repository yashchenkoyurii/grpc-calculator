package app

import (
	"fmt"
	"github.com/yashchenkoyurii/sum-api/internal/service"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator"
	"github.com/yashchenkoyurii/sum-api/internal/transport/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func (a *App) Run() {
	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	gs := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(gs, calculator.NewServer(&service.CalculatorService{}))

	go func() {
		if err = gs.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}
