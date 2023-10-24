//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import octopuspb "github.com/cloudflare/octopus/proto/octopus"

type FrontPort struct {
	Name             string
	RearPort         string
	RearPortPosition uint32
}

func (fp *FrontPort) ToProto() *octopuspb.FrontPort {
	return &octopuspb.FrontPort{
		Name:             fp.Name,
		RearPort:         fp.RearPort,
		RearPortPosition: fp.RearPortPosition,
	}
}
