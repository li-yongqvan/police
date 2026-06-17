/** Normalize API / fetch errors to a readable string (never "[object Object]"). */
export function errorMessageFrom(value, fallback = '请求失败') {
  if (value == null || value === '') return fallback
  if (typeof value === 'string') return value
  if (value instanceof Error) return errorMessageFrom(value.message, fallback)
  if (typeof value === 'object') {
    const nested =
      value.message ?? value.error ?? value.detail ?? value.msg ?? value.reason
    if (nested != null && nested !== value) return errorMessageFrom(nested, fallback)
    try {
      const text = JSON.stringify(value)
      if (text && text !== '{}') return text
    } catch {
      /* ignore */
    }
  }
  const text = String(value)
  return text === '[object Object]' ? fallback : text
}

/** Map backend error text to user-facing copy (R3). */
export function formatApiError(error, fallback = '请求失败') {
  const raw = errorMessageFrom(error?.message ?? error, fallback)
  if (/等级|level/i.test(raw)) return '当前账号等级不足，无法执行此操作。'
  if (/敏感|moderation|违规/i.test(raw)) return '内容触发合规检测，请修改后重试。'
  if (/审核|pending|pending_review/i.test(raw)) return '内容已进入审核队列，通过后会展示。'
  if (/登录|token|未登录|过期/i.test(raw)) return '登录已过期，请重新登录。'
  if (/网络|超时|不稳定|连接服务器/i.test(raw)) return raw
  if (/邀请码|invitation|invite/i.test(raw)) {
    if (/无效|不存在|已使用|过期/i.test(raw)) return '邀请码无效或已使用，请向辅导员重新领取。'
    return raw
  }
  if (/429|过于频繁|too many/i.test(raw)) return '操作过于频繁，请稍后再试。'
  return raw
}
