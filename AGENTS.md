# AGENTS.md

> **项目**: Vue-Go 全栈模板 (Golang Embed + Vue 3 SPA)
> **技术栈**: Go 1.25+ (Gin) + Vue 3.5+ (TypeScript) + Pinia + Vite
> **包管理器**: Bun (前端) + Go Modules (后端)

## 1. 构建与检查命令 (Build & Check Commands)

### 1.1 前端开发 (web/ 目录)
```bash
# 开发服务器 (端口 5173)
bun run dev

# 生产构建 (输出到 web/dist/)
bun run build

# 类型检查
bun run typecheck

# 代码检查 (oxlint + eslint，自动修复)
bun run lint

# 代码格式化 (Prettier)
bun run format
```

### 1.2 后端开发 (根目录)
```bash
# 构建 Go 二进制 (自动 embed 前端 dist)
go build -tags=go_json -ldflags="-s -w" -o app.exe .

# 运行开发服务器
go run .

# 检查代码
go vet ./...

# 整理依赖
go mod tidy
```

### 1.3 完整构建流程 (Windows)
```bash
# 使用批处理脚本自动构建前后端
build.bat
```

---

## 2. 前端代码规范 (Frontend Standards)

### 2.1 语法与类型
- **必须使用**: `<script setup lang="ts">`
- **严格类型**: 严禁使用 `any`，优先使用 TypeScript 接口和类型别名
- **组件绑定**: 使用 `defineModel()` 宏替代 `props` + `emit`
- **导入顺序**: 1. Vue 核心 2. 第三方库 3. 本地组件 4. 类型

### 2.2 组件结构
```vue
<script setup lang="ts">
// 1. 导入 (按顺序排列)
import { ref, computed } from 'vue'
import type { ThemeMode } from '@/types'

// 2. Props 定义
interface Props {
  modelValue: string
  disabled?: boolean
}
const props = withDefaults(defineProps<Props>(), { disabled: false })

// 3. Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: string]
  'submit': []
}>()

// 4. 状态与逻辑
const localValue = ref(props.modelValue)

// 5. 生命周期
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

### 2.3 状态管理 (Pinia)
- **必须使用**: Setup Stores (函数式写法)
- **禁止**: Option Stores
- **持久化**: 使用 `useLocalStorage` (VueUse)

```ts
export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)

  function login() { /* ... */ }

  return { user, isAuthenticated, login }
})
```

### 2.4 目录结构
```
web/src/
├── assets/          # 静态资源
├── components/
│   └── common/      # 全局通用组件 (无业务逻辑)
├── composables/     # 组合式函数 (use 开头)
├── router/          # 路由配置
├── stores/          # Pinia 状态管理
├── types/           # TypeScript 类型定义
├── utils/           # 工具函数 (纯函数)
└── views/           # 页面组件
```

---

## 3. 后端代码规范 (Backend Standards)

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
├── configs/         # 配置加载
├── internal/
│   ├── middleware/  # Gin 中间件
│   └── stream/      # SSE 流处理
├── web/             # 前端源码
└── main.go          # 程序入口
```

### 3.4 嵌入式前端
- 前端构建产物通过 `embed` 编译进二进制
- 支持预压缩: `.zst` > `.br` > `.gz`
- SPA Fallback: 未匹配路由返回 `index.html`

---

## 4. API 规范 (API Standards)

- **路径前缀**: `/api/v1/`
- **请求格式**: JSON
- **响应格式**: JSON
- **命名**: JSON 字段 `snake_case`，Go 结构体 `PascalCase`

```go
type CreateUserRequest struct {
    UserName string `json:"user_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}
```

---

## 5. 开发原则 (Development Principles)

| 原则 | 说明 |
| :--- | :--- |
| **单一职责** | 文件/函数/组件只做一件事，超过 250 行必须拆分 |
| **类型安全** | 前端禁用 `any`，后端必须用 `any` |
| **逻辑分离** | UI 展示，Composables 逻辑，Utils 纯计算 |
| **就近原则** | 私有组件/样式/测试靠近使用位置 |

**注意**: 发现同一文件超过 3 个主要函数或组件臃肿时，立即建议重构或拆分文件。