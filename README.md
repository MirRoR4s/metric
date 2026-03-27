# Metric

A minimal Go demo metrics library that exports Prometheus text-format metrics. This repository contains an example service that demonstrates creating `Counter`/`Gauge`, periodic memory collection, and an HTTP middleware to count requests.

This README provides both English and Chinese sections; English is shown first by default.

## English

### Summary

Minimal Go demo metrics library that exports Prometheus text-format metrics. It demonstrates how to implement `Counter` and `Gauge` primitives, collect memory metrics periodically, and expose metrics via HTTP.

### Key Features

- Core metric types: `Counter`, `Gauge`.
- HTTP request counting middleware via `HttpRequestsTotal()`.
- Periodic memory collector `Memory(ctx)` exporting gauges.
- Lightweight `Registry` serving `/metrics` in Prometheus text format.

### Quick Start (≈ 3 minutes each)

Example A — Run and view metrics

```bash
go mod tidy
go run ./cmd
curl http://localhost:8080/hello
curl http://localhost:8080/metrics
```

Example B — Increment request counter

```bash
go run ./cmd &
for i in 1 2 3 4 5; do curl -s http://localhost:8080/hello >/dev/null; done
curl http://localhost:8080/metrics | grep http_requests_total
```

Example C — Inspect memory gauges

```bash
go run ./cmd &
sleep 1
curl http://localhost:8080/metrics | grep process_virtual_memory
```

### API Quick Reference

- `metric.NewCounter(name, help)` — create a `Counter`.
- `Counter.Inc()` / `Counter.Add(v)` — increment counter.
- `metric.NewGauge(name, help)` — create a `Gauge` (`Set`, `Inc`, `Dec`).
- `metric.HttpRequestsTotal()` — returns `(*Counter, func(http.Handler) http.Handler)`.
- `metric.Memory(ctx)` — starts memory collection; returned object has `Stop()` and `WritePrometheus()`.
- `registry := metric.NewRegistry(); registry.Register(metrics...)` — register metrics and get an HTTP handler via `registry.Handler()`.

### Run & Develop

```bash
go mod tidy
go run ./cmd
```

Extend by adding types that implement `WritePrometheus()` and register them with `Registry`.

---

## 中文

### 简介

该项目旨在提供一个最小可运行的指标采集示例，便于学习 Prometheus 文本格式与自定义指标实现方法。主要用途：教学、原型验证与嵌入式服务的轻量监控。核心实现位于 `pkg`，示例服务在 `cmd/main.go`。

### 主要功能

- 基础指标类型：`Counter`（只增）、`Gauge`（可增减）。
- HTTP 请求计数：`HttpRequestsTotal()` 返回 `*Counter` 与用于包装 `http.Handler` 的中间件。
- 内存采集：`Memory(ctx)` 周期性采集虚拟内存并导出为 `Gauge`。
- 简单的 `Registry`：注册任意实现 `WritePrometheus()` 的对象，并在 `/metrics` 以 Prometheus 文本格式暴露。

### 文件结构（要点）

- `cmd/main.go`：示例服务，监听 `:8080`，注册指标并提供 `/hello` 与 `/metrics`。
- `pkg/metic.go`：定义 `Metric`、`Counter`、`Gauge`、内存采集等。
- `pkg/registry.go`：简单的指标注册与 HTTP Handler。
- `pkg/errors.go`：公共错误定义。

### 快速开始（每个示例约 3 分钟）

示例 A — 启动服务并查看指标

```bash
go mod tidy
go run ./cmd
# 在另一个终端：
curl http://localhost:8080/hello
curl http://localhost:8080/metrics
```

示例 B — 使用中间件统计多次请求

```bash
go run ./cmd &
for i in 1 2 3 4 5; do curl -s http://localhost:8080/hello >/dev/null; done
curl http://localhost:8080/metrics | grep http_requests_total
```

示例 C — 观察内存指标

