// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package ipinfo

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-core-stack/core/errors"
)

type IPInfoClient struct {
	endpoint string
	client   *http.Client
}

func (c *IPInfoClient) GetIPInfo(ip string) (*IPInfoResponse, error) {
	ipEntry := net.ParseIP(ip)

	if ipEntry == nil {
		return nil, errors.Wrapf(errors.InvalidArgument, "invalid IP address: %s", ip)
	}

	if ipEntry.IsPrivate() {
		return nil, errors.Wrapf(errors.InvalidArgument, "private IP address: %s", ip)
	}

	url := c.endpoint + "/" + ip
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errors.Unknown, "failed to get IP Info, got response: %s", resp.Status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	ipInfo := &IPInfoResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ipInfo); err != nil {
		return nil, errors.Wrapf(errors.Unknown, "failed to decode IP Info response: %v", err)
	}
	return ipInfo, nil
}

func NewClient() *IPInfoClient {
	return &IPInfoClient{
		endpoint: "https://ipinfo.io",
		client:   &http.Client{},
	}
}
