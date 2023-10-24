//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

type Circuit struct {
	CID      string
	Provider string
	Type     string
	Status   string
	MetaData *MetaData
}

func NewCircuit(CID string, provider string, cType string, status string) *Circuit {
	return &Circuit{
		CID:      CID,
		Provider: provider,
		Type:     cType,
		Status:   status,
		MetaData: NewMetaData(),
	}
}

func (c *Circuit) ToProto() *octopuspb.Circuit {
	ret := &octopuspb.Circuit{
		Cid:      c.CID,
		Provider: c.Provider,
		Type:     c.Type,
		Status:   c.Status,
		MetaData: c.MetaData.ToProto(),
	}
	return ret
}
