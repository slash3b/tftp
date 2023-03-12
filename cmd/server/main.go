package main

import (
	"fmt"
	"tftp/server"
)

func main() {
	fmt.Println("hello")

	srv := &server.Server{}
	srv.Listen("localhost:8080")
}
