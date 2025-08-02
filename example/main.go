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

package main

import (
	"fmt"

	"github.com/shibukawa/incontainer"
)

func main() {
	// Simple check
	if incontainer.IsInContainer() {
		fmt.Println("Running inside a container!")
	} else {
		fmt.Println("Running on host system")
	}

	// Detailed detection
	result := incontainer.Detect()
	fmt.Printf("In Container: %t\n", result.InContainer)
	fmt.Printf("Container Type: %s\n", result.Type)
	fmt.Printf("Confidence: %.2f\n", result.Confidence)

	// Get just the container type
	containerType := incontainer.GetContainerType()
	if containerType != incontainer.Unknown {
		fmt.Printf("Detected container type: %s\n", containerType)
	}
}
