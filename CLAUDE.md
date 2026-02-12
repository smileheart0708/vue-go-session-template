# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## 项目概述

这是一个单体全栈应用模板，使用 **Go 1.25.6** 后端 + **Vue 3.5** 前端，通过 Go embed 机制将前端构建产物编译进单一可执行文件。

**核心特性:**

- 单一可执行文件部署 (前后端集成)
- 预压缩静态资源 (Brotli + Gzip)
- 实时日志流 (SSE)
- 主题系统 (深色/浅色/自动)
- 响应式设计

---

## 技术栈

### 前端

- **Vue 3.5.27** + **TypeScript 5.9.3**
- **Vite 7.3.1** (构建工具)
- **pnpm 10.29.1** (包管理器)
- **Pinia 3.0.4** (状态管理 - Setup Store 写法)
- **Vue Router 4.6.4** (路由)
- **@vueuse/core 14.2.0** (组合式工具库)
- **ECharts 6.0** (数据可视化)
- **lucide-vue-next 0.563.0** (图标)

### 后端

- **Go 1.25.6**
- **Gin 1.11.0** (Web 框架)
- **log/slog** (结构化日志)
- **tint 1.1.2** (彩色日志输出)

---

## 项目结构

```
/
├── web/                          # 前端 Vue 3 项目
│   ├── src/
│   │   ├── assets/              # 静态资源
│   │   ├── components/
│   │   │   ├── common/          # 通用组件 (Button, Toast, ThemeToggle)
│   │   │   ├── dashboard/       # 仪表板私有组件
│   │   │   ├── icons/           # 自定义图标
│   │   │   └── layout/          # 布局组件 (MainLayout, AppSidebar, AppHeader)
│   │   ├── composables/         # 组合式函数 (useToast, useTheme)
│   │   ├── stores/              # Pinia 状态管理
│   │   ├── utils/               # 工具函数 (format, time)
│   │   ├── views/               # 页面级组件
│   │   ├── router/              # Vue Router 配置
│   │   └── main.ts              # 入口文件
│   ├── dist/                    # 构建输出 (被 Go embed)
│   ├── package.json
│   └── vite.config.ts
├── configs/                      # 后端配置
│   └── config.go                # 环境变量加载
├── internal/                     # 后端私有代码
│   ├── middleware/
│   │   └── logger.go            # 日志中间件 (FanoutHandler)
│   └── stream/
│       └── logger.go            # SSE 日志广播系统
├── main.go                       # 后端入口 (SPA 路由 + 预压缩文件服务)
├── go.mod
├── build.bat                     # Windows 构建脚本
├── Dockerfile                    # Docker 多阶段构建
└── .env.example                  # 环境变量模板
```

---

## 常用命令

### 前端开发

```bash
# 进入前端目录
cd web

# 安装依赖
pnpm install

# 类型检查
pnpm run typecheck

# 代码检查和格式化
pnpm run lint        # 运行 Oxlint + ESLint
pnpm run format      # Prettier 格式化

```

### 后端开发

```bash

# 运行验证
go vet ./...
# 运行格式化
go fmt ./...
```

## 核心架构

### 1. 前后端集成机制

**Go Embed 声明** (main.go:21-22):

```go
//go:embed web/dist
var distFS embed.FS
```

**SPA 路由处理:**

- 支持预压缩文件 (Brotli > Gzip)
- 根据 `Accept-Encoding` 请求头智能选择压缩版本
- 自动 ETag 缓存控制
- 带 hash 的资源: `Cache-Control: public, max-age=31536000, immutable`
- index.html: `Cache-Control: no-cache`
- SPA Fallback: 未匹配的路由返回 index.html

### 2. 日志系统架构

**三层日志处理:**

1. **Console Handler (Tint)** - 彩色终端输出
2. **Broadcast Handler** - SSE 实时流
3. **FanoutHandler** - 分发到多个处理器

**LogBroadcaster 特点** (internal/stream/logger.go):

- 维护 SSE 客户端连接池
- 保存最近 100 条日志历史
- 带缓冲通道 (100) 防止阻塞
- 客户端阻塞时跳过消息，避免影响其他客户端

### 3. 前端路由结构

```
/                    → LoginView (登录页)
/dashboard           → DashboardView (仪表板)
/logs                → LogsView (日志页)
/settings            → SettingsView (设置页)
```

**路由特点:**

- 使用 MainLayout 包装内部路由
- 页面切换使用淡入淡出动画 (Transition name="fade")
- 路由 meta 包含页面标题

### 4. 状态管理模式

