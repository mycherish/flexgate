# FlexGate - Lightweight API Gateway

[English](#english) | [中文](#中文)

---

<a name="english"></a>
## English

**FlexGate** is a high-performance, lightweight API gateway built with Go. It provides core features such as reverse proxying, dynamic path rewriting, and microservice orchestration, designed for modern cloud-native architectures.

### ✨ Key Features
* **Reverse Proxy**: Seamlessly forward requests to backend upstream services.
* **Dynamic Path Rewriting**: Flexible path manipulation using `strip_prefix` (e.g., `/api/users` → `/users`).
* **Robust Path Handling**: Built-in path normalization using `path.Clean` to prevent common routing errors.
* **Docker-Ready**: Pre-configured `Dockerfile` and `docker-compose` for one-click deployment.
* **RESTful Compliant**: Optimized for standard REST API communication.

### 🚀 Quick Start

#### 1. Prerequisites
* Go 1.23+
* Docker & Docker Compose

#### 2. Run with Docker
```bash
# Clone the repository
git clone <your-repo-url>
cd flexgate-project/deployments

# Launch all services (Gateway + User Service + Order Service)
docker-compose up --build -d
````

#### 3\. Verify

```bash
# Test User Service via Gateway
curl http://localhost:8080/api/users/1
```

-----


## 中文

**FlexGate** 是一个基于 Go 语言开发的高性能轻量级 API 网关。它提供了反向代理、动态路径重写和微服务编排等核心功能，专为现代云原生架构设计。

### ✨ 核心特性

  * **反向代理**：将请求无缝转发至后端上游服务。
  * **动态路径重写**：支持通过 `strip_prefix` 灵活修改请求路径（例如：`/api/users` → `/users`）。
  * **健壮的路径处理**：内置 `path.Clean` 路径标准化逻辑，有效防止常见的路由匹配错误。
  * **容器化支持**：预配置 `Dockerfile` 和 `docker-compose`，支持一键部署。
  * **RESTful 兼容**：针对标准的 REST API 通信进行了深度优化。

### 🛠 技术栈

  * **Language**: Go 1.23
  * **Framework**: Gin (Web Engine)
  * **Config**: YAML v3
  * **Deployment**: Docker / Docker Compose

### 📂 项目目录结构

```text
├── flexgate/          # 网关核心代码 (Gateway Core)
├── services/          # 微服务示例 (Microservices: User, Order)
├── deployments/       # 部署配置 (Docker Compose)
└── scripts/           # 辅助脚本 (Helper Scripts)
```

### ⚙️ 配置示例 (`gateway.yaml`)

```yaml
routes:
  - path_prefix: /api/users
    strip_prefix: /api     # 剥离前缀，保留资源路径
    upstream: http://user-service:9001
```

### 🤝 贡献

欢迎提交 Pull Request 或 Issue 来完善 FlexGate！