# Tailwind CSS v4 迁移分析（保持 UI 不变）

## 1. 先回答你的 3 个核心问题

### 1.1 是否必须 100% 从原生 CSS 迁移到 Tailwind？
不必须。

最优、优雅、规范的 Tailwind v4 项目通常不是“零原生 CSS”，而是：
- 布局、间距、排版、状态样式以 utilities 为主。
- 设计令牌、重置、全局元素、滚动条、动画关键帧等保留少量语义 CSS。
- 复杂组件交互（例如多层选择器、`TransitionGroup` 命名动画）允许局部 scoped CSS。

### 1.2 Tailwind 能否实现原生 CSS 的特性？
大部分可以，且 v4 支持任意值和变量语法，覆盖面很高；但“能写”不等于“适合写”。

不建议强行全用 utilities 的场景：
- 复杂关键帧与命名过渡（`@keyframes`、`.xxx-enter-active`）。
- 跨子组件结构选择器（例如 `:deep()`、后代组合器链路）。
- 运行时动态值密集（基于 props 的复杂视觉效果）。
- 图表容器这类“布局壳 + JS 驱动尺寸”的结构样式。

### 1.3 你当前项目离“最优 Tailwind v4”还差什么？
主要是“边界治理”和“清理收口”，不是“再把所有 CSS 全改成类名”。


## 2. 当前项目状态审计（本地代码盘点）

### 2.1 全局层级入口：已正确
你当前 `web/src/assets/main.css` 已是正确的 v4 分层导入思路：
- `theme` 层：Tailwind theme + 设计令牌（tokens）+ `@theme inline` 映射。
- `base` 层：preflight + reset/elements/scrollbar/motion。
- `utilities` 层：Tailwind utilities。

这能避免之前出现的两个核心问题：
- `@import` 顺序导致的规则丢失。
- 非 layer 规则覆盖 utilities 导致错位。

### 2.2 待清理（死文件/文档同步）
- `web/src/assets/base.css`：当前已不再被入口引用。
- `web/src/assets/styles/tokens.css`：当前已不再被入口引用。
- `web/src/assets/styles/TOKENS.md`：仍描述 `tokens.css` 为入口，建议同步更新。

### 2.3 仍含 scoped CSS 的 Vue 文件（13 个）
已盘点出以下文件：
- `web/src/views/LogsView.vue`
- `web/src/views/DashboardView.vue`
- `web/src/views/LoginView.vue`
- `web/src/components/logs/LogBoardToolbar.vue`
- `web/src/components/dashboard/DashboardStats.vue`
- `web/src/components/logs/LogBoard.vue`
- `web/src/components/dashboard/RequestChart.vue`
- `web/src/components/dashboard/ModelDistribution.vue`
- `web/src/components/common/AppDialog.vue`
- `web/src/components/common/StatsCard.vue`
- `web/src/components/common/ThemeToggle.vue`
- `web/src/components/common/ToastMessage.vue`
- `web/src/App.vue`（空 style）


## 3. 迁移优先级（UI 不变前提）

### 3.1 第一批：低风险，建议立即迁移到 utilities
这些几乎是“静态容器布局”，改完风险低、收益高：
- `web/src/views/LogsView.vue`
- `web/src/views/DashboardView.vue`
- `web/src/components/dashboard/DashboardStats.vue`
- `web/src/components/common/ThemeToggle.vue`（仅 `position: relative` + `display`）
- `web/src/App.vue`（删除空 `<style scoped></style>`）

### 3.2 第二批：中风险，建议按组件逐步迁移
这些可以迁，但要先处理组件 API 或动画策略：
- `web/src/views/LoginView.vue`
  - 结构可迁移；建议保留极少量局部样式或直接 utility 化。
- `web/src/components/common/AppDialog.vue`
  - 主体样式可 utility 化；命名过渡类可保留 CSS 或改成 Tailwind 过渡类。
- `web/src/components/dashboard/RequestChart.vue`
- `web/src/components/dashboard/ModelDistribution.vue`
  - 卡片外壳可 utility 化；图表容器尺寸约束和极少媒体规则可保留 CSS。

