import auditSeed from './mock-data/audit-records.json'
import boardsSeed from './mock-data/boards.json'
import commentsSeed from './mock-data/comments.json'
import postsSeed from './mock-data/posts.json'
import sensitiveWordsSeed from './mock-data/sensitive-words.json'
import configSeed from './mock-data/system-config.json'
import usersSeed from './mock-data/users.json'

const STORAGE_KEY = 'ai-forum-standalone-state'
const NETWORK_DELAY_MS = 80

function clone(value) {
  return JSON.parse(JSON.stringify(value))
}

function loadState() {
  const cached = localStorage.getItem(STORAGE_KEY)
  if (cached) {
    try {
      return JSON.parse(cached)
    } catch {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  return {
    boards: clone(boardsSeed),
    posts: clone(postsSeed),
    comments: clone(commentsSeed),
    users: clone(usersSeed),
    auditRecords: clone(auditSeed),
    config: clone(configSeed),
    sensitiveWords: clone(sensitiveWordsSeed),
  }
}

const state = loadState()

function persist() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
}

function delay() {
  return new Promise((resolve) => window.setTimeout(resolve, NETWORK_DELAY_MS))
}

function fail(message) {
  throw new Error(message)
}

function ensureUser(userId) {
  const user = state.users.find((item) => item.id === userId)
  if (!user) fail('用户不存在')
  return user
}

function ensurePost(postId) {
  const post = state.posts.find((item) => item.id === postId)
  if (!post) fail('帖子不存在')
  return post
}

function timestampId(prefix) {
  return `${prefix}-${Date.now()}`
}

function nowIso() {
  return new Date().toISOString()
}

function normalizeAttachment(item, index) {
  return {
    id: item.id || `${timestampId('att')}-${index}`,
    type: item.type || 'link',
    name: item.name || `资源 ${index + 1}`,
    url: item.url || `https://example.com/resource/${index + 1}`,
  }
}

function getVisibleBoardIds(includeDisabled) {
  if (includeDisabled) {
    return new Set(state.boards.map((item) => item.id))
  }

  return new Set(
    state.boards
      .filter((item) => state.config.boardSwitches[item.id] !== false)
      .map((item) => item.id),
  )
}

function moderatePost(post) {
  const body = `${post.title} ${post.content}`.toLowerCase()
  const matchedWord = state.sensitiveWords.find((word) => body.includes(String(word).toLowerCase()))

  if (matchedWord) {
    return {
      status: 'pending',
      reason: `命中敏感词：${matchedWord}`,
    }
  }

  if (state.config.moderationMode === 'manual') {
    return {
      status: 'pending',
      reason: '当前为人工审核模式',
    }
  }

  return {
    status: 'published',
    reason: '',
  }
}

function auditList() {
  return state.auditRecords
    .filter((item) => item.status === 'pending')
    .map((item) => {
      const post = state.posts.find((postItem) => postItem.id === item.postId)
      return {
        id: item.id,
        postId: item.postId,
        reason: item.reason,
        createdAt: item.createdAt,
        title: post?.title || '已删除帖子',
        authorId: post?.authorId || '',
        boardId: post?.boardId || '',
      }
    })
}

function computeOverview() {
  const publishedPosts = state.posts.filter((item) => item.status === 'published')
  const pendingPosts = state.posts.filter((item) => item.status === 'pending')
  const todayKey = new Date().toISOString().slice(0, 10)
  let todayPostCount = state.posts.filter((item) => item.createdAt.slice(0, 10) === todayKey).length

  if (!todayPostCount && state.posts.length) {
    const latestDay = [...state.posts]
      .sort((a, b) => String(b.createdAt).localeCompare(String(a.createdAt)))[0]
      .createdAt.slice(0, 10)
    todayPostCount = state.posts.filter((item) => item.createdAt.slice(0, 10) === latestDay).length
  }

  return {
    userCount: state.users.length,
    todayPostCount,
    pendingAuditCount: pendingPosts.length,
    postCount: publishedPosts.length,
    boardActivity: state.boards.map((board) => ({
      boardId: board.id,
      name: board.name,
      count: publishedPosts.filter((item) => item.boardId === board.id).length,
    })),
  }
}

async function resolve(data) {
  await delay()
  return clone(data)
}

export const userApi = {
  async demoLogin(role) {
    const user = state.users.find((item) => item.role === role)
    if (!user) fail('未找到对应角色账号')
    return resolve({
      token: `demo:${user.id}:${Date.now()}`,
      user,
    })
  },

  async me(token) {
    const userId = token?.split(':')[1]
    if (!userId) fail('登录态已失效')
    return resolve(ensureUser(userId))
  },

  async listUsers() {
    return resolve(state.users)
  },
}

export const forumApi = {
  async getBoards(includeDisabled = false) {
    const visibleIds = getVisibleBoardIds(includeDisabled)
    return resolve(state.boards.filter((item) => visibleIds.has(item.id)))
  },

  async getPosts(boardId = '', includePending = false) {
    const visibleIds = getVisibleBoardIds(false)
    const allowedStatuses = includePending ? new Set(['published', 'pending', 'rejected']) : new Set(['published'])
    const filtered = state.posts.filter((item) => {
      if (boardId && item.boardId !== boardId) return false
      if (!allowedStatuses.has(item.status)) return false
      if (!visibleIds.has(item.boardId) && item.status === 'published') return false
      return true
    })
    return resolve(filtered)
  },

  async getPost(id, includeHidden = false) {
    const post = ensurePost(id)
    if (post.status !== 'published' && !includeHidden) {
      fail('帖子暂不可见')
    }

    const comments = state.comments.filter((item) => item.postId === id)
    return resolve({ post, comments })
  },

  async createPost(payload) {
    const author = ensureUser(payload.authorId)
    if (!state.config.postingEnabled) fail('当前已关闭发帖')
    if (state.config.boardSwitches[payload.boardId] === false) fail('该板块当前未开放')
    if (author.status === 'banned') fail('该用户已被封禁，无法发帖')

    const nextPost = {
      id: timestampId('post'),
      boardId: payload.boardId,
      authorId: payload.authorId,
      title: payload.title.trim(),
      content: payload.content.trim(),
      tags: (payload.tags || []).filter(Boolean),
      attachments: (payload.attachments || []).map(normalizeAttachment),
      status: 'published',
      isFeatured: false,
      likeCount: 0,
      commentCount: 0,
      createdAt: nowIso(),
    }

    const moderation = moderatePost(nextPost)
    nextPost.status = moderation.status
    state.posts.unshift(nextPost)

    if (moderation.status === 'pending') {
      state.auditRecords.unshift({
        id: timestampId('audit'),
        postId: nextPost.id,
        reason: moderation.reason,
        status: 'pending',
        reviewerId: '',
        createdAt: nowIso(),
      })
    }

    persist()
    return resolve(nextPost)
  },

  async createComment(id, payload) {
    const post = ensurePost(id)
    ensureUser(payload.authorId)
    if (post.status !== 'published') fail('帖子不存在或不可评论')

    const nextComment = {
      id: timestampId('comment'),
      postId: id,
      authorId: payload.authorId,
      content: payload.content.trim(),
      createdAt: nowIso(),
    }

    state.comments.push(nextComment)
    post.commentCount += 1
    persist()
    return resolve(nextComment)
  },

  async likePost(id) {
    const post = ensurePost(id)
    post.likeCount += 1
    persist()
    return resolve(post)
  },

  async toggleFeature(id) {
    const post = ensurePost(id)
    post.isFeatured = !post.isFeatured
    persist()
    return resolve(post)
  },
}

export const adminApi = {
  async getOverview() {
    return resolve(computeOverview())
  },

  async getConfig() {
    return resolve(state.config)
  },

  async updateConfig(payload) {
    state.config = {
      postingEnabled: Boolean(payload.postingEnabled),
      boardSwitches: { ...payload.boardSwitches },
      moderationMode: payload.moderationMode || 'auto',
    }
    persist()
    return resolve(state.config)
  },

  async getPendingAudit() {
    return resolve(auditList())
  },

  async approveAudit(id, reviewerId) {
    const record = state.auditRecords.find((item) => item.id === id)
    if (!record) fail('审核记录不存在')
    record.status = 'approve'
    record.reviewerId = reviewerId
    ensurePost(record.postId).status = 'published'
    persist()
    return resolve({ success: true })
  },

  async rejectAudit(id, reviewerId) {
    const record = state.auditRecords.find((item) => item.id === id)
    if (!record) fail('审核记录不存在')
    record.status = 'reject'
    record.reviewerId = reviewerId
    ensurePost(record.postId).status = 'rejected'
    persist()
    return resolve({ success: true })
  },

  async deletePost(id) {
    const post = ensurePost(id)
    post.status = 'deleted'
    persist()
    return resolve(post)
  },

  async setUserStatus(id, status) {
    if (!['active', 'banned'].includes(status)) fail('用户状态参数错误')
    const user = ensureUser(id)
    user.status = status
    persist()
    return resolve(user)
  },
}

