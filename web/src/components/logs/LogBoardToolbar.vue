<template>
  <div class="log-board-header">
    <div class="log-board-title">
      <h2>{{ title }}</h2>
    </div>
    <div class="log-board-actions">
      <div ref="exportAnchorRef" class="log-export">
        <BaseButton
          @click="toggleExportMenu"
          :disabled="exportDisabled"
          width="auto"
          :height="36"
          text="导出"
          :icon="Download"
          title="导出日志"
          aria-label="导出日志"
          aria-haspopup="menu"
          :aria-expanded="showExportMenu"
        />
        <DropdownDrawer v-model="showExportMenu" :anchor-el="exportAnchorEl" :min-width="160">
          <button class="dropdown-item" type="button" @click="handleExport('txt')">
            导出 TXT
          </button>
          <button class="dropdown-item" type="button" @click="handleExport('csv')">
            导出 CSV
          </button>
          <button class="dropdown-item" type="button" @click="handleExport('json')">
            导出 JSON
          </button>
        </DropdownDrawer>
      </div>
      <BaseButton
        @click="handleClear"
        :disabled="exportDisabled"
        width="auto"
        :height="36"
        text="清空日志"
        :icon="Trash2"
        title="清空日志"
        aria-label="清空日志"
      />
      <BaseButton
        @click="handleToggleAutoScroll"
        :primary="autoScroll"
        width="auto"
        :height="36"
        :text="autoScrollText"
        :icon="autoScroll ? ArrowDownToLine : ArrowDown"
        :title="autoScrollText"
        :aria-label="autoScrollText"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Trash2, ArrowDownToLine, ArrowDown, Download } from 'lucide-vue-next'
import BaseButton from '@/components/common/BaseButton.vue'
import DropdownDrawer from '@/components/common/DropdownDrawer.vue'
import type { LogExportType } from '@/composables'

interface Props {
  logsCount: number
  autoScroll: boolean
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '实时日志',
})

const emit = defineEmits<{
  clear: []
  'toggle-auto-scroll': []
  export: [type: LogExportType]
}>()

const showExportMenu = ref(false)
const exportAnchorRef = ref<HTMLElement | null>(null)
const exportAnchorEl = computed<HTMLElement | null>(() => exportAnchorRef.value)

const exportDisabled = computed(() => props.logsCount === 0)
const autoScrollText = computed(() => (props.autoScroll ? '自动滚动: 开' : '自动滚动: 关'))

function toggleExportMenu() {
  if (exportDisabled.value) return
  showExportMenu.value = !showExportMenu.value
}

function handleExport(type: LogExportType) {
  showExportMenu.value = false
  emit('export', type)
}

function handleClear() {
  emit('clear')
}

function handleToggleAutoScroll() {
  emit('toggle-auto-scroll')
}
</script>

<style scoped>
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


.log-board-actions {
  display: flex;
  gap: 0.75rem;
}

.log-export {
  position: relative;
}

@media (max-width: 768px) {
  .log-board-header {
    flex-direction: row;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem;
  }

  .log-board-title {
    flex: 1;
    min-width: 0;
    white-space: nowrap;
  }

  .log-board-actions {
    width: auto;
    flex-direction: row;
    align-items: center;
    justify-content: flex-end;
    gap: 0.5rem;
  }

  .log-board-actions :deep(.button-text) {
    display: none;
  }

  .log-board-actions :deep(.base-button) {
    width: 36px;
    height: 36px;
    padding: 0;
    min-width: 36px;
  }
}
</style>