### 3.3 第三批：建议保留原生 CSS（或只做轻量重构）
这些组件强依赖复杂选择器/动画，强行 utility 化会降低可维护性：
- `web/src/components/logs/LogBoard.vue`
  - `TransitionGroup` 命名动画、日志级别样式、响应式细节多。
- `web/src/components/logs/LogBoardToolbar.vue`
  - 使用了 `:deep()` 针对子组件结构；更推荐先扩展 `BaseButton` API 再迁移。
- `web/src/components/common/StatsCard.vue`
  - 使用 `v-bind('color')` + 多层 hover 联动动画，CSS 语义更清晰。
- `web/src/components/common/ToastMessage.vue`
  - 多状态发光、进度条动画、进入离开动画、移动端适配都偏 CSS 语义化。


## 4. 推荐“最优且优雅”的 Tailwind v4 项目形态

### 4.1 样式分层策略（建议固化为团队约定）
- `theme layer`：只放 tokens 与 `@theme` 映射。
- `base layer`：重置、元素默认样式、全局滚动条、全局动效开关。
- `utilities`：页面与组件主体样式优先使用 utilities。
- `components`（可选）：仅放复用度高且 utility 难表达的语义样式。

### 4.2 组件样式策略
- 原则：先 utilities，后 CSS。
- 出现以下任一情况时允许保留 scoped CSS：
  - 命名过渡、关键帧动画。
  - 多层后代选择器、`:deep()`。
  - 动态视觉逻辑明显（尤其是 `v-bind()` CSS 变量参与动画）。

### 4.3 Tailwind 类名检测策略（避免“构建后丢类”）
避免运行时拼接完整类名字符串，改为静态映射：
- 推荐：`const map = { primary: 'bg-accent text-on-accent', ... }`
- 避免：`` `bg-${color}-500` `` 这类动态拼接。


## 5. 你这个项目下一步的“收口清单”

### 5.1 必做（建议本周完成）
1. 删除未使用文件：`web/src/assets/base.css`、`web/src/assets/styles/tokens.css`。
2. 同步更新 `web/src/assets/styles/TOKENS.md` 的入口说明。
3. 完成第一批低风险组件迁移（第 3.1 节）。

### 5.2 可选优化（提升规范性）
1. 在 PR 模板增加“样式层级检查项”（theme/base/utilities 是否越层）。
2. 统一“是否允许 scoped CSS”决策标准（按第 4.2 节）。
3. 对复杂组件（Toast、LogBoard、StatsCard）明确“保留 CSS 是设计选择，不是迁移未完成”。


## 6. 验收标准（UI 不变）

每次迁移建议执行：
1. `pnpm run lint:css`
2. `pnpm run typecheck`
3. `pnpm run build`
4. 页面对比（Dashboard / Keys / Logs / Settings / Login）：
   - 桌面端 + 移动端
   - 浅色 + 深色
   - 动效开启 + `prefers-reduced-motion`

通过标准：
- 无 `@import` / layer 警告。
- 无明显布局跳变、字号错位、组件尺寸回归。
- 交互状态（hover/focus/disabled/active）与重构前一致。


## 7. 参考资料（官方/标准）

- Tailwind v4 + Vite 安装与入口方式  
  https://tailwindcss.com/docs/installation/using-vite
- Tailwind Theme Variables / `@theme`  
  https://tailwindcss.com/docs/theme
- Tailwind Preflight  
  https://tailwindcss.com/docs/preflight
- Tailwind 自定义样式与指令（`@utility`、`@variant`、`@custom-variant` 等）  
  https://tailwindcss.com/docs/functions-and-directives
- Tailwind 类检测（动态类名注意事项）  
  https://tailwindcss.com/docs/detecting-classes-in-source-files
- MDN `@import`  
  https://developer.mozilla.org/docs/Web/CSS/@import
- MDN `@layer`  
  https://developer.mozilla.org/docs/Web/CSS/@layer
