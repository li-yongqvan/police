import { apiRequest, apiUpload, LOGIN_TIMEOUT_MS, setRefreshToken } from './api/http'
import { formatDisplayDate, formatDisplayTime } from './utils/displayName'

const authHeaders = (json = true, tokenOverride) => {
  const headers = {}
  if (json) headers['Content-Type'] = 'application/json'
  const token = tokenOverride ?? (localStorage.getItem('ai-forum-token') || '')
  if (token) headers.Authorization = `Bearer ${token}`
  return headers
}

function request(base, path, options = {}) {
  return apiRequest(base, path, options)
}

function userAssetUrl(path) {
  if (!path) return ''
  const value = String(path)
  if (value.startsWith('http') || value.startsWith('blob:') || value.startsWith('data:')) return value
  if (value.startsWith('/user-api/')) return value
  if (value.startsWith('/uploads/')) return `/user-api${value}`
  return value
}

function mapBoard(board) {
  return {
    id: String(board.id),
    name: board.name,
    slug: board.slug,
    description: board.description,
    enabled: board.enabled,
    sortOrder: board.sort_order ?? 0,
    postCount: board.post_count ?? 0,
  }
}

function mapPostListItem(post) {
  return {
    id: String(post.id),
    title: post.title,
    content: post.content,
    boardId: String(post.board_id),
    boardName: post.board_name,
    boardSlug: post.board_slug || '',
    authorId: String(post.author_id),
    authorName: post.author_name,
    authorAvatar: userAssetUrl(post.author_avatar || post.author_avatar_url || ''),
    status: post.status,
    isFeatured: post.is_featured,
    isPinned: post.is_pinned,
    likeCount: post.like_count ?? 0,
    dislikeCount: post.dislike_count ?? 0,
    liked: !!post.liked,
    disliked: !!post.disliked,
    commentCount: post.comment_count ?? 0,
    trendDelta: post.trend_delta ?? post.rank_delta ?? post.hot_rank_delta ?? null,
    trendDirection: post.trend_direction ?? post.rank_trend ?? '',
    hotScore: post.hot_score ?? post.hotScore ?? null,
    createdAtIso: post.created_at,
    createdAt: formatDisplayDate(post.created_at),
    tags: post.tags || [],
  }
}

function mapPostDetail(post, comments = []) {
  return {
    post: {
      ...mapPostListItem(post),
      content: post.content,
      attachments: (post.attachments || []).map((item) => ({
        id: String(item.id),
        type: item.file_type,
        title: item.filename || item.title || '附件',
        url: item.file_type === 'link' ? userAssetUrl(item.file_path) : userAssetUrl(item.id),
        fileSize: item.file_size ?? 0,
      })),
      tags: post.tags || [],
      collected: !!post.collected,
    },
    comments: comments.map((item) => ({
      id: String(item.id),
      postId: String(item.post_id),
      parentId: item.parent_id ? String(item.parent_id) : '',
      depth: item.depth ?? 0,
      authorId: String(item.author_id),
      authorName: item.author_name || item.author_id,
      authorAvatar: userAssetUrl(item.author_avatar || item.author_avatar_url || ''),
      content: item.content,
      likeCount: item.like_count ?? 0,
      dislikeCount: item.dislike_count ?? 0,
      liked: !!item.liked,
      disliked: !!item.disliked,
      createdAtIso: item.created_at,
      createdAt: formatDisplayTime(item.created_at),
    })),
  }
}

function mapNotification(row) {
  return {
    id: String(row.id),
    type: row.type,
    title: row.title,
    body: row.body,
    related_post_id: row.related_post_id ? String(row.related_post_id) : '',
    is_read: row.is_read,
    createdAtIso: row.created_at,
    created_at: formatDisplayTime(row.created_at),
  }
}

function mapUser(user, roleFallback) {
  return {
    id: String(user.id),
    username: user.username || '',
    name: user.name || user.nickname || user.username,
    avatar: userAssetUrl(user.avatar),
    role: user.role || roleFallback || 'student',
    level: user.level ?? 1,
    department: user.department || '',
    squad: user.squad || '',
    grade: user.grade || '',
    profileCompleted: user.profile_completed ?? user.profileCompleted ?? false,
    bio: user.bio || '',
    status: user.status || 'active',
  }
}

