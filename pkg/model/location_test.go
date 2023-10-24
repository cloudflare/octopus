//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	"testing"

	octopuspb "github.com/cloudflare/octopus/proto/octopus"
	"github.com/stretchr/testify/assert"
)

func TestSiteToProto(t *testing.T) {
	tests := []struct {
		name      string
		site      *Site
		protoSite *octopuspb.Site
	}{
		{
			name:      "nil",
			site:      nil,
			protoSite: nil,
		},
		{
			name: "Site",
			site: NewSite("foo"),
			protoSite: &octopuspb.Site{
				Name: "foo",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoSite, test.site.ToProto())
	}
}

func TestPopToProto(t *testing.T) {
	tests := []struct {
		name     string
		pop      *Pop
		protoPop *octopuspb.Pop
	}{
		{
			name:     "nil",
			pop:      nil,
			protoPop: nil,
		},
		{
			name: "Site",
			pop:  NewPop("foo"),
			protoPop: &octopuspb.Pop{
				Name: "foo",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoPop, test.pop.ToProto())
	}
}

func TestColoToProto(t *testing.T) {
	tests := []struct {
		name      string
		colo      *Colo
		protoColo *octopuspb.Colo
	}{
		{
			name:      "nil",
			colo:      nil,
			protoColo: nil,
		},
		{
			name: "Colo without Sites",
			colo: &Colo{
				Id:        42,
				Name:      "foo",
				Status:    "V",
				Region:    "WEUR",
				Tier:      1,
				Animal:    "",
				IsMCP:     true,
				IsFedramp: false,
				Pop: &Pop{
					Name: "foo-a",
				},
				Sites: make([]*Site, 0),
			},
			protoColo: &octopuspb.Colo{
				Id:     42,
				Name:   "foo",
				Status: "V",
				Region: "WEUR",
				Tier:   1,
				IsMcp:  true,
				Pop:    "foo-a",
				Sites:  make([]string, 0),
			},
		},
		{
			name: "Colo with Sites",
			colo: &Colo{
				Id:        42,
				Name:      "foo",
				Status:    "V",
				Region:    "WEUR",
				Tier:      1,
				Animal:    "",
				IsMCP:     true,
				IsFedramp: false,
				Pop: &Pop{
					Name: "foo-a",
				},
				Sites: []*Site{
					{
						Name: "foo-dc01",
					},
				},
			},
			protoColo: &octopuspb.Colo{
				Id:     42,
				Name:   "foo",
				Status: "V",
				Region: "WEUR",
				Tier:   1,
				IsMcp:  true,
				Pop:    "foo-a",
				Sites: []string{
					"foo-dc01",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoColo, test.colo.ToProto())
	}
}
