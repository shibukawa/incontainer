// Copyright 2025 Yoshiki Shibukawa
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package incontainer

import (
	"testing"
)

func TestDetect(t *testing.T) {
	result := Detect()

	// Basic validation
	if result.Confidence < 0.0 || result.Confidence > 1.0 {
		t.Errorf("Confidence should be between 0.0 and 1.0, got %f", result.Confidence)
	}

	// If in container, type should not be Unknown
	if result.InContainer && result.Type == Unknown {
		t.Error("If InContainer is true, Type should not be Unknown")
	}

	// If not in container, confidence should be 0
	if !result.InContainer && result.Confidence != 0.0 {
		t.Errorf("If not in container, confidence should be 0.0, got %f", result.Confidence)
	}
}

func TestIsInContainer(t *testing.T) {
	result := IsInContainer()
	detectResult := Detect()

	if result != detectResult.InContainer {
		t.Error("IsInContainer() should match Detect().InContainer")
	}
}

func TestGetContainerType(t *testing.T) {
	containerType := GetContainerType()
	detectResult := Detect()

	if containerType != detectResult.Type {
		t.Error("GetContainerType() should match Detect().Type")
	}
}

func TestIsHexString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC123", true},
		{"123456789abc", true},
		{"xyz123", false},
		{"123-456", false},
		{"", true}, // empty string is technically all hex
	}

	for _, test := range tests {
		result := isHexString(test.input)
		if result != test.expected {
			t.Errorf("isHexString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkDetect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Detect()
	}
}

func BenchmarkIsInContainer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsInContainer()
	}
}
