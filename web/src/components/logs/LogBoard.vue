<template>
  <div class="log-board">
    <div class="log-board-header">
      <div class="log-board-title">
        <h2>实时日志</h2>
        <span v-if="logs.length > 0" class="log-count">{{ logs.length }} 条</span>
      </div>
      <div class="log-board-actions">
        <BaseButton
          @click="clearLogs"
          :disabled="logs.length === 0"
          :width="100"
          :height="36"
          text="清空日志"
        />
        <BaseButton
          @click="toggleAutoScroll"
          :primary="autoScroll"
          :width="120"
          :height="36"
          :text="autoScroll ? '自动滚动: 开' : '自动滚动: 关'"
        />
      </div>
    </div>

    <div ref="logContainer" class="log-board-content">
      <TransitionGroup name="log-list">
        <div
          v-for="(log, index) in logs"
          :key="`${log.time}-${index}`"
          :class="['log-row', `log-row-${getLevelClass(log.level)}`]"
        >
          <span class="log-time">{{ log.time }}</span>
          <span class="log-message">{{ formatLogMessage(log) }}</span>
        </div>
      </TransitionGroup>
      <div v-if="logs.length === 0" class="log-board-empty">
        <Info :size="48" />
        <p>暂无日志数据</p>
        <p class="log-board-empty-hint">等待服务器推送日志...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Info } from 'lucide-vue-next'
import BaseButton from '@/components/common/BaseButton.vue'
import { useAuthStore } from '@/stores/auth'

interface LogEntry {
  time: string
  level: string
  msg: string
  attrs?: Record<string, string | number | boolean>
}

const router = useRouter()
const authStore = useAuthStore()
const logs = ref<LogEntry[]>([])
const logContainer = ref<HTMLElement | null>(null)
const autoScroll = ref(true)
let eventSource: EventSource | null = null

function getLevelClass(level: string): string {
  const upperLevel = level.toUpperCase()
  if (upperLevel === 'ERROR') return 'error'
  if (upperLevel === 'WARN') return 'warn'
  if (upperLevel === 'DEBUG') return 'debug'
  return 'info'
}

function formatLogMessage(log: LogEntry): string {
  let message = log.msg

  // 如果有额外属性，追加到消息后面
  if (log.attrs && Object.keys(log.attrs).length > 0) {
    const attrStrings = Object.entries(log.attrs).map(([key, value]) => {
      // 处理不同类型的值
      if (typeof value === 'string') {
        return `${key}=${value}`
      }
      return `${key}=${JSON.stringify(value)}`
    })
    message += ' ' + attrStrings.join(' ')
  }

  return message
}

function connectSSE() {
  eventSource = new EventSource('/api/logs/stream')

  eventSource.onmessage = (event) => {
    try {
      const logEntry: LogEntry = JSON.parse(event.data)
      logs.value.push(logEntry)

      // 限制日志数量，避免内存溢出
      if (logs.value.length > 500) {
        logs.value.shift()
      }

      // 自动滚动到底部
      if (autoScroll.value) {
        nextTick(() => {
          scrollToBottom()
        })
      }
    } catch (error) {
      console.error('解析日志数据失败:', error)
    }
  }

  eventSource.onerror = (error) => {
    console.error('SSE 连接错误:', error)

    // 检查是否是认证错误
    if (eventSource?.readyState === EventSource.CLOSED) {
      // 尝试验证 session
      authStore.validateSession().then((isValid) => {
        if (!isValid) {
          // session 无效，重定向到登录页
          eventSource?.close()
          router.replace('/login')
          return
        }

        // session 有效但连接断开，5秒后重连
        setTimeout(() => {
          connectSSE()
        }, 5000)
      })
    }
  }
}

function scrollToBottom() {
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
}

function clearLogs() {
  logs.value = []
}

function toggleAutoScroll() {
  autoScroll.value = !autoScroll.value
  if (autoScroll.value) {
    nextTick(() => {
      scrollToBottom()
    })
  }
}

// 监听用户手动滚动
function handleScroll() {
  if (!logContainer.value) return

  const { scrollTop, scrollHeight, clientHeight } = logContainer.value
  const isAtBottom = scrollHeight - scrollTop - clientHeight < 50

  // 如果用户滚动到底部，自动开启自动滚动
  if (isAtBottom && !autoScroll.value) {
    autoScroll.value = true
  }
  // 如果用户向上滚动，关闭自动滚动
  else if (!isAtBottom && autoScroll.value) {
    autoScroll.value = false
  }
}

onMounted(() => {
  connectSSE()
  logContainer.value?.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  eventSource?.close()
  logContainer.value?.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.log-board {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--color-background-elevated);
  border-radius: 12px;
  overflow: hidden;
}

.log-board-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  background: var(--color-background-elevated);
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.log-board-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.log-board-title h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.log-count {
  padding: 0.25rem 0.625rem;
  background: var(--color-primary);
  color: white;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.log-board-actions {
  display: flex;
  gap: 0.75rem;
}

.log-board-content {
  flex: 1;
  overflow: auto;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.log-board-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-tertiary);
  gap: 1rem;
}

.log-board-empty p {
  margin: 0;
  font-size: 1rem;
}

.log-board-empty-hint {
  font-size: 0.875rem;
  color: var(--color-text-tertiary);
}

/* 日志行样式 */
.log-row {
  display: flex;
  align-items: center;
  padding: 0.5rem 1rem;
  white-space: nowrap;
  overflow: hidden;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
}

.log-row:hover {
  background: var(--color-background-secondary);
}

/* 日志等级背景色 */
.log-row-error {
  background: var(--log-row-error-bg);
}

.log-row-warn {
  background: var(--log-row-warn-bg);
}

.log-row-debug {
  background: var(--log-row-debug-bg);
}

.log-row-info {
  background: var(--log-row-info-bg);
}

/* 时间列 */
.log-time {
  flex-shrink: 0;
  width: 180px;
  color: var(--color-text-tertiary);
  font-size: 0.8125rem;
}

/* 消息列 */
.log-message {
  flex: 1;
  min-width: 0;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 日志列表动画 */
.log-list-enter-active {
  transition: all 0.2s ease;
}

.log-list-enter-from {
  opacity: 0;
  transform: translateY(-5px);
}

.log-list-leave-active {
  transition: all 0.15s ease;
}

.log-list-leave-to {
  opacity: 0;
  transform: translateX(-5px);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .log-board-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
    padding: 1rem;
  }

  .log-board-actions {
    width: 100%;
    flex-direction: column;
  }

  .log-time {
    width: 140px;
    font-size: 0.75rem;
  }

  .log-row {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
  }
}

@media (prefers-reduced-motion: reduce) {
  .log-list-enter-active,
  .log-list-leave-active {
    transition: none;
  }
}
</style>
