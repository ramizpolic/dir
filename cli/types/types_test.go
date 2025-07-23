// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

//nolint:testifylint
package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectOASFVersion(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectedVer string
		expectError bool
	}{
		{
			name:        "v1 with schema_version",
			jsonData:    `{"schema_version": "v1", "name": "test-agent"}`,
			expectedVer: "v1",
			expectError: false,
		},
		{
			name:        "v2 with schema_version",
			jsonData:    `{"schema_version": "v2", "name": "test-agent"}`,
			expectedVer: "v2",
			expectError: false,
		},
		{
			name:        "v3 with schema_version",
			jsonData:    `{"schema_version": "v3", "name": "test-record"}`,
			expectedVer: "v3",
			expectError: false,
		},
		{
			name:        "no schema_version defaults to v1",
			jsonData:    `{"name": "test-agent", "version": "1.0"}`,
			expectedVer: "v1",
			expectError: false,
		},
		{
			name:        "invalid json",
			jsonData:    `{"name": "test-agent"`,
			expectedVer: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := DetectOASFVersion([]byte(tt.jsonData))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVer, version)
			}
		})
	}
}

func TestLoadOASFFromReader(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
		expectV1    bool
		expectV2    bool
		expectV3    bool
	}{
		{
			name:        "valid v1 agent",
			jsonData:    `{"schema_version": "v1", "name": "test-agent", "version": "1.0"}`,
			expectError: false,
			expectV1:    true,
		},
		{
			name:        "valid v2 agent record",
			jsonData:    `{"schema_version": "v2", "name": "test-agent", "version": "1.0"}`,
			expectError: false,
			expectV2:    true,
		},
		{
			name:        "valid v3 record",
			jsonData:    `{"schema_version": "v3", "name": "test-record", "version": "1.0"}`,
			expectError: false,
			expectV3:    true,
		},
		{
			name:        "no schema_version defaults to v1",
			jsonData:    `{"name": "test-agent", "version": "1.0"}`,
			expectError: false,
			expectV1:    true,
		},
		{
			name:        "unsupported version",
			jsonData:    `{"schema_version": "v4", "name": "test-record"}`,
			expectError: true,
		},
		{
			name:        "invalid json",
			jsonData:    `{"name": "test-agent"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.jsonData)
			record, err := LoadOASFFromReader(reader)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, record)
			} else {
				require.NoError(t, err)
				require.NotNil(t, record)

				//nolint:gocritic // if-else chain is clearer than switch for boolean flag testing in tests
				if tt.expectV1 {
					assert.NotNil(t, record.GetV1())
					assert.Nil(t, record.GetV2())
					assert.Nil(t, record.GetV3())
				} else if tt.expectV2 {
					assert.Nil(t, record.GetV1())
					assert.NotNil(t, record.GetV2())
					assert.Nil(t, record.GetV3())
				} else if tt.expectV3 {
					assert.Nil(t, record.GetV1())
					assert.Nil(t, record.GetV2())
					assert.NotNil(t, record.GetV3())
				}
			}
		})
	}
}
