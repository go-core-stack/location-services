// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-core-stack/core/db"
	"github.com/go-core-stack/core/values"
	"github.com/go-core-stack/location-services/api"
	"github.com/go-core-stack/location-services/pkg/config"
	"github.com/go-core-stack/location-services/pkg/server"
	"github.com/go-core-stack/location-services/pkg/table"
	"google.golang.org/grpc"
)

var (
	// path to config file
	configFile string

	// Port serving Location services
	GrpcPort = ":8080"
)

func evaluatePorts() {
	port, ok := os.LookupEnv("GRPC_PORT")
	if ok {
		GrpcPort = ":" + port
	}
}

// Parse flags for the process
func parseFlags() {
	// Add String variable flag "-config" allowing option to specify
	// the relevant config file for the process
	flag.StringVar(&configFile, "config", "", "path to the config file")

	// parse the supplied flags
	flag.Parse()
}

func main() {
	// evaluate ports to be used from ENV variables
	// if override ports are provided
	evaluatePorts()

	// Parse the flag options for the process
	parseFlags()
	conf, err := config.ParseConfig(configFile)
	if err != nil {
		log.Panicf("Failed to parse config: %s", err)
	}

	log.Printf("Got Uri config %s", conf.GetConfigDB().Uri)

	// Get mongo configdb database Credentials from environment variables
	// this is done to ensure that the credentials are not stored in plain
	// text as part of the config files
	username, password := values.GetMongoConfigDBCredentials()

	// read the configuration for configdb
	config := &db.MongoConfig{
		Uri:      conf.GetConfigDB().Uri,
		Username: username,
		Password: password,
	}

	// create new client for the mongodb config
	client, err := db.NewMongoClient(config)
	if err != nil {
		log.Panicf("Failed to get handle of mongodb client: %s", err)
	}

	// ensure running heath check to validate that provided mongodb endpoint
	// is usable
	err = client.HealthCheck(context.Background())
	if err != nil {
		log.Panicf("failed to perform Health check with DB Error: %s", err)
	}

	// locate ip location table
	_, err = table.LocateIpLocationTable(client)
	if err != nil {
		log.Panicf("failed to locate ip location table: %s", err)
	}

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