**Pinia Setup Store 写法** (stores/theme.ts):

```typescript
export const useThemeStore = defineStore("theme", () => {
  const mode = ref<ThemeMode>("auto");
  const isDark = computed(() => {
    /* ... */
  });

  function setMode(newMode: ThemeMode) {
    /* ... */
  }

  return { mode, isDark, setMode, cycleMode, init };
});
```

**特点:**

- 使用 `ref` 和 `computed` 定义状态
- 函数式写法，更接近 Composition API
- 支持 localStorage 持久化 (useColorMode from @vueuse/core)

### 5. 主题系统

**View Transitions API 集成:**

- 使用 `document.startViewTransition()` 实现主题切换动画
- 圆形扩展效果 (从点击位置向外扩展)
- 降级方案: 不支持时使用简单的 CSS 过渡

**主题模式:**

- `light` - 浅色
- `dark` - 深色
- `auto` - 自动 (跟随系统)

**CSS 变量系统:**

```css
:root {
  --sidebar-width: 200px;
  --header-height: 60px;
  --color-primary: #0078d4;
}

:root[data-theme="dark"] {
  /* 深色主题覆盖 */
}
```

---

## 开发规范

### 核心原则

| 原则         | 说明                                                                        |
| :----------- | :-------------------------------------------------------------------------- |
| **单一职责** | 每个文件、函数、组件只做一件事。如果组件超过 250 行，必须拆分。             |
| **类型安全** | 前端严禁使用 `any` (TS)，后端 **必须** 使用 `any` 代替 `interface{}` (Go)。 |
| **逻辑分离** | UI 负责展示，Composables 负责状态逻辑，Utils 负责纯计算。                   |
| **就近原则** | 私有组件/样式/测试代码应尽可能靠近使用它们的地方。                          |

### 前端规范

**组件开发:**

- 必须使用 `<script setup lang="ts">`
- 使用 Vue 3.4+ 的 `defineModel()` 宏
- Props 定义使用 TypeScript 接口 + `withDefaults`
- 样式使用 Scoped CSS

**目录划分:**

- `components/common/` - 全局通用组件 (无业务逻辑)
- `components/layout/` - 布局组件
- `views/xxx/components/` - 页面私有组件
- `composables/` - 组合式函数 (以 use 开头)
- `stores/` - Pinia Stores (Setup Store 写法)
- `utils/` - 纯 TS/JS 函数 (不依赖 Vue 实例)

**状态管理:**

- 必须使用 **Setup Stores** (函数式写法)
- 禁止使用 Option Stores
- 持久化使用 `useLocalStorage` (VueUse)

### 后端规范

**现代 Go 语法:**

- **必须使用 `any`** 代替 `interface{}`
- 理由: Go 1.25+ 标准推荐，更简洁

**错误处理:**

- 必须使用 `log/slog` 进行结构化日志，使用中文日志输出
- 禁止吞掉错误，使用 `fmt.Errorf("context: %w", err)` 包装

**示例:**

```go
if err != nil {
    slog.Error("failed to create user", "username", u.Name, "error", err)
    c.JSON(500, gin.H{"error": "Internal Server Error"})
    return
}
```

### API 规范

- **路径**: 所有 API 必须以 `/api` 开头
- **格式**: 请求和响应必须是 JSON
- **命名**: JSON 字段使用 `snake_case`，Go 结构体使用 `PascalCase`

**DTO 示例:**

```go
type CreateUserRequest struct {
    UserName string `json:"user_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}
```

### 时间戳规范

- **统一格式**: 优先使用 **Unix 秒时间戳** (int64)

### Vite 配置

**构建优化:**

- Brotli 和 Gzip 双压缩
- 删除原始文件 (deleteOriginalAssets: true)
- 仅保留压缩版本，减小二进制体积

---

## 组件库

### 通用组件 (components/common/)

**BaseButton.vue:**

- 支持主题色模式 (primary prop)
- 支持图标 (lucide-vue-next)
- 灵活的尺寸控制

**IconButton.vue:**

- 三种尺寸: small (32px), medium (40px), large (48px)
- 支持 active 状态

**Toast.vue:**

- 最多同时显示 5 个
- 支持 4 种类型: success, error, warning, info
- 自动过期 + 手动关闭
- 进度条动画

---

## 注意事项

1. **代码拆分**: 组件超过 250 行必须拆分，避免巨型文件
2. **类型安全**: 前端禁用 `any`，后端使用 `any` 代替 `interface{}`
3. **日志规范**: 使用 `slog.Info/Error/Warn/Debug`，不要使用 `fmt.Println`
