/** 展示用学号/队别昵称（合规：不使用「警号」） */
export function formatAuthorLabel(user, post) {
  if (post?.authorName) return maskStudentLabel(post.authorName)
  if (user?.name) return maskStudentLabel(user.name)
  if (user?.username) return maskStudentLabel(user.username)
  return '匿名同学'
}

export function maskStudentLabel(raw) {
  const text = String(raw || '').trim()
  if (!text) return '匿名同学'
  if (text.includes('·')) return text
  if (/^\d{2,}/.test(text)) return text.replace(/(\d{2})(\d+)/, '$1****')
  return text.length > 8 ? `${text.slice(0, 6)}…` : text
}

export function formatDisplayTime(iso) {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return String(iso)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** 列表/卡片用：仅显示年月日 */
/** Feed 用相对时间 */
export function formatRelativeTime(iso) {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''
  const diffMs = Date.now() - date.getTime()
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days} 天前`
  return formatDisplayDate(iso)
}

export function formatDisplayDate(iso) {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) {
    const match = String(iso).match(/^(\d{4})-(\d{2})-(\d{2})/)
    if (match) return `${match[1]}-${match[2]}-${match[3]}`
    return String(iso)
  }
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}
