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

type Site struct {
	Name  string
	Colos []*Colo
}

type Pop struct {
	Name  string
	Colos []*Colo
}

type Colo struct {
	Id        uint16
	Name      string
	Status    string
	Region    string
	Tier      uint8
	Animal    string
	IsMCP     bool
	IsFedramp bool
	Pop       *Pop
	Sites     []*Site
}

func NewSite(name string) *Site {
	return &Site{
		Name:  name,
		Colos: make([]*Colo, 0),
	}
}

func NewPop(name string) *Pop {
	return &Pop{
		Name:  name,
		Colos: make([]*Colo, 0),
	}
}

func NewColo(id uint16, name string, pop *Pop) *Colo {
	return &Colo{
		Id:    id,
		Name:  name,
		Pop:   pop,
		Sites: make([]*Site, 0),
	}
}

func (s *Site) ToProto() *octopuspb.Site {
	if s == nil {
		return nil
	}

	return &octopuspb.Site{
		Name: s.Name,
	}
}

func (p *Pop) ToProto() *octopuspb.Pop {
	if p == nil {
		return nil
	}

	return &octopuspb.Pop{
		Name: p.Name,
	}
}

func (c *Colo) ToProto() *octopuspb.Colo {
	if c == nil {
		return nil
	}

	colo := &octopuspb.Colo{
		Id:        uint32(c.Id),
		Name:      c.Name,
		Status:    c.Status,
		Region:    c.Region,
		Tier:      uint32(c.Tier),
		Animal:    c.Animal,
		IsMcp:     c.IsMCP,
		IsFedramp: c.IsFedramp,
		Pop:       c.Pop.Name,
		Sites:     make([]string, len(c.Sites)),
	}

	for i, site := range c.Sites {
		colo.Sites[i] = site.Name
	}

	return colo
}
