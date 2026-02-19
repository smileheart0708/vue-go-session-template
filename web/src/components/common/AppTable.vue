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
      <table
        class="app-table__table"
        :class="{ 'app-table__table--compact': isCompact, 'app-table__table--fit': fitContent, 'app-table__table--auto-layout': autoLayout }"
      >
        <thead :class="{ 'is-sticky': stickyHeader }">
          <tr>
            <th
              v-for="(column, index) in columns"
              :key="column.key"
              scope="col"
              :class="[
                getAlignClass(column.align),
                getFixedClass(column.fixed),
                column.headerClass,
              ]"
              :style="getColumnStyle(column, index)"
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
              v-for="(column, colIndex) in columns"
              :key="column.key"
              :class="[
                getAlignClass(column.align),
                getFixedClass(column.fixed),
                column.cellClass,
              ]"
              :style="getColumnStyle(column, colIndex)"
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
export type ColumnFixed = 'left' | 'right'

export interface AppTableColumn<TRow extends object = Record<string, unknown>> {
  key: ColumnKey<TRow>
  label: string
  align?: ColumnAlign
  width?: CssSize
  minWidth?: CssSize
  fixed?: ColumnFixed
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
  fitContent?: boolean
  autoLayout?: boolean
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
  fitContent = false,
  autoLayout = false,
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

function getColumnStyle(column: AppTableColumn<TRow>, index: number): Record<string, string> {
  const style: Record<string, string> = {}

  // Handle fixed column positioning
  if (column.fixed && !isCompact.value) {
    const columnsArray = [...columns]

    if (column.fixed === 'right') {
      // Calculate right offset based on subsequent right-fixed columns
      let rightOffset = 0
      for (let i = index + 1; i < columnsArray.length; i++) {
        if (columnsArray[i]?.fixed === 'right') {
          const colWidth = columnsArray[i]?.width ?? 80
          rightOffset += Number.parseInt(String(colWidth), 10) || 80
        }
      }
      if (rightOffset > 0) {
        style.right = `${rightOffset}px`
      }
    } else if (column.fixed === 'left') {
      // Calculate left offset based on previous left-fixed columns
      let leftOffset = 0
      for (let i = 0; i < index; i++) {
        if (columnsArray[i]?.fixed === 'left') {
          const colWidth = columnsArray[i]?.width ?? 80
          leftOffset += Number.parseInt(String(colWidth), 10) || 80
        }
      }
      if (leftOffset > 0) {
        style.left = `${leftOffset}px`
      }
    }
  }

  // Skip width/minWidth styles in compact/fit mode or auto-layout mode
  if (!isCompact.value && !fitContent && !autoLayout) {
    const width = toCssSize(column.width)
    const minWidth = toCssSize(column.minWidth)
    if (width) style.width = width
    if (minWidth) style.minWidth = minWidth
  }

  return style
}

function getAlignClass(align: AppTableColumn<TRow>['align']): string {
  if (align === 'center') return 'is-center'
  if (align === 'right') return 'is-right'
  return 'is-left'
}

function getFixedClass(fixed: AppTableColumn<TRow>['fixed']): string {
  if (fixed === 'right') return 'is-fixed-right'
  if (fixed === 'left') return 'is-fixed-left'
  return ''
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
  background: var(--sys-color-bg-glass);
  border: 1px solid var(--sys-color-border);
  border-radius: 14px;
  overflow: hidden;
  backdrop-filter: blur(12px) saturate(140%);
}

.app-table__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: var(--sys-color-bg-surface);
  border-bottom: 1px solid var(--sys-color-border);
  backdrop-filter: blur(12px) saturate(140%);
}

.app-table__title {
  min-width: 0;
}

.app-table__title-text {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  color: var(--sys-color-text-primary);
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

.app-table__table--fit {
  width: max-content;
  min-width: max-content;
  display: inline-table;
}

.app-table__table--compact {
  width: max-content;
  min-width: max-content;
  display: inline-table;
}

.app-table__table--auto-layout {
  table-layout: auto;
  min-width: 100%;
}

thead th {
  padding: 0.75rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--sys-color-text-secondary);
  border-bottom: 1px solid var(--sys-color-border);
  white-space: nowrap;
}

thead.is-sticky th {
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--sys-color-bg-surface);
  backdrop-filter: blur(12px) saturate(140%);
}

tbody td {
  padding: 0.875rem 1.25rem;
  font-size: 0.875rem;
  color: var(--sys-color-text-primary);
  border-bottom: 1px solid var(--sys-color-border);
  white-space: nowrap;
}

tbody tr:last-child td {
  border-bottom: none;
}

.app-table__row {
  transition: background-color 0.2s ease;
}

.app-table__row:hover {
  background: var(--sys-color-bg-subtle);
}

.app-table__empty {
  padding: 2rem 1.25rem;
  text-align: center;
  color: var(--sys-color-text-tertiary);
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

/* Fixed column styles */
.is-fixed-right,
.is-fixed-left {
  position: sticky;
  z-index: 3;
  background: var(--sys-color-bg-glass);
  backdrop-filter: blur(12px) saturate(140%);
}

.is-fixed-right {
  right: 0;
  border-left: 1px solid var(--sys-color-border);
}

.is-fixed-left {
  left: 0;
  border-right: 1px solid var(--sys-color-border);
}

/* Sticky header should be above fixed columns */
thead.is-sticky .is-fixed-right,
thead.is-sticky .is-fixed-left {
  z-index: 4;
  background: var(--sys-color-bg-surface);
  backdrop-filter: blur(12px) saturate(140%);
}

/* Row hover: fixed columns should maintain background */
.app-table__row:hover .is-fixed-right,
.app-table__row:hover .is-fixed-left {
  background: var(--sys-color-bg-surface);
}

@media (width <= 768px) {
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
