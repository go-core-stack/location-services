// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package ipinfo

import (
	"encoding/json"
	"fmt"
)

// IPInfoResponse
type IPInfoResponse struct {
	IP        string  `json:"ip,omitempty"`
	City      string  `json:"city,omitempty"`
	Region    string  `json:"region,omitempty"`
	Country   string  `json:"country,omitempty"`
	LocStr    string  `json:"loc,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Org       string  `json:"org,omitempty"`
	Postal    string  `json:"postal,omitempty"`
	Timezone  string  `json:"timezone,omitempty"`
}

func (i *IPInfoResponse) UnmarshalJSON(b []byte) error {
	type Alias IPInfoResponse
	raw := &Alias{}
	if err := json.Unmarshal(b, raw); err != nil {
		return err
	}

	*i = IPInfoResponse(*raw)

	if raw.LocStr != "" {
		var lat, lon float64
		_, err := fmt.Sscanf(raw.LocStr, "%f,%f", &lat, &lon)
		if err == nil {
			i.Latitude = lat
			i.Longitude = lon
		}
	}
	return nil
}
