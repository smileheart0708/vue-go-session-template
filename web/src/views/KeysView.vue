<template>
  <div class="flex w-full flex-col gap-5 max-md:gap-4">
    <header class="flex flex-col gap-1.5">
      <h1 class="m-0 text-2xl font-bold text-text-primary max-md:text-xl">密钥管理</h1>
      <p class="m-0 text-[0.95rem] text-text-secondary">
        模板演示页：前端硬编码数据，使用公共表格与分页组件。
      </p>
    </header>

    <AppTable
      title="API 密钥列表"
      region-label="API 密钥分页列表"
      :columns="columns"
      :rows="pagedRows"
      :format-cell-value="formatCellValue"
      row-key="id"
      empty-text="暂无密钥数据"
    >
      <template #cell-__actions="{ row }">
        <IconButton
          size="medium"
          title="打开操作菜单"
          class="text-text-secondary hover:text-accent"
          aria-label="打开操作菜单"
          aria-haspopup="menu"
          :aria-expanded="isActionMenuOpenForRow(row.id)"
          @click="handleToggleActionMenu(row, $event)"
        >
          <MoreHorizontal />
        </IconButton>
      </template>
    </AppTable>

    <AppPagination
      v-model="currentPage"
      :total="demoKeys.length"
      :page-size="pageSize"
    />

    <DropdownDrawer
      v-model="showActionMenu"
      :anchor-el="actionMenuAnchorEl"
      :min-width="132"
    >
      <button
        class="dropdown-item"
        type="button"
        @click="handleSelectDetail"
      >
        <Eye class="dropdown-icon" />
        <span>详情</span>
      </button>
      <button
        class="dropdown-item"
        type="button"
        @click="handleSelectDelete"
      >
        <Trash2 class="dropdown-icon" />
        <span>删除</span>
      </button>
    </DropdownDrawer>

    <AppDialog
      v-model="showDetailDialog"
      title="密钥详情"
      close-aria-label="关闭详情对话框"
    >
      <div
        v-if="selectedDetailRow"
        class="grid gap-3 text-sm text-text-primary"
      >
        <p class="m-0">
          <span class="text-text-secondary">名称：</span>{{ selectedDetailRow.name }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">项目：</span>{{ selectedDetailRow.project }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">密钥：</span>{{ selectedDetailRow.key_masked }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">权限范围：</span>{{ selectedDetailRow.scopes }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">状态：</span>{{ selectedDetailRow.status_label }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">创建时间：</span>{{ selectedDetailRow.created_at }}
        </p>
        <p class="m-0">
          <span class="text-text-secondary">最近使用：</span>{{ selectedDetailRow.last_used_at }}
        </p>
      </div>
      <template #footer>
        <BaseButton
          text="关闭"
          :height="36"
          @click="showDetailDialog = false"
        />
      </template>
    </AppDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Eye, MoreHorizontal, Trash2 } from 'lucide-vue-next'
import {
  AppDialog,
  AppPagination,
  AppTable,
  BaseButton,
  DropdownDrawer,
  IconButton,
  type AppTableColumn,
} from '@/components/common'
import { useToast } from '@/composables'

defineOptions({ name: 'KeysView' })

type ApiKeyStatus = 'active' | 'disabled' | 'revoked'

interface ApiKeyItem {
  id: string
  name: string
  project: string
  key_masked: string
  scopes: string
  status: ApiKeyStatus
  status_label: string
  created_at: string
  last_used_at: string
}

const { toast } = useToast()
const pageSize = 8
const currentPage = ref<number>(1)
const showActionMenu = ref<boolean>(false)
const actionMenuAnchorEl = ref<HTMLElement | null>(null)
const actionMenuRow = ref<ApiKeyItem | null>(null)
const showDetailDialog = ref<boolean>(false)
const selectedDetailRow = ref<ApiKeyItem | null>(null)

const statusText: Record<ApiKeyStatus, string> = {
  active: '启用',
  disabled: '禁用',
  revoked: '已吊销',
}

const columns: ReadonlyArray<AppTableColumn<ApiKeyItem, '__actions'>> = [
  { key: 'name', label: '名称', minWidth: 10 },
  { key: 'project', label: '项目', minWidth: 10 },
  { key: 'key_masked', label: '密钥', minWidth: 10 },
  { key: 'scopes', label: '权限范围', minWidth: 10 },
  { key: 'status_label', label: '状态', align: 'center', width: 10 },
  { key: 'created_at', label: '创建时间', minWidth: 10 },
  { key: 'last_used_at', label: '最近使用', minWidth: 10 },
  {
    key: '__actions',
    kind: 'display',
    label: '操作',
    align: 'center',
    width: 10,
    fixed: 'right',
    fixedVisibility: 'always',
  },
]

const projectNames = ['支付服务', '日志采集', '报表中心', '多租户代理'] as const
const scopeLabels = [
  'chat.completions, models.read',
  'embeddings.create, models.read',
  'logs.read, logs.write',
  'models.read',
] as const
const statuses: ReadonlyArray<ApiKeyStatus> = ['active', 'active', 'disabled', 'revoked']

function pickCycledValue<T>(list: ReadonlyArray<T>, index: number): T {
  const item = list[index % list.length]
  if (item === undefined) {
    throw new Error('Demo source list must not be empty.')
  }
  return item
}

const demoKeys: ReadonlyArray<ApiKeyItem> = Array.from({ length: 31 }, (_, index) => {
  const order = index + 1
  const month = String((index % 12) + 1).padStart(2, '0')
  const day = String((index % 27) + 1).padStart(2, '0')
  const status = pickCycledValue(statuses, index)
  const maskedTail = String(8400 + order).padStart(4, '0')

  return {
    id: `key_${String(order).padStart(3, '0')}`,
    name: `演示密钥 ${String(order).padStart(2, '0')}`,
    project: pickCycledValue(projectNames, index),
    key_masked: `sk-live-3a9f${String(2000 + order)}****${maskedTail}`,
    scopes: pickCycledValue(scopeLabels, index),
    status,
    status_label: statusText[status],
    created_at: `2025-${month}-${day}`,
    last_used_at: status === 'revoked' ? '从未使用' : `2026-${month}-${day}`,
  }
})

const pagedRows = computed<ReadonlyArray<ApiKeyItem>>(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return demoKeys.slice(start, end)
})

function isActionMenuOpenForRow(rowId: string): boolean {
  return showActionMenu.value && actionMenuRow.value?.id === rowId
}

function handleToggleActionMenu(row: ApiKeyItem, event: MouseEvent): void {
  const target = event.currentTarget
  if (!(target instanceof HTMLElement)) return

  const isSameRow = actionMenuRow.value?.id === row.id
  if (showActionMenu.value && isSameRow) {
    showActionMenu.value = false
    return
  }

  actionMenuAnchorEl.value = target
  actionMenuRow.value = row
  showActionMenu.value = true
}

function handleSelectDetail(): void {
  if (!actionMenuRow.value) return
  selectedDetailRow.value = actionMenuRow.value
  showDetailDialog.value = true
  showActionMenu.value = false
}

function handleSelectDelete(): void {
  if (!actionMenuRow.value) return
  toast.warn(`演示操作：删除 ${actionMenuRow.value.name}`)
  showActionMenu.value = false
}

function formatCellValue(
  value: unknown,
  context: { column: AppTableColumn<ApiKeyItem, '__actions'> },
): string {
  if (context.column.key === '__actions') {
    return ''
  }

  if (context.column.key === 'status_label') {
    return String(value)
  }

  if (value === null || value === undefined || value === '') {
    return '-'
  }

  return String(value)
}
</script>
