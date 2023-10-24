package model

//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColoAndPop(t *testing.T) {
	topology := NewTopology()

	// Add Pop abc-a
	pop_a := topology.AddPopIfNotExists("abc-a")
	assert.Equal(t, topology.Pops["abc-a"], pop_a)
	assert.Equal(t, topology.Pops["abc-a"].Colos, []*Colo{})

	// Re-add Pop abc-a
	topology.AddPopIfNotExists("abc-a")
	assert.Equal(t, topology.Pops["abc-a"], pop_a)
	assert.Equal(t, topology.Pops["abc-a"].Colos, []*Colo{})

	// Add Colo abc01
	colo_a01 := topology.AddColoIfNotExists(1, "abc01", "abc-a")
	assert.Equal(t, topology.Colos[1], colo_a01)
	assert.Equal(t, topology.Colos[1].Pop, pop_a)
	assert.Equal(t, topology.Pops["abc-a"].Colos, []*Colo{colo_a01})

	// Re-add colo abc01
	colo_a01_2nd := topology.AddColoIfNotExists(1, "abc01", "abc-a")
	assert.Equal(t, colo_a01, colo_a01_2nd)
	assert.Equal(t, topology.Colos[1], colo_a01)
	assert.Equal(t, topology.Colos[1].Pop, pop_a)
	assert.Equal(t, topology.Pops["abc-a"].Colos, []*Colo{colo_a01})
}
