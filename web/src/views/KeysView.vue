<template>
  <div class="keys-view">
    <header class="keys-view__header">
      <h1 class="keys-view__title">密钥管理</h1>
      <p class="keys-view__description">模板演示页：前端硬编码数据，使用公共表格与分页组件。</p>
    </header>

    <AppTable
      title="API 密钥列表"
      region-label="API 密钥分页列表"
      :columns="columns"
      :rows="pagedRows"
      :format-cell-value="formatCellValue"
      row-key="id"
      empty-text="暂无密钥数据"
    />

    <AppPagination v-model="currentPage" :total="demoKeys.length" :page-size="pageSize" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import AppPagination from '@/components/common/AppPagination.vue'
import AppTable, { type AppTableColumn } from '@/components/common/AppTable.vue'

defineOptions({
  name: 'KeysView',
})

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

const pageSize = 8
const currentPage = ref<number>(1)

const statusText: Record<ApiKeyStatus, string> = {
  active: '启用',
  disabled: '禁用',
  revoked: '已吊销',
}

const columns: ReadonlyArray<AppTableColumn<ApiKeyItem>> = [
  { key: 'name', label: '名称', minWidth: 180 },
  { key: 'project', label: '项目', minWidth: 140 },
  { key: 'key_masked', label: '密钥', minWidth: 220 },
  { key: 'scopes', label: '权限范围', minWidth: 240 },
  { key: 'status_label', label: '状态', align: 'center', width: 100 },
  { key: 'created_at', label: '创建时间', minWidth: 150 },
  { key: 'last_used_at', label: '最近使用', minWidth: 150 },
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

function formatCellValue(value: unknown, context: { column: AppTableColumn<ApiKeyItem> }): string {
  if (context.column.key === 'status_label') {
    return String(value)
  }
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return String(value)
}
</script>

<style scoped>
.keys-view {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.keys-view__header {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.keys-view__title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
}

.keys-view__description {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 0.95rem;
}

@media (max-width: 768px) {
  .keys-view {
    gap: 1rem;
  }

  .keys-view__title {
    font-size: 1.25rem;
  }
}
</style>
