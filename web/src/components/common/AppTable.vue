<template>
  <section
    class="w-full overflow-hidden rounded-[14px] border border-border bg-bg-glass backdrop-blur-[12px] backdrop-saturate-[140%]"
    v-bind="attrs"
  >
    <header
      v-if="hasHeader"
      class="flex items-center justify-between gap-4 border-b border-border bg-bg-surface px-5 py-4 backdrop-blur-[12px] backdrop-saturate-[140%] max-md:gap-3 max-md:px-4 max-md:py-3"
    >
      <div class="min-w-0">
        <slot name="title">
          <h2
            v-if="title"
            class="m-0 overflow-hidden text-ellipsis whitespace-nowrap text-lg font-semibold text-text-primary"
          >
            {{ title }}
          </h2>
        </slot>
      </div>

      <div v-if="$slots['actions']" class="flex shrink-0 items-center gap-3 max-md:gap-2">
        <slot name="actions" />
      </div>
    </header>

    <div class="max-w-full overflow-auto" role="region" :aria-label="regionLabel || undefined">
      <table :class="tableClass">
        <thead>
          <tr>
            <th
              v-for="(column, index) in columns"
              :key="column.key"
              scope="col"
              :class="[getHeaderCellClass(column), column.headerClass]"
              :style="getColumnStyle(column, index)"
            >
              <slot :name="`header-${column.key}`" :column="column">
                {{ column.label }}
              </slot>
            </th>
          </tr>
        </thead>

        <tbody v-if="rows.length > 0">
          <tr
            v-for="(row, rowIndex) in rows"
            :key="resolveRowKey(row, rowIndex)"
            class="group transition-colors duration-200 hover:bg-bg-subtle last:[&>td]:border-b-0"
          >
            <td
              v-for="(column, colIndex) in columns"
              :key="column.key"
              :class="[getBodyCellClass(column), column.cellClass]"
              :style="getColumnStyle(column, colIndex)"
            >
              <slot
                :name="`cell-${column.key}`"
                :row="row"
                :column="column"
                :value="getCellValue(row, column, rowIndex)"
                :row-index="rowIndex"
              >
                {{ resolveCellText(row, column, rowIndex) }}
              </slot>
            </td>
          </tr>
        </tbody>

        <tbody v-else>
          <tr>
            <td
              class="px-5 py-8 text-center text-text-tertiary max-md:px-4 max-md:py-6"
              :colspan="columns.length"
            >
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
export type ColumnFixedVisibility = 'always' | 'desktop'

interface AppTableColumnBase {
  label: string
  align?: ColumnAlign
  width?: CssSize
  minWidth?: CssSize
  fixed?: ColumnFixed
  fixedVisibility?: ColumnFixedVisibility
  headerClass?: string
  cellClass?: string
}

export interface AppTableDataColumn<TRow extends object> extends AppTableColumnBase {
  kind?: 'data'
  key: ColumnKey<TRow>
}

export type AppTableDisplayValueResolver<TRow extends object> = (
  row: TRow,
  rowIndex: number,
) => unknown

export interface AppTableDisplayColumn<TRow extends object, TExtraKey extends string>
  extends AppTableColumnBase {
  kind: 'display'
  key: TExtraKey
  value?: AppTableDisplayValueResolver<TRow>
}

export type AppTableColumn<
  TRow extends object = Record<string, unknown>,
  TExtraKey extends string = never,
> = AppTableDataColumn<TRow> | AppTableDisplayColumn<TRow, TExtraKey>

// Backward compatible alias for existing imports.
export type GlassTableColumn<
  TRow extends object = Record<string, unknown>,
  TExtraKey extends string = never,
> = AppTableColumn<TRow, TExtraKey>
</script>

<script setup lang="ts" generic="TRow extends object, TExtraKey extends string = never">
import { computed, useAttrs, useSlots } from 'vue'
import { useMediaQuery } from '@vueuse/core'

defineOptions({ name: 'AppTable', inheritAttrs: false })

type RowIdentity = string | number

interface CellContext {
  row: TRow
  rowIndex: number
  column: AppTableColumn<TRow, TExtraKey>
}

