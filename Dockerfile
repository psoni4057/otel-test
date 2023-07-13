############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder

# Install git
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
ADD otel-test /

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /otel-test


EXPOSE 8080

# Run the hello binary.
WORKDIR /
ENTRYPOINT ["/otel-test"]