export const userApi = {
  async me() {
    const payload = await request('/user-api', '/api/v1/users/me', {
      headers: authHeaders(),
    })
    return mapUser(payload.user || payload)
  },

  async login(username, password) {
    const payload = await request('/user-api', '/api/v1/login', {
      method: 'POST',
      headers: authHeaders(false),
      body: JSON.stringify({ username, password }),
      timeoutMs: LOGIN_TIMEOUT_MS,
    })
    return {
      token: payload.access_token,
      user: mapUser(payload.user),
      refreshToken: payload.refresh_token || '',
    }
  },

  async demoLogin(demoKey) {
    const payload = await request('/user-api', '/api/v1/demo-login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ demo_key: demoKey }),
      timeoutMs: LOGIN_TIMEOUT_MS,
    })
    return {
      token: payload.access_token,
      user: mapUser(payload.user),
      refreshToken: payload.refresh_token || '',
    }
  },

  async register(code, username, password) {
    const payload = await request('/user-api', '/api/v1/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ invite_code: code, username, password }),
      timeoutMs: LOGIN_TIMEOUT_MS,
    })
    return {
      token: payload.access_token,
      user: mapUser(payload.user),
      refreshToken: payload.refresh_token || '',
    }
  },

  async getPublicUser(id) {
    const payload = await request('/user-api', `/api/v1/users/public/${id}`, {
      headers: authHeaders(),
    })
    return mapUser(payload.user || payload)
  },

  async updateMe(data) {
    const payload = await request('/user-api', '/api/v1/users/me', {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify(data),
    })
    return mapUser(payload.user || payload)
  },

  async uploadAvatar(file) {
    const form = new FormData()
    form.append('file', file)
    const payload = await apiUpload('/user-api', '/api/v1/users/me/avatar', form)
    return userAssetUrl(payload.url || payload.avatar || '')
  },

  async follow(userId) {
    await request('/user-api', `/api/v1/users/${userId}/follow`, {
      method: 'POST',
      headers: authHeaders(),
    })
  },

  async unfollow(userId) {
    await request('/user-api', `/api/v1/users/${userId}/follow`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },

  async getFollowing(limit = 20, offset = 0) {
    const payload = await request('/user-api', `/api/v1/users/me/following?limit=${limit}&offset=${offset}`, {
      headers: authHeaders(),
    })
    return payload
  },

  async getFollowers(limit = 20, offset = 0) {
    const payload = await request('/user-api', `/api/v1/users/me/followers?limit=${limit}&offset=${offset}`, {
      headers: authHeaders(),
    })
    return payload
  },

  async getPublicFollowing(userId, limit = 20, offset = 0) {
    const payload = await request('/user-api', `/api/v1/users/${userId}/following?limit=${limit}&offset=${offset}`, {
      headers: authHeaders(),
    })
    return payload
  },

  async getPublicFollowers(userId, limit = 20, offset = 0) {
    const payload = await request('/user-api', `/api/v1/users/${userId}/followers?limit=${limit}&offset=${offset}`, {
      headers: authHeaders(),
    })
    return payload
  },

  async getFollowCounts(userId) {
    const payload = await request('/user-api', `/api/v1/users/${userId}/follow-counts`, {
      headers: authHeaders(),
    })
    return payload
  },

  async listUsers(page = 1, limit = 20) {
    const payload = await request(
      '/admin-api',
      `/api/v1/admin/users?limit=${limit}&page=${page}`,
      { headers: authHeaders() },
    )
    return {
      users: (payload.users || []).map((user) => mapUser(user)),
      total: payload.total ?? 0,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },

  async logout() {
    await request('/user-api', '/api/v1/auth/logout', {
      method: 'POST',
      headers: authHeaders(),
    })
  },
}

export const adminApi = {
  async getOverview() {
    const auth = authHeaders()
    const [statsPayload, pendingPayload, boards] = await Promise.all([
      request('/admin-api', '/api/v1/admin/stats/overview', { headers: auth }),
      request('/admin-api', '/api/v1/admin/audit/pending', { headers: auth }).catch(() => ({
        posts: [],
      })),
      request('/admin-api', '/api/v1/admin/boards', { headers: auth }).then(payload => {
        const list = payload.boards || payload || []
        return Array.isArray(list) ? list.map(mapBoard) : []
      }),
    ])
    const data = statsPayload.data || statsPayload
    const pending = pendingPayload.posts || []
    return {
      userCount: data.total_users ?? 0,
      postCount: data.total_posts ?? 0,
      commentCount: data.total_comments ?? 0,
      todayPostCount: data.posts_today ?? 0,
      pendingAuditCount: pending.length,
      boardActivity: boards.map((board) => ({
        boardId: board.id,
        name: board.name,
        count: board.postCount ?? 0,
      })),
    }
  },

  async getDailyStats(days = 7) {
    const payload = await request('/admin-api', `/api/v1/admin/stats/daily?days=${days}`, {
      headers: authHeaders(),
    })
    return payload.data || payload || []
  },

  async getConfig() {
    const payload = await request('/admin-api', '/api/v1/admin/config', {
      headers: authHeaders(),
    })
    const configs = payload.configs || {}
    return {
      postingEnabled: configs.post_requires_level !== '99',
      moderationMode: configs.sensitive_word_action === 'pending_review' ? 'manual' : 'auto',
      boardSwitches: {},
    }
  },

  async updateConfig(config) {
    const updates = []
    if (typeof config.postingEnabled === 'boolean') {
      updates.push(['post_requires_level', config.postingEnabled ? '0' : '99'])
    }
    if (config.moderationMode) {
      updates.push([
        'sensitive_word_action',
        config.moderationMode === 'manual' ? 'pending_review' : 'reject',
      ])
    }
    if (updates.length) {
      await request('/admin-api', '/api/v1/admin/config', {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(Object.fromEntries(updates)),
      })
    }
    return adminApi.getConfig()
  },

  async listAuditPosts(page = 1, limit = 20) {
    const payload = await request(
      '/admin-api',
      `/api/v1/admin/audit?limit=${limit}&page=${page}`,
      { headers: authHeaders() },
    )
    const posts = payload.posts || payload.pending || []
    return {
      posts: posts.map(mapPostListItem),
      total: payload.total ?? 0,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },

  async approveAudit(id) {
    await request('/admin-api', `/api/v1/admin/audit/${id}/approve`, {
      method: 'POST',
      headers: authHeaders(),
    })
  },

  async rejectAudit(id, reason = '') {
    await request('/admin-api', `/api/v1/admin/audit/${id}/reject`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ reason }),
    })
  },

  async listPosts({ page = 1, limit = 20, boardId = '', status = '' } = {}) {
    const params = new URLSearchParams({ page: String(page), limit: String(limit) })
    if (boardId) params.set('board_id', boardId)
    if (status) params.set('status', status)
    const payload = await request('/admin-api', `/api/v1/admin/posts?${params}`, {
      headers: authHeaders(),
    })
    return {
      posts: (payload.posts || []).map(mapPostListItem),
      total: payload.total ?? 0,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },

  async getPostDetail(id) {
    const payload = await request('/admin-api', `/api/v1/admin/posts/${id}`, {
      headers: authHeaders(),
    })
    return mapPostDetail(payload.post || payload, payload.comments || [])
  },

  async deletePost(id) {
    await request('/admin-api', `/api/v1/admin/posts/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },

  async setPostFeatured(id, featured) {
    await request('/admin-api', `/api/v1/admin/posts/${id}/featured`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ featured }),
    })
  },

  async setPostPinned(id, pinned) {
    await request('/admin-api', `/api/v1/admin/posts/${id}/pinned`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ pinned }),
    })
  },

  async listBoards() {
    const payload = await request('/admin-api', '/api/v1/admin/boards', {
      headers: authHeaders(),
    })
    const boards = payload.boards || payload || []
    return Array.isArray(boards) ? boards.map(mapBoard) : []
  },

  async createBoard(data) {
    await request('/admin-api', '/api/v1/admin/boards', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({
        name: data.name,
        slug: data.slug,
        description: data.description,
        sort_order: Number(data.sortOrder) || 0,
      }),
    })
  },

  async updateBoard(id, data) {
    await request('/admin-api', `/api/v1/admin/boards/${id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        name: data.name,
        slug: data.slug,
        description: data.description,
        sort_order: Number(data.sortOrder) || 0,
        enabled: data.enabled,
      }),
    })
  },

  async deleteBoard(id) {
    await request('/admin-api', `/api/v1/admin/boards/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },

  async setUserStatus(id, status) {
    if (status === 'banned') {
      await request('/admin-api', `/api/v1/admin/users/${id}/ban`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ reason: '演示封禁' }),
      })
    } else {
      await request('/admin-api', `/api/v1/admin/users/${id}/unban`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({}),
      })
    }
  },

  async updateUserLevel(id, level) {
    await request('/admin-api', `/api/v1/admin/users/${id}/level`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ level: Number(level) }),
    })
  },

  async getUserLogs(id) {
    const payload = await request('/admin-api', `/api/v1/admin/users/${id}/logs?limit=20`, {
      headers: authHeaders(),
    })
    return payload.logs || []
  },

  async listInviteCodes() {
    const payload = await request('/admin-api', '/api/v1/admin/invite-codes?limit=50', {
      headers: authHeaders(),
    })
    return payload.codes || payload.invite_codes || []
  },

  async getInviteCodeStatus(code) {
    return request(
      '/admin-api',
      `/api/v1/admin/invite-codes/${encodeURIComponent(code)}/status`,
      { headers: authHeaders() },
    )
  },

  async generateInviteCode() {
    const payload = await request('/admin-api', '/api/v1/admin/invite-codes', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({}),
    })
    return payload.code
  },

  async generateInviteBatch(count) {
    const payload = await request('/admin-api', '/api/v1/admin/invite-codes/batch', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ count }),
    })
    return payload.codes || []
  },

  async voidInviteCode(code) {
    await request('/admin-api', `/api/v1/admin/invite-codes/${encodeURIComponent(code)}/void`, {
      method: 'PUT',
      headers: authHeaders(),
    })
  },

  async listSensitiveWords() {
    const payload = await request('/admin-api', '/api/v1/admin/sensitive-words', {
      headers: authHeaders(),
    })
    return Array.isArray(payload) ? payload : payload.words || []
  },

  async addSensitiveWord(word, category) {
    await request('/admin-api', '/api/v1/admin/sensitive-words', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ word, category: category || 'general' }),
    })
  },

  async deleteSensitiveWord(id) {
    await request('/admin-api', `/api/v1/admin/sensitive-words/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },

  async listRoles() {
    const payload = await request('/admin-api', '/api/v1/admin/roles', {
      headers: authHeaders(),
    })
    return payload.roles || []
  },

  async getUserRoles(userId) {
    const payload = await request('/admin-api', `/api/v1/admin/users/${userId}/roles`, {
      headers: authHeaders(),
    })
    return payload.roles || []
  },

  async assignRole(userId, roleId) {
    await request('/admin-api', `/api/v1/admin/users/${userId}/roles`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ role_id: Number(roleId) }),
    })
  },

  async removeRole(userId, roleId) {
    await request('/admin-api', `/api/v1/admin/users/${userId}/roles/${roleId}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },

  async listReports(page = 1, limit = 20) {
    const payload = await request(
      '/admin-api',
      `/api/v1/admin/reports?limit=${limit}&page=${page}`,
      { headers: authHeaders() },
    )
    return {
      reports: payload.reports || [],
      total: payload.total ?? 0,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },

  async resolveReport(id, action, reason = '') {
    await request('/admin-api', `/api/v1/admin/reports/${id}/resolve`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ action, reason }),
    })
  },
}
