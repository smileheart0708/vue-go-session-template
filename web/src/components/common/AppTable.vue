<template>
  <section class="app-table" v-bind="attrs">
    <header v-if="hasHeader" class="app-table__header">
      <div class="app-table__title">
        <slot name="title">
          <h2 v-if="title" class="app-table__title-text">{{ title }}</h2>
        </slot>
      </div>
      <div v-if="$slots.actions" class="app-table__actions">
        <slot name="actions" />
      </div>
    </header>

    <div class="app-table__scroll" role="region" :aria-label="regionLabel || undefined">
      <table class="app-table__table" :class="{ 'app-table__table--compact': isCompact }">
        <thead :class="{ 'is-sticky': stickyHeader }">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              scope="col"
              :class="[getAlignClass(column.align), column.headerClass]"
              :style="getColumnStyle(column)"
            >
              <slot :name="`header-${column.key}`" :column="column">
                {{ column.label }}
              </slot>
            </th>
          </tr>
        </thead>

        <tbody v-if="rows.length > 0">
          <tr v-for="(row, rowIndex) in rows" :key="resolveRowKey(row, rowIndex)" class="app-table__row">
            <td
              v-for="column in columns"
              :key="column.key"
              :class="[getAlignClass(column.align), column.cellClass]"
              :style="getColumnStyle(column)"
            >
              <slot
                :name="`cell-${column.key}`"
                :row="row"
                :column="column"
                :value="getCellValue(row, column.key)"
                :row-index="rowIndex"
              >
                {{ resolveCellText(row, column, rowIndex) }}
              </slot>
            </td>
          </tr>
        </tbody>

        <tbody v-else>
          <tr>
            <td class="app-table__empty" :colspan="columns.length">
              <slot name="empty">{{ emptyText }}</slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script lang="ts">
export type CssSize = string | number
export type ColumnKey<TRow extends object> = Extract<keyof TRow, string>
export type ColumnAlign = 'left' | 'center' | 'right'

export interface AppTableColumn<TRow extends object = Record<string, unknown>> {
  key: ColumnKey<TRow>
  label: string
  align?: ColumnAlign
  width?: CssSize
  minWidth?: CssSize
  headerClass?: string
  cellClass?: string
}

// Backward compatible alias for existing imports.
export type GlassTableColumn<TRow extends object = Record<string, unknown>> = AppTableColumn<TRow>
</script>

<script setup lang="ts" generic="TRow extends object">
import { computed, useAttrs, useSlots } from 'vue'
import { useMediaQuery } from '@vueuse/core'

defineOptions({
  name: 'AppTable',
  inheritAttrs: false,
})

type RowIdentity = string | number

interface CellContext {
  row: TRow
  rowIndex: number
  column: AppTableColumn<TRow>
}

interface Props {
  title?: string
  regionLabel?: string
  columns: ReadonlyArray<AppTableColumn<TRow>>
  rows: ReadonlyArray<TRow>
  emptyText?: string
  emptyCellText?: string
  rowKey?: ColumnKey<TRow> | ((row: TRow, index: number) => RowIdentity)
  stickyHeader?: boolean
  compactMediaQuery?: string
  formatCellValue?: (value: unknown, context: CellContext) => string
}

const {
  title = '',
  regionLabel = '',
  columns,
  rows,
  emptyText = '',
  emptyCellText = '',
  rowKey,
  stickyHeader = true,
  compactMediaQuery = '(max-width: 768px)',
  formatCellValue,
} = defineProps<Props>()

const attrs = useAttrs()
const slots = useSlots()

const isCompact = useMediaQuery(computed(() => compactMediaQuery))

const hasHeader = computed<boolean>(() => {
  return Boolean(title) || Boolean(slots.title) || Boolean(slots.actions)
})

function toCssSize(value: CssSize | undefined): string | undefined {
  if (value === undefined) return undefined
  return typeof value === 'number' ? `${value}px` : value
}

function getColumnStyle(column: AppTableColumn<TRow>): Record<string, string> {
  if (isCompact.value) return {}
  const style: Record<string, string> = {}
  const width = toCssSize(column.width)
  const minWidth = toCssSize(column.minWidth)
  if (width) style.width = width
  if (minWidth) style.minWidth = minWidth
  return style
}

function getAlignClass(align: AppTableColumn<TRow>['align']): string {
  if (align === 'center') return 'is-center'
  if (align === 'right') return 'is-right'
  return 'is-left'
}

function getCellValue<TKey extends ColumnKey<TRow>>(row: TRow, key: TKey): TRow[TKey] {
  return row[key]
}

function normalizeCellValue(value: unknown): string {
  if (value === null || value === undefined || value === '') return emptyCellText
  if (value instanceof Date) return value.toISOString()
  return String(value)
}

function resolveCellText(row: TRow, column: AppTableColumn<TRow>, rowIndex: number): string {
  const value = getCellValue(row, column.key)
  if (formatCellValue) {
    return formatCellValue(value, { row, rowIndex, column })
  }
  return normalizeCellValue(value)
}

function resolveRowKey(row: TRow, index: number): RowIdentity {
  if (typeof rowKey === 'function') return rowKey(row, index)
  if (typeof rowKey === 'string' && rowKey.length > 0) {
    const candidate = getCellValue(row, rowKey)
    if (typeof candidate === 'string' || typeof candidate === 'number') return candidate
  }
  return index
}
</script>

<style scoped>
.app-table {
  width: 100%;
  background: var(--color-background-glass);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  overflow: hidden;
  backdrop-filter: blur(18px) saturate(140%);
  -webkit-backdrop-filter: blur(18px) saturate(140%);
}

.app-table__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: var(--color-background-elevated);
  border-bottom: 1px solid var(--color-border);
  backdrop-filter: blur(20px) saturate(160%);
  -webkit-backdrop-filter: blur(20px) saturate(160%);
}

.app-table__title {
  min-width: 0;
}

.app-table__title-text {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-table__actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-shrink: 0;
}

.app-table__scroll {
  overflow: auto;
  max-width: 100%;
}

.app-table__table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}

.app-table__table--compact {
  width: max-content;
  min-width: max-content;
  display: inline-table;
}

thead th {
  padding: 0.75rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap;
}

thead.is-sticky th {
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--color-background-elevated);
  backdrop-filter: blur(16px) saturate(160%);
  -webkit-backdrop-filter: blur(16px) saturate(160%);
}

tbody td {
  padding: 0.875rem 1.25rem;
  font-size: 0.875rem;
  color: var(--color-text);
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap;
}

tbody tr:last-child td {
  border-bottom: none;
}

.app-table__row {
  transition: background-color 0.2s ease;
}

.app-table__row:hover {
  background: var(--color-background-secondary);
}

.app-table__empty {
  padding: 2rem 1.25rem;
  text-align: center;
  color: var(--color-text-tertiary);
}

.is-left {
  text-align: left;
}

.is-center {
  text-align: center;
}

.is-right {
  text-align: right;
}

@media (max-width: 768px) {
  .app-table__header {
    padding: 0.75rem 1rem;
    gap: 0.75rem;
  }

  .app-table__actions {
    gap: 0.5rem;
  }

  thead th,
  tbody td {
    padding: 0.75rem 1rem;
    font-size: 0.8125rem;
  }
}
</style>
