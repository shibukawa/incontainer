# syntax=docker/dockerfile:1

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.25rc2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/iscontainer ./cmd/iscontainer

################################################################################
# Create a new stage for running the application using distroless base image.
# Distroless images contain only your application and its runtime dependencies.
# They do not contain package managers, shells or any other programs you would
# expect to find in a standard Linux distribution.
FROM gcr.io/distroless/static-debian12:latest AS final

# Copy the executable from the "build" stage.
COPY --from=build /bin/iscontainer /iscontainer

# What the container should run when it is started.
ENTRYPOINT [ "/iscontainer" ]
CMD [ "-v" ]
