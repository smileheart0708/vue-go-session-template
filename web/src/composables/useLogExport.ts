import type { Ref } from 'vue'
import type { LogEntry } from '@/utils/logs'

export type LogExportType = 'txt' | 'csv' | 'json'
export type { LogEntry } from '@/utils/logs'

interface UseLogExportOptions {
  logs: Ref<LogEntry[]>
  formatMessage: (log: LogEntry) => string
}

export function useLogExport(options: UseLogExportOptions) {
  const { logs, formatMessage } = options

  function formatFileTimestamp(date: Date): string {
    const pad = (value: number): string => value.toString().padStart(2, '0')
    const year = date.getFullYear()
    const month = pad(date.getMonth() + 1)
    const day = pad(date.getDate())
    const hours = pad(date.getHours())
    const minutes = pad(date.getMinutes())
    const seconds = pad(date.getSeconds())
    return `${year}${month}${day}-${hours}${minutes}${seconds}`
  }

  function downloadFile(content: string, mimeType: string, extension: LogExportType) {
    const timestamp = formatFileTimestamp(new Date())
    const filename = `logs-${timestamp}.${extension}`
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  }

  function buildTxtContent(entries: LogEntry[]): string {
    return entries.map((log) => `${log.time} ${formatMessage(log)}`).join('\n')
  }

  function escapeCsvValue(value: string): string {
    const escaped = value.replace(/"/g, '""')
    if (/[",\n]/.test(escaped)) {
      return `"${escaped}"`
    }
    return escaped
  }

  function buildCsvContent(entries: LogEntry[]): string {
    const header = ['time', 'level', 'message', 'attrs']
    const rows = entries.map((log) => {
      const message = formatMessage(log)
      const attrs = log.attrs ? JSON.stringify(log.attrs) : ''
      return [log.time, log.level, message, attrs].map(escapeCsvValue).join(',')
    })
    return [header.join(','), ...rows].join('\n')
  }

  function buildJsonContent(entries: LogEntry[]): string {
    return JSON.stringify(entries, null, 2)
  }

  function exportLogs(type: LogExportType) {
    if (logs.value.length === 0) return
    const entries = [...logs.value]

    if (type === 'txt') {
      downloadFile(buildTxtContent(entries), 'text/plain;charset=utf-8', type)
      return
    }
    if (type === 'csv') {
      downloadFile(buildCsvContent(entries), 'text/csv;charset=utf-8', type)
      return
    }
    downloadFile(buildJsonContent(entries), 'application/json;charset=utf-8', type)
  }

  return { exportLogs }
}
