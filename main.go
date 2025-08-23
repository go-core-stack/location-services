// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-core-stack/location-services/api"
	"github.com/go-core-stack/location-services/pkg/server"
	"google.golang.org/grpc"
)

var (
	// Port serving Location services
	GrpcPort = ":8080"
)

func evaluatePorts() {
	port, ok := os.LookupEnv("GRPC_PORT")
	if ok {
		GrpcPort = ":" + port
	}
}

func main() {
	evaluatePorts()

	grpcServer := grpc.NewServer()
	api.RegisterIpLocationServer(grpcServer, server.NewIpLocationServer())

	go func() {
		lis, err := net.Listen("tcp", GrpcPort)
		if err != nil {
			log.Panicf("failed to start GRPC Server")
		}
		log.Panic(grpcServer.Serve(lis))
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	s := <-sigc
	log.Printf("Terminating Process got signal: %s", s)
}
