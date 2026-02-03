/**
 * 时间格式化工具函数
 */

/**
 * 格式化启动时间为可读的时间字符串
 * 格式：YYYY/MM/DD HH:mm:ss
 * @param startTime - Unix时间戳（秒）
 * @returns 格式化后的时间字符串
 */
export function formatStartTime(startTime: number): string {
  const date = new Date(startTime * 1000)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
}

/**
 * 格式化运行时长
 * 最小单位是秒，后面是分、时、天、月、年
 * @param startTime - Unix时间戳（秒）
 * @param nowMs - 当前时间戳（毫秒），可选，用于外部控制更新
 * @returns 格式化的运行时长字符串
 */
export function formatUptime(startTime: number, nowMs?: number): string {
  const now = nowMs || Date.now()
  const diff = Math.floor((now - startTime * 1000) / 1000) // 转换为秒

  const seconds = diff % 60
  const minutes = Math.floor(diff / 60) % 60
  const hours = Math.floor(diff / 3600) % 24
  const days = Math.floor(diff / 86400) % 30
  const months = Math.floor(diff / 2592000) % 12
  const years = Math.floor(diff / 31536000)

  const parts: string[] = []

  if (years > 0) parts.push(`${years}年`)
  if (months > 0) parts.push(`${months}月`)
  if (days > 0) parts.push(`${days}天`)
  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分`)
  if (seconds > 0 || parts.length === 0) parts.push(`${seconds}秒`)

  // 只返回前两个单位
  return parts.slice(0, 2).join('')
}
