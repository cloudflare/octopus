//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	"testing"

	bnet "github.com/bio-routing/bio-rd/net"
	"github.com/stretchr/testify/assert"
)

func TestNewIP(t *testing.T) {
	bioIP := bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32)

	ip := NewIP(bioIP)
	expectedIP := IP{
		Address:  bioIP,
		MetaData: NewMetaData(),
	}

	assert.Equal(t, expectedIP, ip)
}

func TestAppendIPIfNotExists(t *testing.T) {
	IPs := make([]IP, 0)

	tests := []struct {
		name     string
		IP       IP
		expected []IP
	}{
		{
			name: "Add 1st IP",
			IP:   NewIP(bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32)),
			expected: []IP{
				{
					Address:  bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32),
					MetaData: NewMetaData(),
				},
			},
		},
		{
			name: "Add 1st IP again",
			IP:   NewIP(bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32)),
			expected: []IP{
				{
					Address:  bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32),
					MetaData: NewMetaData(),
				},
			},
		},
		{
			name: "Add 2nd IP",
			IP:   NewIP(bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 2), 32)),
			expected: []IP{
				{
					Address:  bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32),
					MetaData: NewMetaData(),
				},
				{
					Address:  bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 2), 32),
					MetaData: NewMetaData(),
				},
			},
		},
	}

	for _, test := range tests {
		IPs = appendIPIfNotExists(IPs, test.IP)
		assert.Equal(t, test.expected, IPs)
	}
}
