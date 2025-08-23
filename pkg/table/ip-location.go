// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package table

import (
	"github.com/go-core-stack/core/db"
	"github.com/go-core-stack/core/errors"
	"github.com/go-core-stack/core/table"
)

var ipLocationTable *IpLocationTable

type IpKey struct {
	// Ip as a key for the Ip address v4 or v6
	Ip string `bson:"ip,omitempty"`
}

type LocationInfo struct {
	Latitude  float64 `bson:"latitude,omitempty"`
	Longitude float64 `bson:"longitude,omitempty"`
	Country   string  `bson:"country,omitempty"`
	Region    string  `bson:"region,omitempty"`
	City      string  `bson:"city,omitempty"`
	Postal    string  `bson:"postal,omitempty"`
}

type IpLocationEntry struct {
	// Ip Key information
	Key *IpKey `bson:"key,omitempty"`

	// last update to the entry
	Updated int64 `bson:"updated,omitempty"`

	// location info
	Location *LocationInfo `bson:"location,omitempty"`
}

type IpLocationTable struct {
	table.Table[IpKey, IpLocationEntry]
	col db.StoreCollection
}

func GetIpLocationTable() (*IpLocationTable, error) {
	if ipLocationTable != nil {
		return ipLocationTable, nil
	}

	return nil, errors.Wrapf(errors.NotFound, "tenant table not found")
}

func LocateIpLocationTable(client db.StoreClient) (*IpLocationTable, error) {
	if ipLocationTable != nil {
		return ipLocationTable, nil
	}

	col := client.GetCollection(LocationServicesDatabaseName, IpLocationCollectionName)
	tbl := &IpLocationTable{
		col: col,
	}

	err := tbl.Initialize(col)
	if err != nil {
		return nil, err
	}

	ipLocationTable = tbl

	return ipLocationTable, nil
}
