<template>
  <div
    class="flex w-full flex-col items-center gap-2"
    :class="{ 'opacity-70': disabled }"
  >
    <p
      v-if="showSummary"
      class="m-0 w-full text-center text-[clamp(0.875rem,1vw,1.125rem)] leading-5 tracking-[0.02em] text-text-secondary"
    >
      <slot
        name="summary"
        :start="rangeStart"
        :end="rangeEnd"
        :total="total"
      >
        {{ summaryText }}
      </slot>
    </p>

    <nav
      class="flex w-full flex-wrap items-center justify-center gap-1"
      :aria-label="ariaLabel"
    >
      <button
        class="inline-flex size-8 items-center justify-center rounded-md border border-border bg-bg-surface text-text-primary transition-all duration-200 enabled:hover:border-accent enabled:hover:text-accent enabled:active:translate-y-px disabled:cursor-not-allowed disabled:text-text-tertiary disabled:opacity-70"
        type="button"
        :disabled="!canGoPrevious"
        :aria-label="previousAriaLabel"
        @click="goPrevious"
      >
        <ChevronLeft :size="14" />
      </button>

      <template
        v-for="(item, index) in pageItems"
        :key="`${item.type}-${index}`"
      >
        <button
          v-if="item.type === 'page'"
          :class="getPageButtonClass(item.page)"
          type="button"
          :disabled="disabled"
          :aria-current="item.page === page ? 'page' : undefined"
          @click="goToPage(item.page)"
        >
          {{ item.page }}
        </button>
        <span
          v-else
          class="inline-flex w-3 justify-center text-[0.5625rem] leading-none text-text-tertiary"
          aria-hidden="true"
          >...</span
        >
      </template>

      <button
        class="inline-flex size-8 items-center justify-center rounded-md border border-border bg-bg-surface text-text-primary transition-all duration-200 enabled:hover:border-accent enabled:hover:text-accent enabled:active:translate-y-px disabled:cursor-not-allowed disabled:text-text-tertiary disabled:opacity-70"
        type="button"
        :disabled="!canGoNext"
        :aria-label="nextAriaLabel"
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
  previousAriaLabel?: string
  nextAriaLabel?: string
  summaryFormatter?: (context: { start: number; end: number; total: number }) => string
}

const {
  total,
  pageSize = 50,
  maxVisible = 5,
  disabled = false,
  showSummary = true,
  ariaLabel = '分页导航',
  previousAriaLabel = '上一页',
  nextAriaLabel = '下一页',
  summaryFormatter,
} = defineProps<Props>()

const page = defineModel<number>({ required: true })
const pageButtonBaseClass =
  'inline-flex size-8 items-center justify-center rounded-md border text-xs leading-none font-semibold transition-all duration-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-(--sys-color-focus-ring) disabled:cursor-not-allowed disabled:text-text-tertiary disabled:opacity-70'
const inactivePageButtonClass = `${pageButtonBaseClass} border-border bg-bg-surface text-text-primary enabled:hover:border-accent enabled:hover:text-accent enabled:active:translate-y-px`
const activePageButtonClass = `${pageButtonBaseClass} border-accent bg-accent text-on-accent enabled:hover:border-accent-hover enabled:hover:bg-accent-hover enabled:active:border-accent-active enabled:active:bg-accent-active`

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

const summaryText = computed<string>(() => {
  const context = { start: rangeStart.value, end: rangeEnd.value, total }
  if (summaryFormatter) {
    return summaryFormatter(context)
  }
  return `显示 ${context.start}-${context.end} 条 / 共 ${context.total} 条`
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

function getPageButtonClass(targetPage: number): string {
  if (targetPage === page.value) {
    return activePageButtonClass
  }
  return inactivePageButtonClass
}
</script>
