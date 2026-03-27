# mini-metrics

> A small, lightweight Go metrics library and example server that exposes basic metrics in Prometheus text format.

## Overview

`mini-metrics` provides simple, concurrency-safe Counter and Gauge metrics, a small Registry for registering metrics, an HTTP `/metrics` handler that emits Prometheus-compatible text exposition, and example helpers (HTTP request counter and periodic process memory gauges).

This repository is implemented in Go (module: `github.com/MirRoR4s/metric`) and uses `gopsutil` to collect memory stats.

## Features

- Counter: `NewCounter`, `Inc`, `Add`, `Value`
- Gauge: `NewGauge`, `Set`, `Inc`, `Dec`, `Add`, `Value`
- Registry: `NewRegistry`, `Register`, `Handler()` to serve `/metrics`
- Example middleware: `HttpRequestsTotal()` returns a counter and a middleware that increments it per request
- Memory metrics: `Memory(ctx)` returns memory-related gauges updated periodically

## Requirements

- Go 1.24 (see `go.mod`)

## Quick Start

1. Clone the repository and download modules:

```
go mod download
```

2. Run the example server:

```
go run ./cmd
```

The example server listens on `:8080`. Open `http://localhost:8080/hello` and visit `http://localhost:8080/metrics` to see metrics in Prometheus text format.

Example using curl:

```
curl http://localhost:8080/hello
curl http://localhost:8080/metrics
```

## Package Overview

- Package path: `pkg` (imported as `github.com/MirRoR4s/metric/pkg` in the example)
- Key types and functions:
  - `Metric` base type with `WritePrometheus(metricType string)`
  - `Counter` (use `NewCounter(name, help)`)
  - `Gauge` (use `NewGauge(name, help)`)
  - `Registry` (use `NewRegistry()`, `Register(...)`, `Handler()`)
  - `HttpRequestsTotal()` helper returns `(*Counter, middleware func(http.Handler) http.Handler)`
  - `Memory(ctx)` helper returns memory gauges and updates them periodically using `gopsutil`

## Prometheus Exposition

The registry's handler sets `Content-Type: text/plain; version=0.0.4` and emits metrics with `# HELP` and `# TYPE` lines followed by metric samples. This is compatible with Prometheus text exposition format used by `prometheus` server scrapers.

## Build & Test

Run all tests:

```
go test ./...
```

Build the module/executable:

```
go build ./...
```

Tip: enable the race detector during development:

```
go test -race ./...
```

## Notes & Design

- Metric name validation and simple error variables live under `pkg/errors.go`.
- The library prioritizes simplicity and teaching value; it is not a drop-in replacement for `prometheus/client_golang` but is useful for small projects and learning.
- The implementation uses mutexes for thread-safety.

## Dependencies

- `github.com/shirou/gopsutil/v4` — used for collecting process memory statistics.

## Contributing

Contributions welcome. Please open issues or PRs. Suggested workflow:

1. Fork the repo
2. Create a feature branch
3. Run `go test ./...` and ensure checks pass
4. Open a PR with a description and tests where appropriate

## License

No license file present in this repository. Add a LICENSE file to clarify terms for reuse.

---

If you'd like, I can add a LICENSE (MIT recommended), a CONTRIBUTING.md, or expand the README with examples showing how to create and register custom metrics. Would you like any of those added?
