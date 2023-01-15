package main

import (
	"fmt"
	"os"

	"github.com/hizani/crud_service/http_server"
)

const USAGE string = "[ADDR]:PORT [ADDR]:PORT"

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s %s\n", os.Args[0], USAGE)
		return
	}

	srv := http_server.New(os.Args[1], os.Args[2])
	srv.Start()
}
