package main

import (
	"context"
	"log"

	"github.com/go-core-stack/location-services/pkg/client"
)

func main() {
	cl, err := client.NewIpLocationClient("localhost", "8088")
	if err != nil {
		log.Panicf("failed to create client: %v", err)
	}
	defer func() {
		_ = cl.Close()
	}()

	resp, err := cl.GetLocation(context.Background(), "49.207.201.217")

	log.Printf("got response: %+v, err: %v", resp, err)
}
