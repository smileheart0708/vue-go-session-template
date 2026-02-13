<template>
  <div class="app-pagination" :class="{ 'is-disabled': disabled }">
    <p v-if="showSummary" class="app-pagination__summary">
      <slot name="summary" :start="rangeStart" :end="rangeEnd" :total="total">
        显示 {{ rangeStart }}-{{ rangeEnd }} / 共 {{ total }} 条
      </slot>
    </p>

    <nav class="app-pagination__nav" :aria-label="ariaLabel">
      <button
        class="app-pagination__button app-pagination__button--icon"
        type="button"
        :disabled="!canGoPrevious"
        aria-label="上一页"
        @click="goPrevious"
      >
        <ChevronLeft :size="14" />
      </button>

      <template v-for="(item, index) in pageItems" :key="`${item.type}-${index}`">
        <button
          v-if="item.type === 'page'"
          class="app-pagination__button app-pagination__button--page"
          :class="{ 'is-active': item.page === page }"
          type="button"
          :disabled="disabled"
          :aria-current="item.page === page ? 'page' : undefined"
          @click="goToPage(item.page)"
        >
          {{ item.page }}
        </button>
        <span v-else class="app-pagination__ellipsis" aria-hidden="true">...</span>
      </template>

      <button
        class="app-pagination__button app-pagination__button--icon"
        type="button"
        :disabled="!canGoNext"
        aria-label="下一页"
        @click="goNext"
      >
        <ChevronRight :size="14" />
      </button>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

defineOptions({ name: 'AppPagination' })

type PaginationToken = { type: 'page'; page: number } | { type: 'ellipsis'; key: 'left' | 'right' }

interface Props {
  total: number
  pageSize?: number
  maxVisible?: number
  disabled?: boolean
  showSummary?: boolean
  ariaLabel?: string
}

const {
  total,
  pageSize = 50,
  maxVisible = 5,
  disabled = false,
  showSummary = true,
  ariaLabel = 'Pagination',
} = defineProps<Props>()

const page = defineModel<number>({ default: 1 })

const normalizedPageSize = computed<number>(() => {
  const value = Math.trunc(pageSize)
  return value > 0 ? value : 1
})

const totalPages = computed<number>(() => {
  if (total <= 0) return 1
  return Math.ceil(total / normalizedPageSize.value)
})

const rangeStart = computed<number>(() => {
  if (total <= 0) return 0
  return (page.value - 1) * normalizedPageSize.value + 1
})

const rangeEnd = computed<number>(() => {
  if (total <= 0) return 0
  return Math.min(total, page.value * normalizedPageSize.value)
})

const canGoPrevious = computed<boolean>(() => !disabled && page.value > 1)
const canGoNext = computed<boolean>(() => !disabled && page.value < totalPages.value)

const pageItems = computed<ReadonlyArray<PaginationToken>>(() => {
  const items: PaginationToken[] = []
  const pages = totalPages.value
  const current = page.value
  const normalizedVisibleBase = Math.max(5, Math.trunc(maxVisible))
  const normalizedVisible =
    normalizedVisibleBase % 2 === 0 ? normalizedVisibleBase + 1 : normalizedVisibleBase

  if (pages <= normalizedVisible) {
    for (let value = 1; value <= pages; value += 1) {
      items.push({ type: 'page', page: value })
    }
    return items
  }

  const middleSlots = normalizedVisible - 2
  let start = current - Math.floor(middleSlots / 2)
  let end = current + Math.floor(middleSlots / 2)

  if (start <= 2) {
    start = 2
    end = start + middleSlots - 1
  }
  if (end >= pages - 1) {
    end = pages - 1
    start = end - middleSlots + 1
  }

  items.push({ type: 'page', page: 1 })
  if (start > 2) items.push({ type: 'ellipsis', key: 'left' })

  for (let value = start; value <= end; value += 1) {
    items.push({ type: 'page', page: value })
  }

  if (end < pages - 1) items.push({ type: 'ellipsis', key: 'right' })
  items.push({ type: 'page', page: pages })
  return items
})

watch(
  [() => page.value, totalPages],
  ([nextPage]) => {
    const normalized = clampPage(nextPage)
    if (nextPage !== normalized) {
      page.value = normalized
    }
  },
  { immediate: true },
)

function clampPage(value: number): number {
  if (!Number.isFinite(value)) return 1
  const normalized = Math.trunc(value)
  if (normalized < 1) return 1
  if (normalized > totalPages.value) return totalPages.value
  return normalized
}

function goToPage(nextPage: number): void {
  if (disabled) return
  const normalized = clampPage(nextPage)
  if (normalized !== page.value) {
    page.value = normalized
  }
}

function goPrevious(): void {
  goToPage(page.value - 1)
}

function goNext(): void {
  goToPage(page.value + 1)
}
</script>

<style scoped>
.app-pagination {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.625rem;
}

.app-pagination__summary {
  width: 100%;
  margin: 0;
  text-align: center;
  color: var(--color-text-secondary);
  font-size: clamp(0.875rem, 1vw, 1.125rem);
  line-height: 1.25;
  letter-spacing: 0.02em;
}

.app-pagination__nav {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.app-pagination__button {
  width: 2rem;
  height: 2rem;
  border-radius: 6px;
  border: 1px solid var(--color-border);
  background: var(--color-background-elevated);
  color: var(--color-text);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1;
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    color 0.2s ease,
    background-color 0.2s ease,
    transform 0.2s ease;
}

.app-pagination__button:not(:disabled):hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.app-pagination__button:not(:disabled):active {
  transform: translateY(1px);
}

.app-pagination__button.is-active,
.app-pagination__button.is-active:hover {
  border-color: var(--color-primary);
  background: var(--color-primary);
  color: var(--color-on-primary);
}

.app-pagination__button:disabled {
  cursor: not-allowed;
  color: var(--color-text-tertiary);
  border-color: var(--color-border);
  opacity: 0.7;
}

.app-pagination__button--icon {
  font-size: 0.4375rem;
}

.app-pagination__ellipsis {
  width: 0.75rem;
  text-align: center;
  color: var(--color-text-tertiary);
  font-size: 0.5625rem;
  line-height: 1;
}

.app-pagination.is-disabled .app-pagination__summary {
  opacity: 0.7;
}
</style>
