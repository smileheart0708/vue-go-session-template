export interface LogEntry {
  time: string
  level: string
  msg: string
  attrs?: Record<string, string | number | boolean>
}

export function getLevelClass(level: string): string {
  const upperLevel = level.toUpperCase()
  if (upperLevel === 'ERROR') return 'error'
  if (upperLevel === 'WARN') return 'warn'
  if (upperLevel === 'DEBUG') return 'debug'
  return 'info'
}

export function formatLogMessage(log: LogEntry): string {
  let message = log.msg

  if (log.attrs && Object.keys(log.attrs).length > 0) {
    const attrStrings = Object.entries(log.attrs).map(([key, value]) => {
      if (typeof value === 'string') {
        return `${key}=${value}`
      }
      return `${key}=${JSON.stringify(value)}`
    })
    message += ' ' + attrStrings.join(' ')
  }

  return message
}
