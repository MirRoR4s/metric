# mini-metrics — 项目需求文档（简易指标收集器）

版本：1.0  
作者：MirRoR4s / 项目经理（由 AI 协助）  
日期：2026-03-25

---

## 1. 概述

- 项目名称：mini-metrics  
- 简介：一个轻量的、并发安全的 Go 指标库，用于在应用程序中记录并导出基础指标（Counter、Gauge、可选 Histogram），并以 Prometheus 文本格式通过 `/metrics` HTTP 端点暴露，方便监控与告警系统抓取。  
- 目标用户：Go 开发者、运维人员、学习监控系统的初学者。  
- 项目目标（高层）：提供简单、易用、可扩展的指标 API，方便在任意 Go 应用中快速埋点并被 Prometheus 等监控系统抓取。

> 关于metrics和Prometheus，可以参看这篇[知乎](https://zhuanlan.zhihu.com/p/512696957)文章。

## 2. 背景与动机

许多开发者希望理解指标基本概念并在小型/学习项目中使用指标，但完整的库（如 `prometheus/client_golang`）对新手和教学场景略显复杂。`mini-metrics` 提供简单的学习与入门工具，同时具备生产可用的基本功能。

## 3. 关键目标（优先级排序）

### 必须实现（MVP）

1. 并发安全的 Counter（`Inc`, `Add`, `Value`）。  
2. 并发安全的 Gauge（`Set`, `Inc`, `Dec`, `Value`）。  
3. Registry：注册指标实例并统一管理。  
4. HTTP `/metrics` Handler：以 Prometheus 文本格式输出已注册指标。  
5. 示例应用（`example/`）展示如何埋点与查看 `/metrics`。  
6. 单元测试覆盖关键逻辑（线程安全、数值准确性）。  
7. CI：GitHub Actions，运行 `go test`、`go vet`、`go fmt`。

### 后续/可选（非 MVP）

1. Histogram（bucket-based）与简单分位估算支持。  
2. 标签（labels）支持（`metric{key="v"}`）。  
3. 原子操作优化以提升性能（使用 `sync/atomic`）。  
4. Benchmarks（性能基准）。  
5. 更完善的 Prometheus 指标模型兼容性（帮助迁移到 `client_golang`）。

## 4. 范围说明

- 包含：**内存**存储指标，进程内暴露 `/metrics`，Go 语言实现，基础文档与样例。  
- 不包含（初期）：持久化存储、跨进程同步、集群聚合、客户端库以外语言绑定、生产级性能优化（可在后续迭代）。

## 5. 利益相关方

- 产品负责人/项目经理：定义优先级、发布节奏。  
- 开发者：实现功能、编写示例与测试。  
- 维护者/社区贡献者：审查 PR、维护 issue。  
- 用户（最终使用者）：采用库并将 `/metrics` 暴露给监控系统。

## 6. 功能需求（详细）

### 6.1 指标类型

- **Counter**
  - 方法：`NewCounter(name, help string) *Counter`
  - API：`Inc()`, `Add(float64)`, `Value() float64`
  - 语义：只能增加（`Add` 允许正数），实现应保证负值不会产生或应被拒绝/忽略（在设计文档中明确选择）。
- **Gauge**
  - 方法：`NewGauge(name, help string) *Gauge`
  - API：`Set(float64)`, `Inc()`, `Dec()`, `Add(float64)`, `Value() float64`
  - 语义：可增可减，支持设置任意浮点值。
- **Histogram（可选）**
  - bucket-based 实现（用户可传 `buckets []float64`）
  - API：`Observe(float64)`，基于 buckets 输出 counts 与 sum。

### 6.2 Registry（管理）

- 提供 `NewRegistry()` 和 `Register(metric)` 方法。  
- 支持列出已注册指标（内部使用）。  
- 提供 `http.Handler`（`/metrics`）统一输出所有已注册指标内容。

### 6.3 Prometheus 文本格式输出

- 输出必须包含每个 metric 的 `# HELP` 和 `# TYPE` 行，及 `metric{labels} value` 行。  
- `Content-Type: text/plain; version=0.0.4`。  
- 输出示例：

  ```text
  # HELP example_http_requests_total Total example HTTP requests
  # TYPE example_http_requests_total counter
  example_http_requests_total 123
  ```

### 6.4 并发与线程安全

- 所有指标操作必须对并发调用安全（使用 `sync.Mutex` 或 `sync/atomic`）。  
- 注册与导出期间，Metrics 的读取应尽量减少阻塞（使用 `RWMutex` / 原子读取策略）。

### 6.5 错误处理

- 指标创建/注册应对非法 `name`/`help` 进行简单校验（非空、字符集限制）。  
- `Add`/`Set` 等非法参数（例如对 Counter 的 `Add` 负数）应返回 `error` 或记录并忽略（在设计文档中明确选择）。

### 6.6 可观察性与调试

- 提供日志或 debug 模式（可选），便于查看注册的指标列表。

## 7. 非功能需求

- 可用性：README 提供“3 分钟上手”示例。  
- 可测试性：提供单元测试覆盖核心逻辑，目标覆盖率 >= 70%（MVP）。  
- 可维护性：代码清晰、模块划分，遵循 Go 习惯（`pkg/`、`cmd/`、`example/`）。  
- 性能：在基本场景（少量指标、高请求频率）表现合理；初期目标：每秒处理至少 10k 次 `Inc` 调用（粗略目标，需基准测试验证）。  
- 兼容性：Go 1.20+（或当前 LTS），模块化（`go.mod`）。  
- 安全性：HTTP handler 应避免任意格式注入（输出需做格式化，不执行用户输入）。

## 8. API / 接口设计（示例）

### Go API（library）

- package `metrics`
  - `func NewCounter(name, help string) *Counter`
  - `func NewGauge(name, help string) *Gauge`
  - `type Registry struct { ... }`
  - `func NewRegistry() *Registry`
  - `func (r *Registry) Register(m PrometheusFormatter)`
  - `func (r *Registry) Handler() http.Handler`
- PrometheusFormatter 接口：
  - `WritePrometheus() string`

### HTTP

- `GET /metrics` -> 返回所有已注册指标的 Prometheus 文本格式。  
- 示例路径： `/metrics` （`Content-Type: text/plain; version=0.0.4`）

## 9. 数据模型与命名规范

- 指标命名遵循 Prometheus 建议：snake_case，包含单位后缀（如 `_bytes`、`_seconds`）。  
- `name` 与 `help` 必需提供并记录。  
- labels（如果支持）为 `key="value"` 样式，value 需做适当转义。

## 10. 验收标准（每条均应可测试）

- **Counter 功能验证**：
  - 单元测试：并发多 goroutine 并增加，最终 `Value` 与期望匹配。  
  - 文本输出在 `/metrics` 中能看到对应 `# HELP`/`# TYPE`/metric lines。  
- **Gauge 功能验证**：
  - `Set`/`Inc`/`Dec` 操作正确反映 `Value`。  
- **Registry & Handler 验证**：
  - 注册后 `/metrics` 返回包含该指标。  
- **CI 验证**：
  - 在 PR 中通过 `go test`、`go vet`、`go fmt` 检查。  
- **文档验收**：
  - README 能在 3 步内让用户运行 example 并查看 `/metrics` 输出。

## 11. 测试策略

- 单元测试：覆盖 Counter/Gauge 的并发读写、边界条件、注册逻辑。  
- 集成测试：启动 example http 服务，发起若干请求，检查 `/metrics` 输出包含预期值。  
- Benchmarks：为关键操作（`Inc`/`Set`/`Observe`）编写 benchmark（后续）。  
- CI：自动运行所有测试并报告。

## 12. CI/CD 与仓库惯例

- GitHub Actions：
  - 检查：`go fmt`、`go vet`、`go test -v ./...`。  
  - 触发：`pull_request`、`push` 到 `main`。  
- 代码风格：`go fmt`、`golangci-lint`（可选后续添加）。  
- PR 模板：包含变更描述、如何验证、是否包含测试。  
- Issue 模板：bug/feature/good-first-issue。

## 13. 文档与示例

- README 包含：
  - 项目简介、设计哲学。  
  - 快速开始（`go run example/main.go`）。  
  - 示例代码片段（如何创建 Counter/Gauge，如何注册）。  
  - 贡献指南（`CONTRIBUTING.md` 链接）。  
- `CONTRIBUTING.md`：如何运行测试、如何提交 PR、编码规范。  
- `CODE_OF_CONDUCT.md`：推荐使用 Contributor Covenant���

## 14. 里程碑（建议拆为 Sprint）

- **M1（2 天）**：仓库初始化（`go.mod`, README skeleton, LICENSE, .gitignore）、实现 Counter（含测试）。  
- **M2（1–2 天）**：实现 Gauge（含测试）。  
- **M3（1 天）**：实现 Registry 与 `/metrics` Handler，添加 example 服务。  
- **M4（1 天）**：CI（GitHub Actions）、基线文档完善。  
- **M5（后续，2–4 天）**：实现 Histogram、labels 支持、benchmarks。

## 15. 初始 Issue 列表（可直接贴到 GitHub，带标签 `good first issue`）

1. **Issue 1（good first issue）**：实现 Counter（`NewCounter`、`Inc`、`Add`、`Value`）并添加并发测试（路径：`metrics/counter.go`）。  
2. **Issue 2（good first issue）**：实现 Gauge（`Set`、`Inc`、`Dec`、`Value`）并添加单元测试（路径：`metrics/gauge.go`）。  
3. **Issue 3（good first issue）**：实现 `Registry.Register` 和 `Registry.Handler`，确保 `/metrics` 返回已注册指标（路径：`metrics/registry.go`、`example/main.go`）。  
4. **Issue 4（good first issue）**：添加 `example/main.go`，演示 HTTP endpoint（`/hello`）触发指标并在 `/metrics` 查看（路径：`example/main.go`）。  
5. **Issue 5（good first issue）**：添加 GitHub Actions CI（`go test`、`go vet`、`go fmt`）并把 badge 添加到 README。

## 16. 风险与缓解

- 风险：��发实现出错导致竞态（race）和数据错乱。  
  - 缓解：编写并发测试，CI 中启用 `-race` 检查（`go test -race`）。  
- 风险：Prometheus 格式实现不规范导致抓取失败。  
  - 缓解：对输出格式做单元测试并对照 Prometheus 文档。  
- 风险：性能在高并发场景下不足。  
  - 缓解：后续用 `sync/atomic` 优化并添加 benchmark，明确不在 MVP 强求高吞吐。  
- 风险：设计与 `prometheus/client_golang` 差异导致迁移成本高。  
  - 缓解：记录设计决定与兼容性文档，尽量采用类似语义。

## 17. 里程碑验收示例（可执行）

- **M1 验收**：Counter 单元测试通过且无 race 报告。  
- **M3 验收**：启动 example，调用 `/hello` 多次后 `/metrics` 能正确显示 `example_http_requests_total` 的值。  
- **M4 验收**：CI 在 PR 流程中自动通过测试与格式检查。

## 18. 术语表

- Metric / 指标：可测量的数值（Counter/Gauge/Histogram）。  
- Registry：管理已注册指标并导出之组件。  
- Prometheus 文本格式：监控系统 Prometheus 支持的文本导出格式。

## 19. 交付物（MVP）

- 完整仓库（`go.mod`、`pkg/metrics`、`example`、README、CONTRIBUTING、LICENSE）。  
- 单元测试与 CI（GitHub Actions）。  
- README 的快速上手示例。  
- 至少 5 个标注为 `good-first-issue` 的 issue。

## 20. 后续建议（发布后）

- 收集社区反馈以决定添加 labels、Histogram、labels 支持等功能。  
- 添加示例集成（与 Docker Compose + Prometheus + Grafana）用于演示端到端监控流程。  
- 在 README 中增加“迁移指南”以阐明与 `prometheus/client_golang` 的差异与兼容性。

---

（文档结束）
