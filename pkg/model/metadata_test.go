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

func TestMetaDataToProto(t *testing.T) {

	tests := []struct {
		name     string
		metadata *MetaData
		protoMD  *octopuspb.MetaData
	}{
		{
			name:     "nil",
			metadata: nil,
			protoMD:  nil,
		},

		{
			name:     "Empty MetaData (constructor)",
			metadata: NewMetaData(),
			protoMD:  nil,
		},

		{
			name: "Empty MetaData (manual)",
			metadata: &MetaData{
				Tags:         []string{},
				SemanticTags: map[string]string{},
			},
			protoMD: nil,
		},

		{
			name: "MetaData with one Tag",
			metadata: &MetaData{
				Tags: []string{
					"foo",
				},
				SemanticTags: map[string]string{},
			},
			protoMD: &octopuspb.MetaData{
				Tags: []string{
					"foo",
				},
				SemanticTags: map[string]string{},
			},
		},

		{
			name: "MetaData with SemanticTag",
			metadata: &MetaData{
				Tags: []string{},
				SemanticTags: map[string]string{
					"foo": "bar",
				},
			},
			protoMD: &octopuspb.MetaData{
				Tags: []string{},
				SemanticTags: map[string]string{
					"foo": "bar",
				},
			},
		},

		{
			name: "MetaData with CustomFieldData",
			metadata: &MetaData{
				Tags:            []string{},
				SemanticTags:    map[string]string{},
				CustomFieldData: "{\"foo\": \"bar\"}",
			},
			protoMD: &octopuspb.MetaData{
				Tags:            []string{},
				SemanticTags:    map[string]string{},
				CustomFieldData: "{\"foo\": \"bar\"}",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoMD, test.metadata.ToProto(), test.name)
	}
}