interface Props {
  title?: string
  regionLabel?: string
  columns: ReadonlyArray<AppTableColumn<TRow, TExtraKey>>
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

const tableClass = computed<string>(() => {
  if (autoLayout) {
    return 'w-full min-w-full table-auto border-separate border-spacing-0'
  }
  if (fitContent || isCompact.value) {
    return 'inline-table w-max min-w-max border-separate border-spacing-0'
  }
  return 'w-full border-separate border-spacing-0'
})

const hasHeader = computed<boolean>(() => {
  return Boolean(title) || Boolean(slots['title']) || Boolean(slots['actions'])
})

const alignClassMap: Record<ColumnAlign, string> = {
  left: 'text-left',
  center: 'text-center',
  right: 'text-right',
}

const DEFAULT_FIXED_COLUMN_WIDTH = '80px'

function isDisplayColumn(
  column: AppTableColumn<TRow, TExtraKey>,
): column is AppTableDisplayColumn<TRow, TExtraKey> {
  return column.kind === 'display'
}

function toCssSize(value: CssSize | undefined): string | undefined {
  if (value === undefined) return undefined
  return typeof value === 'number' ? `${value}px` : value
}

function resolveAlignClass(align: ColumnAlign | undefined): string {
  if (align === undefined) {
    return alignClassMap.left
  }
  return alignClassMap[align]
}

function resolveFixedVisibility(column: AppTableColumn<TRow, TExtraKey>): ColumnFixedVisibility {
  return column.fixedVisibility ?? 'always'
}

function shouldPinColumn(column: AppTableColumn<TRow, TExtraKey>): boolean {
  if (!column.fixed) return false
  const visibility = resolveFixedVisibility(column)
  if (visibility === 'desktop' && isCompact.value) return false
  return true
}

function resolveFixedWidth(column: AppTableColumn<TRow, TExtraKey>): string {
  return toCssSize(column.width) ?? toCssSize(column.minWidth) ?? DEFAULT_FIXED_COLUMN_WIDTH
}

function composeCalc(lengths: ReadonlyArray<string>): string | undefined {
  if (lengths.length === 0) return undefined
  if (lengths.length === 1) return lengths[0]
  return `calc(${lengths.join(' + ')})`
}

function resolveLeftOffset(index: number): string | undefined {
  const widths: string[] = []
  for (let cursor = 0; cursor < index; cursor += 1) {
    const column = columns[cursor]
    if (!column || !shouldPinColumn(column) || column.fixed !== 'left') continue
    widths.push(resolveFixedWidth(column))
  }
  return composeCalc(widths)
}

function resolveRightOffset(index: number): string | undefined {
  const widths: string[] = []
  for (let cursor = index + 1; cursor < columns.length; cursor += 1) {
    const column = columns[cursor]
    if (!column || !shouldPinColumn(column) || column.fixed !== 'right') continue
    widths.push(resolveFixedWidth(column))
  }
  return composeCalc(widths)
}

function getColumnStyle(
  column: AppTableColumn<TRow, TExtraKey>,
  index: number,
): Record<string, string> {
  const style: Record<string, string> = {}
  const shouldPin = shouldPinColumn(column)

  if (shouldPin && column.fixed) {
    const offset = column.fixed === 'left' ? resolveLeftOffset(index) : resolveRightOffset(index)
    if (offset) {
      style[column.fixed] = offset
    }
  }

  const shouldApplyWidths = (!isCompact.value && !fitContent && !autoLayout) || shouldPin
  if (shouldApplyWidths) {
    const width = toCssSize(column.width)
    const minWidth = toCssSize(column.minWidth)
    if (width) style['width'] = width
    if (minWidth) style['minWidth'] = minWidth
  }

  return style
}

function getPinnedCellClass(column: AppTableColumn<TRow, TExtraKey>, header: boolean): string {
  if (!shouldPinColumn(column) || !column.fixed) return ''

  const sideClass = column.fixed === 'right' ? 'right-0 border-l border-border' : 'left-0 border-r border-border'
  const sharedClass = 'sticky backdrop-blur-[12px] backdrop-saturate-[140%]'
  if (header) {
    return `${sharedClass} ${sideClass} z-[4] bg-bg-surface`
  }
  return `${sharedClass} ${sideClass} z-[3] bg-bg-glass group-hover:bg-bg-subtle`
}

function getHeaderCellClass(column: AppTableColumn<TRow, TExtraKey>): string {
  const classes = [
    'px-5 py-3 text-sm font-semibold whitespace-nowrap border-b border-border text-text-secondary max-md:px-4 max-md:py-3 max-md:text-[13px]',
    resolveAlignClass(column.align),
  ]

  if (stickyHeader) {
    classes.push('sticky top-0 z-[2] bg-bg-surface backdrop-blur-[12px] backdrop-saturate-[140%]')
  }

  const pinnedClass = getPinnedCellClass(column, true)
  if (pinnedClass) {
    classes.push(pinnedClass)
  }

  return classes.join(' ')
}

function getBodyCellClass(column: AppTableColumn<TRow, TExtraKey>): string {
  const classes = [
    'px-5 py-3.5 text-sm whitespace-nowrap border-b border-border text-text-primary max-md:px-4 max-md:py-3 max-md:text-[13px]',
    resolveAlignClass(column.align),
  ]

  const pinnedClass = getPinnedCellClass(column, false)
  if (pinnedClass) {
    classes.push(pinnedClass)
  }

  return classes.join(' ')
}

function getDataCellValue<TKey extends ColumnKey<TRow>>(row: TRow, key: TKey): TRow[TKey] {
  return row[key]
}

function getCellValue(
  row: TRow,
  column: AppTableColumn<TRow, TExtraKey>,
  rowIndex: number,
): unknown {
  if (isDisplayColumn(column)) {
    return column.value ? column.value(row, rowIndex) : undefined
  }
  return getDataCellValue(row, column.key)
}

function normalizeCellValue(value: unknown): string {
  if (value === null || value === undefined || value === '') return emptyCellText
  if (value instanceof Date) return value.toISOString()
  return String(value)
}

function resolveCellText(
  row: TRow,
  column: AppTableColumn<TRow, TExtraKey>,
  rowIndex: number,
): string {
  const value = getCellValue(row, column, rowIndex)
  if (formatCellValue) {
    return formatCellValue(value, { row, rowIndex, column })
  }
  return normalizeCellValue(value)
}

function resolveRowKey(row: TRow, index: number): RowIdentity {
  if (typeof rowKey === 'function') return rowKey(row, index)
  if (typeof rowKey === 'string' && rowKey.length > 0) {
    const candidate = getDataCellValue(row, rowKey)
    if (typeof candidate === 'string' || typeof candidate === 'number') return candidate
  }
  return index
}
</script>
