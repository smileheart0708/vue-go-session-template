/**
 * 格式化数字，添加千位分隔符
 * @param num 数字
 * @returns 格式化后的字符串
 */
export function formatNumber(num: number): string {
  if (num === 0) return '0'
  return num.toLocaleString('zh-CN')
}

/**
 * 格式化字节数为人类可读格式
 * @param bytes 字节数
 * @returns 格式化后的字符串
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'

  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${units[i]}`
}

/**
 * 格式化 Token 数量
 * @param tokens Token 数量
 * @returns 格式化后的字符串
 */
export function formatTokens(tokens: number): string {
  if (tokens === 0) return '0'

  const units = ['', 'K', 'M', 'B', 'T']
  const k = 1000
  const i = Math.floor(Math.log(tokens) / Math.log(k))

  if (i === 0) {
    return formatNumber(tokens)
  }

  return `${(tokens / Math.pow(k, i)).toFixed(2)}${units[i]}`
}
