import type { AppTableColumn } from '@/components/common'

export interface PoeAccountListItem {
  id: string
  account_name: string
  status: 'active' | 'disabled'
}

export const poeTableColumnsWithActions: ReadonlyArray<AppTableColumn<PoeAccountListItem, '__actions'>> = [
  { key: 'account_name', label: '账号' },
  { key: 'status', label: '状态' },
  { key: '__actions', kind: 'display', label: '操作', fixed: 'right' },
]

export const poeTableColumnsDataOnly: ReadonlyArray<AppTableColumn<PoeAccountListItem>> = [
  { key: 'account_name', label: '账号' },
  { key: 'status', label: '状态' },
  // @ts-expect-error "__actions" requires an explicit extra key generic.
  { key: '__actions', kind: 'display', label: '操作' },
]
