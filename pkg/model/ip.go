//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	octopuspb "github.com/cloudflare/octopus/proto/octopus"

	bnet "github.com/bio-routing/bio-rd/net"
)

// Our IP data type which will get more attributes in the future
type IP struct {
	Address  bnet.Prefix
	MetaData *MetaData
}

func NewIP(ip bnet.Prefix) IP {
	return IP{
		Address:  ip,
		MetaData: NewMetaData(),
	}
}

func (ip *IP) ToProto() *octopuspb.IPAddress {
	if ip == nil {
		return nil
	}

	return &octopuspb.IPAddress{
		IP:       ip.Address.ToProto(),
		MetaData: ip.MetaData.ToProto(),
	}
}

func appendIPIfNotExists(slice []IP, newIP IP) []IP {
	for _, ip := range slice {
		if ip.Address.Equal(&newIP.Address) {
			return slice
		}
	}

	return append(slice, newIP)
}