```bash
go run ./cmd &
sleep 1
curl http://localhost:8080/metrics | grep process_virtual_memory
```

### API 快速参考

- `metric.NewCounter(name, help)` — 创建 `Counter`。
- `Counter.Inc()` / `Counter.Add(v)` — 增加计数。
- `metric.NewGauge(name, help)` — 创建 `Gauge`，可 `Set`、`Inc`、`Dec`。
- `metric.HttpRequestsTotal()` — 返回 `(*Counter, func(http.Handler) http.Handler)`，用于统计请求总数。
- `metric.Memory(ctx)` — 启动内存采集，返回包含 `WritePrometheus()` 的对象，`Stop()` 停止采集。
- `registry := metric.NewRegistry(); registry.Register(metrics...)` — 注册指标并通过 `registry.Handler()` 获取 HTTP Handler。

示例：在代码中注册并使用中间件

```go
counter, middleware := metric.HttpRequestsTotal()
registry := metric.NewRegistry()
registry.Register(counter)
mux := http.NewServeMux()
mux.Handle('/metrics', registry.Handler())
srv := &http.Server{Addr: ':8080', Handler: middleware(mux)}
```

### 运行与开发

1. 获取依赖并运行：

```bash
go mod tidy
go run ./cmd
```

2. 调试/扩展：在 `pkg` 中添加新的类型并实现 `WritePrometheus()` 即可被 `Registry` 收集。

### 注意点

- 指标名称受限于正则表达式 `[a-zA-Z_:][a-zA-Z0-9_:]*`，无效名称会导致 panic。
- `Counter.Add` 对负数会 panic，因为 Counter 只允许递增。
- 内存采集默认间隔 500ms，可在 `pkg/metic.go` 中调整。

### 许可证

该仓库未包含明确许可证。用于实验请遵循个人或组织政策。

---

## English

### Summary

Minimal Go demo metrics library that exports Prometheus text-format metrics. It demonstrates how to implement `Counter` and `Gauge` primitives, collect memory metrics periodically, and expose metrics via HTTP.

### Key Features

- Core metric types: `Counter`, `Gauge`.
- HTTP request counting middleware via `HttpRequestsTotal()`.
- Periodic memory collector `Memory(ctx)` exporting gauges.
- Lightweight `Registry` serving `/metrics` in Prometheus text format.

### Quick Start (≈ 3 minutes each)

Example A — Run and view metrics

```bash
go mod tidy
go run ./cmd
curl http://localhost:8080/hello
curl http://localhost:8080/metrics
```

Example B — Increment request counter

```bash
go run ./cmd &
for i in 1 2 3 4 5; do curl -s http://localhost:8080/hello >/dev/null; done
curl http://localhost:8080/metrics | grep http_requests_total
```

Example C — Inspect memory gauges

```bash
go run ./cmd &
sleep 1
curl http://localhost:8080/metrics | grep process_virtual_memory
```

### API Quick Reference

- `metric.NewCounter(name, help)` — create a `Counter`.
- `Counter.Inc()` / `Counter.Add(v)` — increment counter.
- `metric.NewGauge(name, help)` — create a `Gauge` (`Set`, `Inc`, `Dec`).
- `metric.HttpRequestsTotal()` — returns `(*Counter, func(http.Handler) http.Handler)`.
- `metric.Memory(ctx)` — starts memory collection; returned object has `Stop()` and `WritePrometheus()`.
- `registry := metric.NewRegistry(); registry.Register(metrics...)` — register metrics and get an HTTP handler via `registry.Handler()`.

### Run & Develop

```bash
go mod tidy
go run ./cmd
```

Extend by adding types that implement `WritePrometheus()` and register them with `Registry`.

---

如需我将 README 进一步改为更详尽的 API 文档、添加示例代码文件或生成 GoDoc 风格文档，请告诉我你想要的深度（例如：函数签名+示例、全量注释、或导出至 gh-pages）。

