# mini-metrics

A small, lightweight Go metrics library that provides simple concurrency-safe Counters and Gauges, a minimal Registry, and an HTTP `/metrics` handler that emits Prometheus-compatible text exposition.

## Overview

mini-metrics focuses on clarity and minimalism: easy-to-use metric primitives for small services, learning projects, or cases where integrating a full Prometheus client is unnecessary.

The library lives under the module `github.com/MirRoR4s/metric` and uses `gopsutil` for optional process memory gauges.

## Usage Example

This minimal runnable example shows the typical usage: create a registry, register metrics, expose `/metrics`, and use an HTTP middleware counter.

Save as `cmd/main.go` (or integrate into your `main` package):

```go
package main

import (
    "context"
    "log"
    "net/http"

    metric "github.com/MirRoR4s/metric/pkg"
)

func main() {
    ctx := context.Background()

    // Create registry
    registry := metric.NewRegistry()

    // HTTP requests counter + middleware
    requests, mw := metric.HttpRequestsTotal()

    // Optional: start process memory gauges
    memMetrics := metric.Memory(ctx)

    // Register metrics
    registry.Register(requests, memMetrics)

    // HTTP handlers
    mux := http.NewServeMux()
    mux.Handle("/metrics", registry.Handler())
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })

    // Wrap with middleware so each request increments the counter
    handler := mw(mux)

    log.Println("listening :8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}
```

Run it and query the endpoints:

```bash
go run ./cmd
curl http://localhost:8080/hello
curl http://localhost:8080/metrics
```

## Quick Start

Add the module to your project:

```bash
go get github.com/MirRoR4s/metric@latest
```

Import the package where you need metrics:

```go
import metric "github.com/MirRoR4s/metric/pkg"
```

Then follow the usage example to register and expose metrics.

## Simple API snippets

Create and use a Counter:

```go
c := metric.NewCounter("example_requests_total", "Total example requests.")
c.Inc()
c.Add(5)
_ = c.Value()
```

Create and use a Gauge:

```go
g := metric.NewGauge("current_workers", "Number of active workers.")
g.Set(3)
g.Inc()
g.Dec()
_ = g.Value()
```

Register metrics with the registry and serve them:

```go
r := metric.NewRegistry()
r.Register(c, g)
http.Handle("/metrics", r.Handler())
```

## Prometheus exposition details

- The handler sets `Content-Type: text/plain; version=0.0.4`.
- Each metric includes `# HELP` and `# TYPE` header lines followed by the metric sample lines.

This format is compatible with Prometheus scrapers.

## Design notes

- Validation: metric names are validated; creating a metric with an invalid name or an empty help string will panic (see `pkg/metic.go`).
- Concurrency: the implementation uses mutexes for thread-safety; operations are safe for concurrent use.
- Scope: this library is intentionally minimal — it aims at pedagogy and small services, not full feature parity with `prometheus/client_golang`.

## Dependencies

- `github.com/shirou/gopsutil/v4` — used for collecting process memory statistics in the optional memory helper.

## Contributing & License

Contributions are welcome. Please open issues or PRs; include tests for non-trivial changes.

There is no LICENSE file in the repository; consider adding an MIT or Apache-2.0 license for clarity.
