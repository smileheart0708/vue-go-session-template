<template>
  <div class="log-board">
    <LogBoardToolbar
      :logs-count="logs.length"
      :auto-scroll="autoScroll"
      @clear="clearLogs"
      @toggle-auto-scroll="toggleAutoScroll"
      @export="exportLogs"
    />

    <div ref="logContainer" class="log-board-content">
      <TransitionGroup name="log-list" tag="div" class="log-list">
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
import { ref, onMounted, onUnmounted, nextTick, useTemplateRef } from 'vue'
import { Info } from 'lucide-vue-next'
import LogBoardToolbar from '@/components/logs/LogBoardToolbar.vue'
import { useLogExport, useLogStream } from '@/composables'
import { formatLogMessage, getLevelClass } from '@/utils/logs'

const logContainer = useTemplateRef<HTMLElement>('logContainer')
const autoScroll = ref(true)

function scrollToBottom() {
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
}

const { logs, connect, disconnect, clearLogs } = useLogStream({
  onLog: () => {
    if (autoScroll.value) {
      nextTick(() => {
        scrollToBottom()
      })
    }
  },
})

function toggleAutoScroll() {
  autoScroll.value = !autoScroll.value
  if (autoScroll.value) {
    nextTick(() => {
      scrollToBottom()
    })
  }
}

const { exportLogs } = useLogExport({ logs, formatMessage: formatLogMessage })

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
  connect()
  logContainer.value?.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  disconnect()
  logContainer.value?.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.log-board {
  display: flex;
  flex-direction: column;
  flex: 1;
  height: 100%;
  min-width: 0;
  min-height: 0;
  background: var(--color-background-elevated);
  border-radius: 12px;
  overflow: hidden;
}

.log-board-content {
  flex: 1;
  overflow: auto;
  overscroll-behavior: contain;
  scrollbar-gutter: stable both-edges;
  padding: 0;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
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

.log-list {
  width: max-content;
  min-width: 100%;
}

/* 日志行样式 */
.log-row {
  display: flex;
  align-items: center;
  padding: 0.5rem 1rem;
  white-space: nowrap;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
  min-width: 100%;
  width: max-content;
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
  color: var(--color-text);
  flex: 0 0 auto;
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
