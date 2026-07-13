/** 顶栏 / 侧栏与后端板块的映射（设计规范 v1.1） */
export const GX_NAV_ITEMS = [
  { key: 'study', label: '学业研讨', icon: 'book', keywords: ['学业', '研讨', '学习', '课程'] },
  { key: 'training', label: '警务实训', icon: 'shield', keywords: ['实训', '警务', '实务', '技能'] },
  { key: 'notice', label: '校园公告', icon: 'bell', keywords: ['公告', '通知', '活动'] },
  { key: 'club', label: '社团风采', icon: 'flag', keywords: ['社团', '风采', '文化'] },
]

/** 顶栏主导航（概念稿） */
export const GX_HEADER_NAV = [
  { id: 'home', label: '首页', to: '/community', match: (r) => r.name === 'community-home' },
  {
    id: 'boards',
    label: '板块',
    to: '/community/boards/study',
    match: (r) => r.name === 'board',
  },
  { id: 'circle', label: '校园圈', to: '/community/circle', match: (r) => r.name === 'campus-circle' },
  { id: 'rank', label: '排行榜', to: '/community/rank', match: (r) => r.name === 'rank' },
  { id: 'help', label: '帮助中心', to: '/community/about', match: (r) => r.name === 'about' },
]

/** 校园圈专属板块 slug（与 DB 迁移一致） */
export const CAMPUS_CIRCLE_SLUG = 'campus-circle'

/** 侧栏个人导航（概念稿） */
export const GX_SIDEBAR_PERSONAL = [
  { id: 'home', name: 'community-home', label: '首页', icon: 'home', to: '/community' },
  { id: 'follow', name: 'messages', label: '我的关注', icon: 'flag', to: '/community/messages' },
  { id: 'my-posts', name: 'my-posts', label: '我的帖子', icon: 'edit', to: '/community/my/posts' },
  { id: 'favorites', name: 'my-favorites', label: '我的收藏', icon: 'star', to: '/community/my/favorites' },
  { id: 'history', name: 'my-history', label: '浏览历史', icon: 'clock', to: '/community/my/history' },
]

export const GX_SIDEBAR_PRIMARY = [
  { id: 'home', name: 'community-home', label: '首页', icon: 'home', to: '/community' },
  ...GX_NAV_ITEMS.map((item) => ({
    id: `board-${item.key}`,
    name: `board-${item.key}`,
    label: item.label,
    icon: item.icon,
    to: `/community/boards/${item.key}`,
    boardKey: item.key,
  })),
]

export const GX_SIDEBAR_SECONDARY = [
  { id: 'new-post', name: 'new-post', label: '发帖', icon: 'edit', to: '/community/posts/new' },
  { id: 'messages', name: 'messages', label: '消息', icon: 'message', to: '/community/messages' },
  { id: 'profile', name: 'profile', label: '个人中心', icon: 'user', to: '/community/profile' },
  { id: 'about', name: 'about', label: '关于本站', icon: 'info', to: '/community/about' },
]

export const GX_SIDEBAR_NAV = [...GX_SIDEBAR_PRIMARY, ...GX_SIDEBAR_SECONDARY]

/** @deprecated use GX_SIDEBAR_PRIMARY */
export const GX_TOP_TABS = [
  { name: 'community-home', label: '首页', to: '/community' },
  ...GX_NAV_ITEMS.map((item) => ({
    name: `board-${item.key}`,
    label: item.label,
    to: `/community/boards/${item.key}`,
  })),
  { name: 'profile', label: '个人中心', to: '/community/profile' },
]

export const GX_MOBILE_TABS = [
  { name: 'community-home', label: '首页', to: '/community', icon: 'home' },
  { name: 'board', label: '板块', to: '/community/boards/study', icon: 'book', matchBoard: true },
  { name: 'new-post', label: '发帖', to: '/community/posts/new', icon: 'edit' },
  { name: 'messages', label: '消息', to: '/community/messages', icon: 'message' },
  { name: 'profile', label: '我的', to: '/community/profile', icon: 'user' },
]

export function isNavActive(route, item) {
  if (item.boardKey) {
    return route.name === 'board' && route.params.slug === item.boardKey
  }
  if (item.matchBoard) {
    return route.name === 'board'
  }
  if (item.id === 'home') {
    return route.name === 'community-home'
  }
  return route.name === item.name
}

export function resolveBoardByKey(boards, key) {
  if (!boards?.length || !key) return null
  const direct = boards.find((b) => b.slug === key)
  if (direct) return direct
  const item = GX_NAV_ITEMS.find((n) => n.key === key)
  if (item) {
    const matched = boards.find((b) => item.keywords.some((k) => b.name?.includes(k)))
    if (matched) return matched
  }
  return null
}

export function resolveCircleBoard(boards) {
  if (!boards?.length) return null
  return (
    boards.find((b) => b.slug === CAMPUS_CIRCLE_SLUG) ||
    boards.find((b) => /校园圈|生活|日常/.test(b.name || '')) ||
    null
  )
}

export function boardTagClass(boardName = '') {
  if (/校园圈|生活|日常/.test(boardName)) return 'gx-tag--life'
  if (/社团/.test(boardName)) return 'gx-tag--club'
  if (/公告|通知/.test(boardName)) return 'gx-tag--notice'
  return 'gx-tag--study'
}

export function boardKeyFromName(boardName = '') {
  const item = GX_NAV_ITEMS.find((n) => n.keywords.some((k) => boardName.includes(k)))
  return item?.key || 'study'
}

export const GX_LAST_BOARD_KEY = 'gx-last-board'

export function getLastBoardPath() {
  try {
    const key = localStorage.getItem(GX_LAST_BOARD_KEY) || 'study'
    return `/community/boards/${key}`
  } catch {
    return '/community/boards/study'
  }
}

export function rememberBoardSlug(slug) {
  if (!slug) return
  try {
    localStorage.setItem(GX_LAST_BOARD_KEY, String(slug))
  } catch {
    /* ignore */
  }
}
