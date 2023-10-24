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

type MetaData struct {
	Tags            []string
	SemanticTags    map[string]string
	CustomFieldData string
}

func NewMetaData() *MetaData {
	return &MetaData{
		Tags:         make([]string, 0),
		SemanticTags: make(map[string]string),
	}
}

func (m *MetaData) ToProto() *octopuspb.MetaData {
	if m == nil {
		return nil
	}

	if len(m.SemanticTags) == 0 && len(m.Tags) == 0 && m.CustomFieldData == "" {
		return nil
	}

	return &octopuspb.MetaData{
		Tags:            m.Tags,
		SemanticTags:    m.SemanticTags,
		CustomFieldData: m.CustomFieldData,
	}
}
