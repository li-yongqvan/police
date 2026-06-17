const STORAGE_KEY = 'gx-browse-history'
const MAX_ITEMS = 50

function readAll() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    const list = raw ? JSON.parse(raw) : []
    return Array.isArray(list) ? list : []
  } catch {
    return []
  }
}

function writeAll(list) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(list.slice(0, MAX_ITEMS)))
  } catch {
    /* ignore quota */
  }
}

export function recordBrowseHistory(post) {
  if (!post?.id) return
  const entry = {
    id: String(post.id),
    title: post.title || '帖子',
    boardName: post.boardName || '',
    boardSlug: post.boardSlug || '',
    authorName: post.authorName || '',
    visitedAt: new Date().toISOString(),
  }
  const next = [entry, ...readAll().filter((item) => item.id !== entry.id)]
  writeAll(next)
}

export function getBrowseHistory() {
  return readAll()
}

export function clearBrowseHistory() {
  localStorage.removeItem(STORAGE_KEY)
}
