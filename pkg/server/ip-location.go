// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-core-stack/location-services/api"
	"github.com/go-core-stack/location-services/pkg/ipinfo"
)

type IpLocationServer struct {
	api.UnimplementedIpLocationServer
	ipInfoClient *ipinfo.IPInfoClient
}

func (s *IpLocationServer) GetLocation(ctx context.Context, req *api.LocationReq) (*api.LocationResp, error) {
	log.Printf("got request for IP: %s", req.Ip)
	if s.ipInfoClient == nil {
		return nil, status.Errorf(codes.Internal, "Server not initialized properly")
	}

	if req.Ip == "" {
		return nil, status.Errorf(codes.InvalidArgument, "IP address is required")
	}

	data, err := s.ipInfoClient.GetIPInfo(req.Ip)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get IP info: %v", err)
	}

	resp := &api.LocationResp{
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		Country:   data.Country,
		Region:    data.Region,
		City:      data.City,
		Postal:    data.Postal,
	}

	return resp, nil
}

func NewIpLocationServer() *IpLocationServer {
	return &IpLocationServer{
		ipInfoClient: ipinfo.NewClient(),
	}
}
