//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import octopuspb "github.com/cloudflare/octopus/proto/octopus"

type RearPort struct {
	Name      string
	Positions int16
}

func (rp *RearPort) ToProto() *octopuspb.RearPort {
	return &octopuspb.RearPort{
		Name:      rp.Name,
		Positions: uint32(rp.Positions),
	}
}
