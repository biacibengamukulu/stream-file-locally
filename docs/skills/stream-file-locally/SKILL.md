---
name: stream-file-locally
description: Use when working on the stream-file-locally Go service, especially upload, disk/Cassandra storage, public file streaming, Docker Compose, and API documentation tasks.
---

# Stream File Locally

## Project Defaults

- Repository: `github.com/biacibengamukulu/stream-file-locally`
- Runtime: Go with Fiber.
- Architecture: follow the existing DDD-style layout under `internal/domain` and `internal/interfaces`.
- Docker image: `010309/stream-file-locally:latest`.
- Default Cassandra host: `safer.easipath.com`, but prefer environment variables.
- Disk storage must use a Docker-mounted volume, defaulting to a stable container path such as `/data/stream-file-locally`.

## Implementation Guidance

- Complete the existing scaffold instead of replacing it.
- Keep storage behind a domain-facing abstraction.
- Support `STORAGE_DRIVER=disk|cassandra`.
- Run Cassandra keyspace/table migrations from the application when Cassandra storage is enabled.
- Return upload metadata with a public stream URL.
- Keep public streaming separate from metadata fetches.
- Update `docs/my-spec.md`, OpenAPI docs, and integration docs when API behavior changes.

## Validation

- Run `gofmt` on changed Go files.
- Run `go test ./...`.
- Run a build for `./cmd/api` when possible.
- Do not deploy to `ssh safer` without explicit approval.
