package utils

import (
	"testing"

	"github.com/cloudflare/octopus/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestEnrichment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFail bool
		expected model.VLANTag
	}{
		/*
		 * Errors
		 */
		{
			name:     "invalid unit string",
			input:    "a.b.c",
			wantFail: true,
			expected: model.VLANTag{},
		},
		{
			name:     "invalid outer tag",
			input:    "a.1",
			wantFail: true,
			expected: model.VLANTag{},
		},
		{
			name:     "invalid inner tag",
			input:    "1.a",
			wantFail: true,
			expected: model.VLANTag{},
		},
		{
			name:     "invalid single tag",
			input:    "something",
			wantFail: true,
			expected: model.VLANTag{},
		},
		{
			name:     "valid single tag",
			input:    "42",
			wantFail: false,
			expected: model.NewVLANTag(0, 42),
		},
		{
			name:     "valid double tag",
			input:    "23.42",
			wantFail: false,
			expected: model.NewVLANTag(23, 42),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vt, err := ParseUnitStr(test.input)
			if err != nil {
				if test.wantFail {
					return
				}

				t.Errorf("unexpected failure of test %q: %v", test.name, err)
				return
			}

			if test.wantFail {
				t.Errorf("unexpected success of test %q", test.name)
				return
			}

			assert.Equal(t, test.expected, vt, test.name)
		})
	}
}

func TestSanitizeIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "address w/ prefix-length",
			input:    "192.0.2.42/24",
			expected: "192.0.2.42/24",
		},
		{
			name:     "IPv4 address w/o prefix-length",
			input:    "192.0.2.42",
			expected: "192.0.2.42/32",
		},
		{
			name:     "IPv6 address w/o prefix-length",
			input:    "2001:db8::42",
			expected: "2001:db8::42/128",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, SanitizeIPAddress(test.input), test.name)
		})
	}
}

func TestExtractInterfaceAndUnit(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedName    string
		expectedVLANTag model.VLANTag
		wantFail        bool
	}{
		{
			name:            "dot1q",
			input:           "vlan.900",
			expectedName:    "vlan",
			expectedVLANTag: model.NewVLANTag(0, 900),
		},
		{
			name:            "Q-in-Q",
			input:           "xe-0/0/0.100.200",
			expectedName:    "xe-0/0/0",
			expectedVLANTag: model.NewVLANTag(100, 200),
		},
		{
			name:            "no unit",
			input:           "xe-0/0/0",
			expectedName:    "",
			expectedVLANTag: model.NewVLANTag(0, 0),
			wantFail:        true,
		},
		{
			name:            "invalid qot1q",
			input:           "xe-0/0/0.a",
			expectedName:    "",
			expectedVLANTag: model.NewVLANTag(0, 0),
			wantFail:        true,
		},
		{
			name:            "invalid qot1q",
			input:           "xe-0/0/0.0.a",
			expectedName:    "",
			expectedVLANTag: model.NewVLANTag(0, 0),
			wantFail:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ifName, vt, err := extractInterfaceAndUnit(test.input)
			if err != nil && !test.wantFail {
				t.Errorf("unexpected failure of test %q: %v", test.name, err)
				return
			}

			if err == nil && test.wantFail {
				t.Errorf("unexpected success of test %q", test.name)
				return
			}

			assert.Equal(t, test.expectedName, ifName, test.name)
			assert.Equal(t, test.expectedVLANTag, vt, test.name)
		})
	}
}
