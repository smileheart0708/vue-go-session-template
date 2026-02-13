# AGENTS.md

> **项目**: Vue-Go 全栈模板 (Golang Embed + Vue 3 SPA)
> **技术栈**: Go 1.25+ (Gin) + Vue 3.5+ (TypeScript) + Pinia + Vite
> **包管理器**: pnpm (前端) + Go Modules (后端)

---

## 1. 构建与检查命令

### 1.1 前端开发 (web/ 目录)

```bash
# 开发服务器 (端口 5173)
pnpm run dev

# 生产构建 (输出到 web/dist/)
pnpm run build

# 类型检查
pnpm run typecheck

# 代码检查 (oxlint + eslint，自动修复)
pnpm run lint

# 代码格式化 (Prettier)
pnpm run format

# 运行单个 ESLint 规则检查
pnpm oxlint --no-fix src/components/Button.vue
pnpm eslint --no-cache src/utils/http.ts
```

### 1.2 后端开发 (根目录)

```bash
# 构建 Go 二进制 (自动 embed 前端 dist)
go build -tags=go_json -ldflags="-s -w" -o app.exe .

# 运行开发服务器
go run .

# 代码检查
go vet ./...

# 整理依赖
go mod tidy
```

### 1.3 完整构建 (Windows)

```powershell
# PowerShell 脚本 (推荐)
.\build.ps1
```

---

## 2. 前端代码规范

### 2.1 语法与类型

- **必须使用**: `<script setup lang="ts">`
- **严格类型**: 严禁使用 `any`，必须使用 TypeScript 接口/类型别名
- **组件绑定**: 使用 `defineModel()` 宏替代 `props` + `emit`
- **编辑器配置**: `.editorconfig` 已配置 (2 空格缩进, UTF-8, LF 换行)

### 2.2 导入顺序

```ts
// 1. Vue 核心
import { ref, computed, onMounted } from 'vue'

// 2. 第三方库
import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'

// 3. 本地组件/模块
import { useDashboardStore } from '@/stores/dashboard'
import { http } from '@/utils/http'

// 4. 类型定义
import type { User } from '@/types'
```

### 2.3 组件结构

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ThemeMode } from '@/types'

interface Props {
  modelValue: string
  disabled?: boolean
}
const props = withDefaults(defineProps<Props>(), { disabled: false })

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'submit': []
}>()

const localValue = ref(props.modelValue)

onMounted(() => {})
</script>

<template>
  <div class="container">
    <!-- 模板内容 -->
  </div>
</template>

<style scoped>
/* 样式 */
</style>
```

### 2.4 状态管理 (Pinia)

- **必须使用**: Setup Stores (函数式写法)
- **禁止**: Option Stores
- **持久化**: 使用 `@vueuse/core` 的 `useLocalStorage`

```ts
export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)

  function login() { /* ... */ }

  return { user, isAuthenticated, login }
})
```

### 2.5 目录结构

```
web/src/
├── assets/styles/       # CSS 令牌 tokens.*.css
├── components/
│   ├── common/          # 全局通用组件 (无业务逻辑)
│   ├── dashboard/       # 仪表盘相关组件
│   ├── layout/          # 布局组件
│   ├── logs/            # 日志组件
│   └── icons/           # SVG 图标组件
├── composables/         # 组合式函数 (use* 开头)
├── mocks/               # MSW 模拟服务
├── router/              # Vue Router 配置
├── stores/              # Pinia 状态管理
├── types/               # TypeScript 类型定义
├── utils/               # 工具函数 (纯函数)
└── views/               # 页面组件
    └── settings/        # 设置子页面
```

---

## 3. 后端代码规范

### 3.1 现代 Go 语法

- **必须使用**: `any` 替代 `interface{}`
- **禁止**: `interface{}` (除非兼容旧库)
- **日志**: 使用 `log/slog` 结构化日志
- **错误处理**: 使用 `fmt.Errorf("context: %w", err)` 包装错误

### 3.2 代码风格

```go
// ✅ 正确
func ProcessData(v any) error {
    slog.Info("processing", "id", v)
    return nil
}

// ❌ 错误
func ProcessData(v interface{}) error {
    fmt.Println(v)
    return nil
}
```

### 3.3 项目布局

```
/
├── configs/             # 配置加载 (config.go)
├── internal/
│   ├── handlers/       # Gin 处理器 (auth, logs, system)
│   ├── middleware/     # 中间件 (auth, logger)
│   ├── session/        # 会话管理
│   └── stream/         # SSE 流处理
├── web/dist/           # 前端构建产物 (embed)
├── main.go             # 程序入口
└── build.ps1           # 构建脚本
```

### 3.4 嵌入式前端

- 前端构建产物通过 `//go:embed web/dist` 编译进二进制
- 支持预压缩: `.br` (Brotli) > `.gz` (Gzip)
- SPA Fallback: 未匹配路由返回 `index.html`

---

## 4. API 规范

- **路径前缀**: `/api/`
- **请求/响应格式**: JSON
- **命名规则**: JSON 字段 `snake_case`，Go 结构体 `PascalCase`

```go
type CreateUserRequest struct {
    UserName string `json:"user_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}
```

---

## 5. 开发原则

| 原则 | 说明 |
| :--- | :--- |
| **单一职责** | 文件/函数/组件只做一件事，超过 250 行必须拆分 |
| **类型安全** | 前端禁用 `any`，后端必须用 `any` |
| **逻辑分离** | UI 展示，Composables 逻辑，Utils 纯计算 |
| **就近原则** | 私有组件/样式靠近使用位置 |

**注意**: 发现同一文件超过 3 个主要函数或组件臃肿时，应立即重构或拆分文件。

---

## 6. 注意事项

- 当前项目**无测试文件**，如需添加测试请使用 `go test` 和 `vitest`
- 前端 Mock 服务使用 MSW (`web/src/mocks/`)
- 开发环境可通过 `.env.development.local` 配置环境变量
