package main

import "github.com/yashchenkoyurii/sum-api/internal/transport/grpc/app"

func main() {
	application := &app.App{}

	application.Run()
}
