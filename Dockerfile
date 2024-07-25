# syntax=docker/dockerfile:1

# Build Container
FROM golang:1.22-alpine AS build

WORKDIR /workspace
COPY . /workspace
RUN go mod download
RUN go build -o fhenix .

# Runtime Container
FROM alpine:latest AS runtime

COPY --from=build /workspace/fhenix /fhenix

ENTRYPOINT ["/fhenix"]
