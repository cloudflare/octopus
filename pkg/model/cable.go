//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	"fmt"

	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

type Cable struct {
	AEnd CableEnd
	BEnd CableEnd
}

type CableEnd struct {
	DeviceName   string
	EndpointName string
	EndpointType octopuspb.CableEndpointType
}

func (c *Cable) ToProto() *octopuspb.Cable {
	if c == nil {
		return nil
	}

	return &octopuspb.Cable{
		AEnd: c.AEnd.toProto(),
		BEnd: c.BEnd.toProto(),
	}
}

func (ce CableEnd) toProto() *octopuspb.CableEnd {
	return &octopuspb.CableEnd{
		DeviceName:   ce.DeviceName,
		EndpointType: ce.EndpointType,
		EndpointName: ce.EndpointName,
	}
}

func (c Cable) String() string {
	return fmt.Sprintf("%s:%s:%d<->%s:%s:%d",
		c.AEnd.DeviceName, c.AEnd.EndpointName, c.AEnd.EndpointType,
		c.BEnd.DeviceName, c.BEnd.EndpointName, c.BEnd.EndpointType)
}
