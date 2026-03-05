# Vue-Go Session Template

Go + Vue 的单体全栈模板。后端使用 Gin，前端使用 Vue 3 + TypeScript，生产构建后前端静态资源通过 `go:embed` 打包进单个二进制。

## 1. 架构与特性

- 前端：Vue 3.5+、TypeScript、Pinia、Vite
- 后端：Go 1.26+、Gin、`log/slog`
- 认证：`AUTH_KEY` + `gin-contrib/sessions`（filesystem store）
- 日志：SSE 实时推送 + 历史日志接口
- 构建：前端 `web/dist` 嵌入 Go 可执行文件，支持 `.br/.gz`
- 安全增强：前端 API 响应统一 Zod Schema 运行时校验（含 SSE 日志数据）

## 2. 目录结构

```text
.
├── internal/               # Go 后端核心逻辑
│   ├── handlers/           # API 处理器
│   ├── middleware/         # 认证/日志中间件
│   ├── server/             # 路由与 SPA fallback
│   └── stream/             # SSE 日志广播
├── web/                    # Vue 前端工程
│   ├── src/types/api.ts    # API 类型(z.infer) + Zod Schema 校验中心
│   ├── src/utils/api-client.ts # ky 客户端与校验入口
│   └── src/composables/    # 组合式逻辑（含日志流）
├── main.go                 # 应用入口
└── build.ps1               # Windows 一键构建脚本
```

## 3. 快速开始

### 3.1 环境要求

- Go `1.25+`
- Node `20.19+` 或 `22.12+`
- pnpm `10+`

### 3.2 开发模式

1. 配置后端环境变量：

```powershell
Copy-Item .env.example .env.development.local
```

2. 启动后端（项目根目录）：

```powershell
go run .
```

3. 启动前端（新终端）：

```powershell
cd web
pnpm install
pnpm run dev
```

### 3.3 构建与检查

前端：

```powershell
cd web
pnpm run lint
pnpm run typecheck
pnpm run build
```

后端：

```powershell
go vet ./...
go build -tags=go_json -ldflags="-s -w" -o app.exe .
```

Windows 全量构建：

```powershell
.\build.ps1
```

## 4. 环境变量

### 4.1 后端（根目录）

| 变量名 | 默认值 | 说明 |
| --- | --- | --- |
| `PORT` | `8080` | 服务端口 |
| `DATA_DIR` | `.data` | 数据目录 |
| `LOG_LEVEL` | `info` | 日志等级：`debug/info/warn/error` |
| `AUTH_KEY` | 空 | 管理认证密钥；为空时启动自动生成 |
| `COOKIE_SECURE` | `false` | Session Cookie 的 `Secure` 属性，生产环境必须 `true` |
| `SESSION_NAME` | `session_id` | Session Cookie 名称 |
| `SESSION_AUTH_KEY` | 空 | Session 签名密钥；为空时仅开发环境自动生成并告警 |
| `SESSION_ENC_KEY` | 空 | Session 加密密钥（可选，长度必须是 16/24/32 字节） |

### 4.2 前端（`web/`）

| 变量名 | 默认值 | 说明 |
| --- | --- | --- |
| `VITE_API_BASE_URL` | 空 | API 基地址，未设置时默认 `/api` |
| `VITE_API_MODE` | `real` | `real` 或 `mock` |
| `VITE_MOCK_AUTH` | `false` | 是否启用前端模拟认证 |

## 5. API Schema 规范（强制）

所有 JSON 响应都必须在 `web/src/types/api.ts` 定义 Zod schema，并通过 `z.infer` 导出类型；在 `ky` 调用后执行 schema 解析。

### 5.1 统一入口

- HTTP 校验入口：`web/src/utils/api-client.ts`
- 响应异常类型：`ApiResponseValidationError`
- 类型定义与解析：`web/src/types/api.ts`

### 5.2 代码约束

1. 新增接口时，先在 `web/src/types/api.ts` 增加 `zod schema + z.infer type`。
2. 调用时必须执行 `parseWithSchema()`，禁止仅依赖泛型断言。
3. SSE 数据必须使用同一套 schema 做 `parse`（当前日志流已接入）。
4. 发生 `ApiResponseValidationError` 时，按“服务端响应格式异常”处理并记录错误日志。

示例：

```ts
const response = await api.get('dashboard/stats')
const payload = await response.json<unknown>()
const data = parseWithSchema(payload, dashboardStatsResponseSchema, response.url)
```

## 6. 前后端开发规范

### 6.1 前端

- 使用 `<script setup lang="ts">`
- 禁止 `any`
- 禁止类型断言（项目 ESLint 已强制）
- Pinia 使用 Setup Store
- 业务 API 类型集中维护在 `web/src/types/api.ts`

### 6.2 后端

- 使用 `any`，避免 `interface{}`
- 统一 `log/slog` 结构化日志
- 错误包装使用 `fmt.Errorf("context: %w", err)`
- API 路径统一 `/api/*`

## 7. 安全开放与上线建议

### 7.1 最低生产基线

1. 必须通过 HTTPS 暴露服务（建议反向代理终止 TLS）。
2. 设置强随机 `AUTH_KEY`（建议至少 32 字符）。
3. 设置强随机 `SESSION_AUTH_KEY`（建议 32 字节以上）。
4. 设置 `COOKIE_SECURE=true`。
4. 限制公网暴露面：仅暴露网关端口，不直接暴露内部调试端口。
5. 为 `DATA_DIR` 配置最小权限（仅服务账户可读写）。

### 7.2 认证与会话

- Session Cookie：`HttpOnly` + `SameSite=Lax`（已在后端设置）
- Session 默认有效期：7 天（`internal/server/router.go`）
- 轮换 `SESSION_AUTH_KEY` 会使旧会话失效，需规划维护窗口

### 7.3 前端安全

- `VITE_*` 变量会进入前端产物，不能放密钥
- 禁止在 LocalStorage 存储后端真实密钥
- 所有接口响应必须经过 schema 校验，防御后端异常返回或中间链路污染

### 7.4 依赖与供应链安全

建议在 CI 增加：

```powershell
go vet ./...
go test ./...
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

cd web
pnpm run lint
pnpm run typecheck
pnpm audit --prod
pnpm run build
```

### 7.5 反向代理建议

- 仅信任明确的代理来源 IP（Gin `SetTrustedProxies`）
- 在网关层增加速率限制、IP 白名单（按业务需要）
- 设置标准安全响应头（HSTS、X-Content-Type-Options、CSP 等）

## 8. 开放给他人使用时的 Checklist

1. 复制 `.env.example` 并填写 `AUTH_KEY`、`SESSION_AUTH_KEY`、`COOKIE_SECURE`。
2. 执行前后端静态检查和构建。
3. 确认 `VITE_API_MODE=real`，关闭 mock。
4. 使用 HTTPS 网关对外开放，配置代理和限流。
5. 在预发环境验证登录、会话过期、SSE 重连、日志导出。
6. 发布后监控 `slog` 错误与响应校验异常。

## 9. 参考资料（官方/权威）

- Vue Security: https://vuejs.org/guide/best-practices/security.html
- Vite Env and Mode: https://vite.dev/guide/env-and-mode
- Go Security Best Practices: https://go.dev/doc/security/best-practices
- govulncheck: https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
- Gin Trusted Proxies: https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
- OWASP API Security: https://owasp.org/API-Security/
