//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	bnet "github.com/bio-routing/bio-rd/net"
	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

type Prefix struct {
	Prefix   bnet.Prefix
	Tags     []string
	MetaData *MetaData
}

func NewPrefix(pfx bnet.Prefix) *Prefix {
	return &Prefix{
		Prefix:   pfx,
		MetaData: NewMetaData(),
	}
}

func (p *Prefix) ToProto() *octopuspb.Prefix {
	return &octopuspb.Prefix{
		Prefix:   p.Prefix.ToProto(),
		MetaData: p.MetaData.ToProto(),
	}
}
