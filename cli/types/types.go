// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"fmt"
	"io"

	corev1 "github.com/agntcy/dir/api/core/v1"
	objectsv1 "github.com/agntcy/dir/api/objects/v1"
	objectsv2 "github.com/agntcy/dir/api/objects/v2"
	objectsv3 "github.com/agntcy/dir/api/objects/v3"
)

// VersionDetector is used to detect OASF schema version from JSON data
type VersionDetector struct {
	SchemaVersion string `json:"schema_version"`
}

// DetectOASFVersion detects the OASF schema version from JSON data
func DetectOASFVersion(data []byte) (string, error) {
	var detector VersionDetector
	err := json.Unmarshal(data, &detector)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON for version detection: %w", err)
	}

	if detector.SchemaVersion == "" {
		// Default to v1 if no schema_version specified for backward compatibility
		return "v1", nil
	}

	return detector.SchemaVersion, nil
}

// LoadOASFFromReader loads OASF data from reader and returns a Record with proper version detection
func LoadOASFFromReader(reader io.Reader) (*corev1.Record, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	version, err := DetectOASFVersion(data)
	if err != nil {
		return nil, err
	}

	switch version {
	case "v1":
		agent := &objectsv1.Agent{}
		err := json.Unmarshal(data, agent)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal v1 Agent: %w", err)
		}
		return &corev1.Record{Data: &corev1.Record_V1{V1: agent}}, nil

	case "v2":
		agentRecord := &objectsv2.AgentRecord{}
		err := json.Unmarshal(data, agentRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal v2 AgentRecord: %w", err)
		}
		return &corev1.Record{Data: &corev1.Record_V2{V2: agentRecord}}, nil

	case "v3":
		record := &objectsv3.Record{}
		err := json.Unmarshal(data, record)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal v3 Record: %w", err)
		}
		return &corev1.Record{Data: &corev1.Record_V3{V3: record}}, nil

	default:
		return nil, fmt.Errorf("unsupported OASF version: %s", version)
	}
}
