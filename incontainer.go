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

// Package incontainer provides utilities to detect if the current process is running inside a container.
package incontainer

import (
	"bufio"
	"os"
	"strings"
)

// ContainerType represents the type of container detected
type ContainerType string

const (
	// Docker represents Docker container
	Docker ContainerType = "docker"
	// Kubernetes represents Kubernetes pod
	Kubernetes ContainerType = "kubernetes"
	// Podman represents Podman container
	Podman ContainerType = "podman"
	// LXC represents LXC container
	LXC ContainerType = "lxc"
	// Colima represents Colima container
	Colima ContainerType = "colima"
	// OrbStack represents OrbStack container
	OrbStack ContainerType = "orbstack"
	// RancherDesktop represents Rancher Desktop container
	RancherDesktop ContainerType = "rancher-desktop"
	// Unknown represents unknown container type
	Unknown ContainerType = "unknown"
)

// Result contains the detection result
type Result struct {
	InContainer bool
	Type        ContainerType
	Confidence  float64 // 0.0 to 1.0
}

// Detect checks if the current process is running inside a container
func Detect() Result {
	result := Result{
		InContainer: false,
		Type:        Unknown,
		Confidence:  0.0,
	}

	// Check multiple indicators
	indicators := []func() (bool, ContainerType, float64){
		CheckDockerEnv,
		CheckCgroup,
		CheckKubernetes,
		CheckPodman,
		CheckColima,
		CheckOrbStack,
		CheckRancherDesktop,
	}

	maxConfidence := 0.0
	detectedType := Unknown

	for _, check := range indicators {
		if found, containerType, confidence := check(); found {
			result.InContainer = true
			if confidence > maxConfidence {
				maxConfidence = confidence
				detectedType = containerType
			}
		}
	}

	result.Type = detectedType
	result.Confidence = maxConfidence

	return result
}

// CheckDockerEnv checks for Docker-specific indicators
func CheckDockerEnv() (bool, ContainerType, float64) {
	// Check for .dockerenv file
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true, Docker, 0.9
	}

	// Check for docker in hostname
	if hostname, err := os.Hostname(); err == nil {
		if len(hostname) == 12 && isHexString(hostname) {
			return true, Docker, 0.7
		}
	}

	return false, Unknown, 0.0
}

// CheckCgroup checks /proc/1/cgroup for container indicators
func CheckCgroup() (bool, ContainerType, float64) {
	file, err := os.Open("/proc/1/cgroup")
	if err != nil {
		return false, Unknown, 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "docker") {
			return true, Docker, 0.8
		}
		if strings.Contains(line, "kubepods") {
			return true, Kubernetes, 0.8
		}
		if strings.Contains(line, "lxc") {
			return true, LXC, 0.8
		}
		if strings.Contains(line, "podman") {
			return true, Podman, 0.8
		}
		if strings.Contains(line, "colima") {
			return true, Colima, 0.8
		}
	}

	return false, Unknown, 0.0
}

// CheckKubernetes checks for Kubernetes-specific indicators
func CheckKubernetes() (bool, ContainerType, float64) {
	// Check for Kubernetes service account
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
		return true, Kubernetes, 0.9
	}

	// Check for Kubernetes environment variables
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true, Kubernetes, 0.8
	}

	return false, Unknown, 0.0
}

// CheckPodman checks for Podman-specific indicators
func CheckPodman() (bool, ContainerType, float64) {
	// Check for Podman environment variable
	if os.Getenv("container") == "podman" {
		return true, Podman, 0.9
	}

	return false, Unknown, 0.0
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// IsInContainer is a convenience function that returns true if running in any container
func IsInContainer() bool {
	return Detect().InContainer
}

// GetContainerType returns the detected container type
func GetContainerType() ContainerType {
	return Detect().Type
}

// CheckColima checks for Colima-specific indicators
func CheckColima() (bool, ContainerType, float64) {
	// Check for Colima environment variables
	if os.Getenv("COLIMA") != "" {
		return true, Colima, 0.9
	}

	// Check for Colima socket path
	if _, err := os.Stat("/var/run/colima.sock"); err == nil {
		return true, Colima, 0.8
	}

	// Check for Colima in hostname patterns
	if hostname, err := os.Hostname(); err == nil {
		if strings.Contains(hostname, "colima") {
			return true, Colima, 0.7
		}
	}

	return false, Unknown, 0.0
}

// CheckOrbStack checks for OrbStack-specific indicators
func CheckOrbStack() (bool, ContainerType, float64) {
	// Check for OrbStack environment variables
	if os.Getenv("ORBSTACK") != "" {
		return true, OrbStack, 0.9
	}

	// Check for OrbStack socket path
	if _, err := os.Stat("/var/run/orbstack.sock"); err == nil {
		return true, OrbStack, 0.8
	}

	// Check for OrbStack in hostname patterns
	if hostname, err := os.Hostname(); err == nil {
		if strings.Contains(hostname, "orbstack") {
			return true, OrbStack, 0.7
		}
	}

	// Check for OrbStack-specific mount points
	if _, err := os.Stat("/opt/orbstack"); err == nil {
		return true, OrbStack, 0.6
	}

	return false, Unknown, 0.0
}

// CheckRancherDesktop checks for Rancher Desktop-specific indicators
func CheckRancherDesktop() (bool, ContainerType, float64) {
	// Check for Rancher Desktop environment variables
	if os.Getenv("RANCHER_DESKTOP") != "" {
		return true, RancherDesktop, 0.9
	}

	// Check for Rancher Desktop socket path
	if _, err := os.Stat("/var/run/rancher-desktop.sock"); err == nil {
		return true, RancherDesktop, 0.8
	}

	// Check for Rancher Desktop in hostname patterns
	if hostname, err := os.Hostname(); err == nil {
		if strings.Contains(hostname, "rancher") || strings.Contains(hostname, "rd-") {
			return true, RancherDesktop, 0.7
		}
	}

	// Check for k3s (used by Rancher Desktop) in processes
	if _, err := os.Stat("/usr/local/bin/k3s"); err == nil {
		return true, RancherDesktop, 0.6
	}

	return false, Unknown, 0.0
}
