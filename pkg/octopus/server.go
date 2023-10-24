//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package octopus

import (
	"context"

	api "github.com/cloudflare/octopus/proto/octopus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ocotopusServer struct {
	octopus *Octopus
}

func newOctopusServer(octopus *Octopus) *ocotopusServer {
	return &ocotopusServer{
		octopus: octopus,
	}
}

func (os *ocotopusServer) GetTopology(context.Context, *api.TopologyRequest) (*api.TopologyResponse, error) {
	topology := os.octopus.GetTopology()
	if topology == nil {
		return nil, status.New(codes.Unavailable, "Octopus not ready.").Err()
	}

	return &api.TopologyResponse{
		Topology: topology.ToProto(),
	}, nil
}
