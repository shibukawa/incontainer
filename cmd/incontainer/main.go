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
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/incontainer"
)

type DetailedResult struct {
	InContainer bool                      `json:"in_container"`
	Type        incontainer.ContainerType `json:"type"`
	Confidence  float64                   `json:"confidence"`
	Details     map[string]CheckResult    `json:"details"`
}

type CheckResult struct {
	Found      bool                      `json:"found"`
	Type       incontainer.ContainerType `json:"type"`
	Confidence float64                   `json:"confidence"`
}

func main() {
	var (
		jsonOutput = flag.Bool("json", false, "Output in JSON format")
		verbose    = flag.Bool("v", false, "Verbose output with detailed checks")
		help       = flag.Bool("h", false, "Show help")
	)
	flag.Parse()

	if *help {
		fmt.Println("incontainer - Detect if running inside a container")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  incontainer [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  -json    Output in JSON format")
		fmt.Println("  -v       Verbose output with detailed checks")
		fmt.Println("  -h       Show this help")
		fmt.Println()
		fmt.Println("Exit codes:")
		fmt.Println("  0        Running in a container")
		fmt.Println("  1        Not running in a container")
		fmt.Println("  2        Error occurred")
		return
	}

	result := incontainer.Detect()

	if *verbose || *jsonOutput {
		detailed := getDetailedResult(result)

		if *jsonOutput {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(detailed); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(2)
			}
		} else {
			printVerboseResult(detailed)
		}
	} else {
		// Simple output
		if result.InContainer {
			fmt.Printf("Container detected: %s (confidence: %.2f)\n", result.Type, result.Confidence)
		} else {
			fmt.Println("Not running in a container")
		}
	}

	// Exit with appropriate code
	if result.InContainer {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func getDetailedResult(result incontainer.Result) DetailedResult {
	details := make(map[string]CheckResult)

	// Run individual checks to get detailed results
	checks := map[string]func() (bool, incontainer.ContainerType, float64){
		"docker_env": incontainer.CheckDockerEnv,
		"cgroup":     incontainer.CheckCgroup,
		"kubernetes": incontainer.CheckKubernetes,
		"podman":     incontainer.CheckPodman,
	}

	for name, checkFunc := range checks {
		found, containerType, confidence := checkFunc()
		details[name] = CheckResult{
			Found:      found,
			Type:       containerType,
			Confidence: confidence,
		}
	}

	return DetailedResult{
		InContainer: result.InContainer,
		Type:        result.Type,
		Confidence:  result.Confidence,
		Details:     details,
	}
}

func printVerboseResult(result DetailedResult) {
	fmt.Printf("Container Detection Results\n")
	fmt.Printf("===========================\n")
	fmt.Printf("In Container: %t\n", result.InContainer)
	fmt.Printf("Detected Type: %s\n", result.Type)
	fmt.Printf("Confidence: %.2f\n", result.Confidence)
	fmt.Printf("\nDetailed Checks:\n")

	for name, check := range result.Details {
		status := "❌"
		if check.Found {
			status = "✅"
		}
		fmt.Printf("  %s %s: found=%t, type=%s, confidence=%.2f\n",
			status, name, check.Found, check.Type, check.Confidence)
	}
}
