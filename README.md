# incontainer

A Go library to detect if the current process is running inside a container.

## Features

- Detects various container types: Docker, Kubernetes, Podman, LXC, Colima, OrbStack, Rancher Desktop
- Provides confidence levels for detection accuracy
- Lightweight with no external dependencies
- Cross-platform support

## Installation

```bash
go get github.com/shibukawa/incontainer
```

## Usage

### Simple Detection

```go
package main

import (
    "fmt"
    "github.com/shibukawa/incontainer"
)

func main() {
    if incontainer.IsInContainer() {
        fmt.Println("Running inside a container!")
    } else {
        fmt.Println("Running on host system")
    }
}
```

### Detailed Detection

```go
package main

import (
    "fmt"
    "github.com/shibukawa/incontainer"
)

func main() {
    result := incontainer.Detect()
    
    fmt.Printf("In Container: %t\n", result.InContainer)
    fmt.Printf("Container Type: %s\n", result.Type)
    fmt.Printf("Confidence: %.2f\n", result.Confidence)
}
```

## API Reference

### Types

#### `ContainerType`
```go
type ContainerType string
```

Supported container types:
- `Docker`: Docker container
- `Kubernetes`: Kubernetes pod  
- `Podman`: Podman container
- `LXC`: LXC container
- `Colima`: Colima container
- `OrbStack`: OrbStack container
- `RancherDesktop`: Rancher Desktop container
- `Unknown`: Unknown or undetected container type

#### `Result`
```go
type Result struct {
    InContainer bool          // Whether running in a container
    Type        ContainerType // Detected container type
    Confidence  float64       // Detection confidence (0.0 to 1.0)
}
```

### Functions

#### `Detect() Result`
Performs comprehensive container detection and returns detailed results.

#### `IsInContainer() bool`
Convenience function that returns `true` if running in any container.

#### `GetContainerType() ContainerType`
Returns the detected container type.

## Detection Methods

The library uses multiple detection methods to identify containers:

| Detection Method    | Description                                                                            |
|---------------------|----------------------------------------------------------------------------------------|
| Docker Environment  | Checks for `.dockerenv` file and Docker-specific hostname patterns                    |
| Control Groups      | Analyzes `/proc/1/cgroup` for container-specific entries                               |
| Kubernetes          | Looks for Kubernetes service account files and environment variables                  |
| Podman              | Checks for Podman-specific environment variables                                      |
| Colima              | Detects Colima environment variables, socket paths, and hostname patterns             |
| OrbStack            | Identifies OrbStack through environment variables, socket paths, and mount points     |
| Rancher Desktop     | Recognizes Rancher Desktop via environment variables, socket paths, and k3s binaries  |

## Command Line Tool

A CLI tool is available for testing and debugging:

```bash
# Build the CLI tool
go build -o incontainer ./cmd/incontainer

# Simple check
./incontainer

# Verbose output with detailed checks
./incontainer -v

# JSON output
./incontainer -json
```

### CLI Exit Codes
- `0`: Running in a container
- `1`: Not running in a container  
- `2`: Error occurred

## Testing

### Unit Tests
```bash
go test ./...
```

Run benchmarks:
```bash
go test -bench=.
```

### Container Testing
Test the library in different container environments:

```bash
# Build and test with Docker
make docker-test

# Test locally
make local-test

# Run all tests
make test-all
```

### GitHub Actions
The repository includes GitHub Actions workflows that automatically test the library in various container environments:
- Docker
- Podman  
- Kubernetes

The workflow builds a container image and runs detection tests in each environment to verify accuracy.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.