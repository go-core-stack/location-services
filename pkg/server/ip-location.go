// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-core-stack/location-services/api"
	"github.com/go-core-stack/location-services/pkg/ipinfo"
	"github.com/go-core-stack/location-services/pkg/table"
)

var (
	// Cache Timeout
	CacheTimeOut = int64(86400)
)

type IpLocationServer struct {
	api.UnimplementedIpLocationServer
	ipInfoClient *ipinfo.IPInfoClient
	ipLocTable   *table.IpLocationTable
}

func (s *IpLocationServer) GetLocation(ctx context.Context, req *api.LocationReq) (*api.LocationResp, error) {
	if s.ipInfoClient == nil {
		return nil, status.Errorf(codes.Internal, "Server not initialized properly")
	}

	if req.Ip == "" {
		return nil, status.Errorf(codes.InvalidArgument, "IP address is required")
	}

	ipKey := &table.IpKey{
		Ip: req.Ip,
	}

	now := time.Now().Unix()
	var resp *api.LocationResp
	entry, err := s.ipLocTable.Find(ctx, ipKey)
	if err != nil || entry == nil || entry.Location == nil || (entry.Updated+CacheTimeOut < now) {
		log.Printf("performing request for IP: %s", req.Ip)
		data, err := s.ipInfoClient.GetIPInfo(req.Ip)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get IP info: %v", err)
		}

		resp = &api.LocationResp{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
			Country:   data.Country,
			Region:    data.Region,
			City:      data.City,
			Postal:    data.Postal,
		}

		update := &table.IpLocationEntry{
			Key:     ipKey,
			Updated: now,
			Location: &table.LocationInfo{
				Latitude:  data.Latitude,
				Longitude: data.Longitude,
				Country:   data.Country,
				Region:    data.Region,
				City:      data.City,
				Postal:    data.Postal,
			},
		}
		err = s.ipLocTable.Locate(ctx, update.Key, update)
		if err != nil {
			log.Printf("failed to locate ip location entry: %s", err)
		}
	} else {
		resp = &api.LocationResp{
			Latitude:  entry.Location.Latitude,
			Longitude: entry.Location.Longitude,
			Country:   entry.Location.Country,
			Region:    entry.Location.Region,
			City:      entry.Location.City,
			Postal:    entry.Location.Postal,
		}
	}
	return resp, nil
}

func NewIpLocationServer() *IpLocationServer {
	ipLocTable, err := table.GetIpLocationTable()
	if err != nil {
		log.Panicf("failed to get ip location table: %s", err)
	}
	return &IpLocationServer{
		ipInfoClient: ipinfo.NewClient(),
		ipLocTable:   ipLocTable,
	}
}
