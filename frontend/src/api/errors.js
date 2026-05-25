/** Map backend error text to user-facing copy (R3). */
export function formatApiError(error, fallback = '请求失败') {
  const raw = error?.message || String(error || fallback)
  if (/等级|level/i.test(raw)) return '当前账号等级不足，无法执行此操作。'
  if (/敏感|moderation|违规/i.test(raw)) return '内容触发合规检测，请修改后重试。'
  if (/审核|pending|pending_review/i.test(raw)) return '内容已进入审核队列，通过后会展示。'
  if (/登录|token|未登录|过期/i.test(raw)) return '登录已过期，请重新登录。'
  if (/邀请码|invitation/i.test(raw)) return raw
  return raw
}
