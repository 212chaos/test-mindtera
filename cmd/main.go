package main

import (
	"sync"

	"github.com/mindtera/corporate-service/common"
	router "github.com/mindtera/go-common-module/common/v2/configuration/gin/router"
)

func main() {
	// call dependencies injection
	c, grpcServer, err := BuildInRuntime()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	// go routine for grpc server
	go func() {
		grpcServer.SERVE()
		wg.Done()
	}()

	// go routine for Gin gonic server
	go func() {
		c.SERVE(router.WithPort(common.SERVICE_PORT))
		wg.Done()
	}()
	wg.Wait()
}
