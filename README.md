# FlexGate - Lightweight API Gateway

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

FlexGate is a lightweight API gateway prototype built with Go, Gin, and the standard library `httputil.ReverseProxy`. It provides configuration-driven routing, path rewriting, per-route rate limiting, request IDs, unified error handling, and a Docker Compose development environment.

> This project is currently a single-node prototype for learning and demonstrating gateway fundamentals. Benchmark results are local test data and do not represent production capacity.

### Features

- **Configuration-driven routing**: Registers routes from YAML at startup.
- **Reverse proxy**: Forwards requests to upstream services with Go's `ReverseProxy`.
- **Path rewriting**: Strips a configured prefix, for example `/api/users/1` to `/users/1`.
- **Token Bucket rate limiting**: Limits average QPS per route and allows configured bursts.
- **Request tracing and errors**: Generates or forwards `X-Request-ID`, recovers panics, and returns unified upstream errors.
- **Containerized environment**: Starts the gateway and two demo services with Docker Compose and health checks.
- **Repeatable benchmarks**: Runs standard `wrk` scenarios through `scripts/benchmark.sh`.

### Requirements

- Go 1.25+
- Docker Desktop with Docker Compose v2
- `wrk` for local benchmarks

### Quick Start with Docker Compose

Run these commands from the repository root:

```bash
docker compose -f deployments/docker-compose.yml config
make dev-docker
docker compose -f deployments/docker-compose.yml ps
```

Verify the routes:

```bash
curl -i http://127.0.0.1:8080/api/users/1
curl -i http://127.0.0.1:8080/api/orders
```

View logs and stop the environment:

```bash
make logs
make stop-docker
```

Compose uses the service names `user-service` and `order-service` for upstreams. For direct host execution, use a local configuration with `localhost` upstreams instead.

### Make Commands

```bash
make dev-docker   # Build and start all three containers
make logs         # Follow logs from all services
make stop-docker  # Stop and remove containers and network
```

### Benchmarking

After the services are running:

```bash
./scripts/benchmark.sh
```

The script defaults to 4 threads, 100 connections, and 30 seconds. Override them with:

```bash
THREADS=8 CONNECTIONS=200 DURATION=60s ./scripts/benchmark.sh
```

`wrk`'s `Requests/sec` is total response throughput and includes 2xx, 429, and 5xx responses. Always record success rate, status codes, and P99 latency. A high QPS with rate limiting enabled may mostly represent fast 429 responses. Use a separate no-limit configuration or test branch for proxy-capacity benchmarks.

### Configuration

```yaml
server:
  port: 8080

routes:
  - path_prefix: /api/users
    strip_prefix: /api
    upstream: http://user-service:9001
    ratelimit:
      rate: 10
      capacity: 20
```

Configuration is loaded at startup; changes require a gateway restart. Hot reload, service discovery, health-aware load balancing, circuit breaking, Prometheus metrics, and Redis-based distributed rate limiting are not implemented yet.

### Project Layout

```text
flexgate/                  # Gateway source, configuration, and Dockerfile
services/user-service/     # Demo user service
services/order-service/    # Demo order service
deployments/               # Docker Compose orchestration
scripts/                   # Development and benchmark scripts
docs/                      # Architecture documentation
```

### Roadmap

1. Add configuration validation, graceful shutdown, upstream pooling, and timeouts.
2. Add health checks, load balancing, circuit breaking, and Prometheus metrics.
3. Add hot configuration reload and Redis + Lua distributed rate limiting.
4. Add unit, integration, benchmark, and CI coverage.

---

<a name="中文"></a>
## 中文

FlexGate 是一个基于 Go、Gin 和标准库 `httputil.ReverseProxy` 实现的轻量级 API 网关原型，提供配置驱动路由、路径重写、按路由限流、请求 ID、统一错误处理和 Docker Compose 开发环境。

> 当前定位是用于学习和演示网关核心链路的单机原型。压测数据是本地测试结果，不代表生产容量。

### 核心特性

- **配置驱动路由**：启动时从 YAML 注册路由。
- **反向代理**：使用 Go 标准库 `ReverseProxy` 转发到上游服务。
- **路径重写**：支持剥离配置前缀，例如将 `/api/users/1` 转发为 `/users/1`。
- **Token Bucket 限流**：按路由限制平均 QPS，并允许配置突发容量。
- **请求追踪与错误处理**：生成或透传 `X-Request-ID`，恢复 panic，并统一处理上游错误。
- **容器化环境**：通过 Docker Compose 启动网关、两个示例服务和健康检查。
- **可重复压测**：通过 `scripts/benchmark.sh` 运行标准 `wrk` 场景。

### 环境要求

- Go 1.25+
- Docker Desktop（包含 Docker Compose v2）
- `wrk`（仅本地压测需要）

### Docker Compose 快速开始

在项目根目录执行：

```bash
docker compose -f deployments/docker-compose.yml config
make dev-docker
docker compose -f deployments/docker-compose.yml ps
```

验证网关路由：

```bash
curl -i http://127.0.0.1:8080/api/users/1
curl -i http://127.0.0.1:8080/api/orders
```

查看日志并停止环境：

```bash
make logs
make stop-docker
```

Compose 使用 `user-service` 和 `order-service` 作为上游服务名。若直接在宿主机运行服务，请使用上游地址为 `localhost` 的本地配置。

### Make 命令

```bash
make dev-docker   # 构建并启动三个容器
make logs         # 查看全部服务日志
make stop-docker  # 停止并移除容器和网络
```

### 压测

服务启动后执行：

```bash
./scripts/benchmark.sh
```

脚本默认使用 4 个线程、100 个连接、持续 30 秒，也可以覆盖参数：

```bash
THREADS=8 CONNECTIONS=200 DURATION=60s ./scripts/benchmark.sh
```

`wrk` 的 `Requests/sec` 是总响应吞吐，包含 2xx、429 和 5xx。评估网关时应同时记录成功率、状态码和 P99 延迟。开启限流时，高 QPS 可能主要来自快速返回的 429；测试代理能力时应使用单独的不限流配置或测试分支。

### 配置

```yaml
server:
  port: 8080

routes:
  - path_prefix: /api/users
    strip_prefix: /api
    upstream: http://user-service:9001
    ratelimit:
      rate: 10
      capacity: 20
```

配置在启动时加载，修改后需要重启网关。当前尚未实现配置热更新、服务发现、健康节点负载均衡、熔断、Prometheus 指标和 Redis 分布式限流。

### 项目结构

```text
flexgate/                  # 网关核心代码、配置和 Dockerfile
services/user-service/     # 用户示例服务
services/order-service/    # 订单示例服务
deployments/               # Docker Compose 编排
scripts/                   # 开发和压测脚本
docs/                      # 架构文档
```

### 后续计划

1. 增加配置校验、优雅停机、上游连接池和超时控制。
2. 增加健康检查、负载均衡、熔断和 Prometheus 指标。
3. 增加配置热更新和 Redis + Lua 多实例限流。
4. 补充单元测试、集成测试、基准测试和 CI。
