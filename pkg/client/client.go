// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-core-stack/location-services/api"
)

type IpLocationClient struct {
	conn      *grpc.ClientConn
	locClient api.IpLocationClient
}

func (c *IpLocationClient) Close() error {
	return c.conn.Close()
}

func (c *IpLocationClient) GetLocation(ctx context.Context, ip string) (*api.LocationResp, error) {
	req := &api.LocationReq{
		Ip: ip,
	}
	return c.locClient.GetLocation(ctx, req)
}

func NewIpLocationClient(endpoint, port string) (*IpLocationClient, error) {
	conn, err := grpc.NewClient(endpoint+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &IpLocationClient{
		conn:      conn,
		locClient: api.NewIpLocationClient(conn),
	}, nil
}
