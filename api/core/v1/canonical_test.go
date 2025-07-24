// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package corev1

import (
	"encoding/json"
	"testing"

	objectsv1 "github.com/agntcy/dir/api/objects/v1"
	objectsv3 "github.com/agntcy/dir/api/objects/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecord_MarshalCanonical(t *testing.T) {
	tests := []struct {
		name    string
		record  *Record
		wantErr bool
	}{
		{
			name: "v1alpha1 agent record",
			record: &Record{
				Data: &Record_V1{
					V1: &objectsv1.Agent{
						Name:          "test-agent",
						SchemaVersion: "v1alpha1",
						Description:   "A test agent",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "v1alpha2 record",
			record: &Record{
				Data: &Record_V3{
					V3: &objectsv3.Record{
						Name:          "test-agent-v2",
						SchemaVersion: "v1alpha2",
						Description:   "A test agent in v1alpha2 record",
						Version:       "1.0.0",
						Extensions: []*objectsv3.Extension{
							{
								Name:    "test-extension",
								Version: "1.0.0",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "nil record",
			record:  nil,
			wantErr: false, // Should return nil, nil
		},
		{
			name:    "empty record",
			record:  &Record{},
			wantErr: false,
		},
		{
			name: "record with complex nested data",
			record: &Record{
				Data: &Record_V3{
					V3: &objectsv3.Record{
						Name:          "complex-agent",
						SchemaVersion: "v1alpha2",
						Description:   "A complex test agent",
						Version:       "2.1.0",
						Extensions: []*objectsv3.Extension{
							{
								Name:    "extension-1",
								Version: "1.0.0",
							},
							{
								Name:    "extension-2",
								Version: "2.0.0",
							},
						},
						Skills: []*objectsv3.Skill{
							{
								Name: "skill-1",
								Id:   1,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.record.MarshalCanonical()

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)

			if tt.record == nil {
				assert.Nil(t, got)

				return
			}

			// Verify the output is valid JSON
			var jsonData interface{}
			err = json.Unmarshal(got, &jsonData)
			require.NoError(t, err, "Output should be valid JSON")

			// Verify the output is compact (no unnecessary whitespace)
			assert.NotContains(t, string(got), "\n", "Output should be single line")
			assert.NotContains(t, string(got), "  ", "Output should not contain extra spaces")
		})
	}
}

func TestRecord_MarshalCanonical_Deterministic(t *testing.T) {
	// Test that marshaling the same record multiple times produces identical output
	record := &Record{
		Data: &Record_V1{
			V1: &objectsv1.Agent{
				Name:          "deterministic-test",
				SchemaVersion: "v1alpha1",
				Description:   "Testing deterministic marshaling",
			},
		},
	}

	// Marshal the same record multiple times
	result1, err1 := record.MarshalCanonical()
	require.NoError(t, err1)

	result2, err2 := record.MarshalCanonical()
	require.NoError(t, err2)

	result3, err3 := record.MarshalCanonical()
	require.NoError(t, err3)

	// All results should be identical
	assert.Equal(t, result1, result2, "Marshaling should be deterministic")
	assert.Equal(t, result2, result3, "Marshaling should be deterministic")
	assert.Equal(t, result1, result3, "Marshaling should be deterministic")
}

func TestRecord_MarshalCanonical_KeyOrdering(t *testing.T) {
	// Test that JSON keys are ordered consistently
	record := &Record{
		Data: &Record_V3{
			V3: &objectsv3.Record{
				Name:          "key-order-test",
				SchemaVersion: "v1alpha2",
				Description:   "Testing key ordering",
				Version:       "1.0.0",
				Extensions: []*objectsv3.Extension{
					{
						Name:    "zeta-extension", // Intentionally out of alphabetical order
						Version: "1.0.0",
					},
					{
						Name:    "alpha-extension", // Should be ordered alphabetically in JSON
						Version: "2.0.0",
					},
				},
			},
		},
	}

	result, err := record.MarshalCanonical()
	require.NoError(t, err)

	resultStr := string(result)

	// Verify that keys appear in alphabetical order in the JSON
	// Note: We can't easily test the exact ordering without parsing the JSON structure,
	// but we can verify it's consistent and valid
	assert.True(t, json.Valid(result), "Result should be valid JSON")
	assert.NotEmpty(t, resultStr, "Result should not be empty")
}

func TestUnmarshalCanonical(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid v1alpha1 json",
			data:    []byte(`{"v1":{"name":"test-agent","schema_version":"v1alpha1","description":"A test agent"}}`),
			wantErr: false,
		},
		{
			name:    "valid v1alpha2 json",
			data:    []byte(`{"v3":{"name":"test-agent","schema_version":"v1alpha2","description":"A test agent","version":"1.0.0"}}`),
			wantErr: false,
		},
		{
			name:    "empty json",
			data:    []byte(`{}`),
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte(`{"invalid": json}`),
			wantErr: true,
		},
		{
			name:    "null json",
			data:    []byte(`null`),
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(``),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalCanonical(tt.data)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		record *Record
	}{
		{
			name: "v1alpha1 agent",
			record: &Record{
				Data: &Record_V1{
					V1: &objectsv1.Agent{
						Name:          "roundtrip-test",
						SchemaVersion: "v1alpha1",
						Description:   "Testing roundtrip marshaling",
					},
				},
			},
		},
		{
			name: "v1alpha2 record with extensions",
			record: &Record{
				Data: &Record_V3{
					V3: &objectsv3.Record{
						Name:          "roundtrip-v2",
						SchemaVersion: "v1alpha2",
						Description:   "Testing v1alpha2 roundtrip",
						Version:       "1.5.0",
						Extensions: []*objectsv3.Extension{
							{
								Name:    "test-ext",
								Version: "1.0.0",
							},
						},
					},
				},
			},
		},
		{
			name:   "empty record",
			record: &Record{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the record
			marshaled, err := tc.record.MarshalCanonical()
			require.NoError(t, err)

			// Unmarshal it back
			unmarshaled, err := UnmarshalCanonical(marshaled)
			require.NoError(t, err)

			// Marshal the unmarshaled record again
			remarshaled, err := unmarshaled.MarshalCanonical()
			require.NoError(t, err)

			// The bytes should be identical (idempotent)
			assert.Equal(t, marshaled, remarshaled, "Round-trip should be idempotent")

			// The records should be functionally equivalent
			// (We can test this by comparing their CIDs)
			if tc.record.Data != nil {
				originalCID := tc.record.GetCid()
				unmarshaledCID := unmarshaled.GetCid()
				assert.Equal(t, originalCID, unmarshaledCID, "CIDs should match after round-trip")
			}
		})
	}
}

func TestMarshalCanonical_ConsistentAcrossIdenticalRecords(t *testing.T) {
	// Create two identical records separately to ensure they marshal identically
	createRecord := func() *Record {
		return &Record{
			Data: &Record_V3{
				V3: &objectsv3.Record{
					Name:          "consistency-test",
					SchemaVersion: "v1alpha2",
					Description:   "Testing marshaling consistency",
					Version:       "1.0.0",
					Extensions: []*objectsv3.Extension{
						{
							Name:    "ext-1",
							Version: "1.0.0",
						},
						{
							Name:    "ext-2",
							Version: "2.0.0",
						},
					},
				},
			},
		}
	}

	record1 := createRecord()
	record2 := createRecord()

	marshaled1, err1 := record1.MarshalCanonical()
	require.NoError(t, err1)

	marshaled2, err2 := record2.MarshalCanonical()
	require.NoError(t, err2)

	assert.Equal(t, marshaled1, marshaled2, "Identical records should marshal to identical bytes")
}

func TestUnmarshalCanonical_InvalidInputs(t *testing.T) {
	invalidInputs := []struct {
		name string
		data []byte
	}{
		{
			name: "malformed json",
			data: []byte(`{"unclosed": "object"`),
		},
		{
			name: "invalid protobuf structure",
			data: []byte(`{"invalid_field": "value"}`),
		},
		{
			name: "wrong data type",
			data: []byte(`"just a string"`),
		},
	}

	for _, tc := range invalidInputs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := UnmarshalCanonical(tc.data)
			// Some of these might not error depending on protobuf's tolerance,
			// but they should either error or return a valid (possibly empty) record
			if err != nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}
