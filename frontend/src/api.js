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

function viewerRequest(base, path, options = {}) {
  const token = localStorage.getItem('ai-forum-token')
  if (token) {
    options.headers = { ...authHeaders(), ...options.headers }
  }
  return request(base, path, options)
}

function attachmentUrl(pathOrId) {
  if (!pathOrId) return '#'
  if (String(pathOrId).startsWith('http')) return pathOrId
  return `/forum-api/api/v1/attachments/${pathOrId}`
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
        url: item.file_type === 'link' ? attachmentUrl(item.file_path) : attachmentUrl(item.id),
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

function mapConfigFromBackend(configs = {}, boards = []) {
  const boardSwitches = {}
  boards.forEach((board) => {
    boardSwitches[board.id] = configs[`board_${board.slug}_enabled`] !== 'false'
  })
  return {
    postingEnabled: configs.post_requires_level !== '99',
    moderationMode: configs.sensitive_word_action === 'pending_review' ? 'manual' : 'auto',
    boardSwitches,
    _raw: configs,
  }
}

function persistTokens(payload) {
  const token = payload.access_token || payload.token
  if (token) localStorage.setItem('ai-forum-token', token)
  if (payload.refresh_token) setRefreshToken(payload.refresh_token)
  return token
}

export const userApi = {
  async register({ username, password, invitationCode }) {
    const payload = await request('/user-api', '/api/v1/register', {
      method: 'POST',
      headers: authHeaders(true, null),
      body: JSON.stringify({
        username,
        password,
        invitation_code: invitationCode?.trim().toUpperCase(),
      }),
    })
    const token = persistTokens(payload)
    if (token && payload.user) {
      return { token, user: mapUser(payload.user), refreshToken: payload.refresh_token }
    }
    return userApi.login(username, password)
  },
  async login(username, password) {
    const payload = await request('/user-api', '/api/v1/login', {
      method: 'POST',
      headers: authHeaders(true, null),
      body: JSON.stringify({ username, password }),
      timeoutMs: LOGIN_TIMEOUT_MS,
      retries: 2,
      retryOnPost: true,
    })
    const token = persistTokens(payload)
    const user = await userApi.me()
    return { token, user, refreshToken: payload.refresh_token }
  },
  async demoLogin(role) {
    if (!import.meta.env.DEV) throw new Error('demo-login 仅用于开发环境')
    const payload = await request('/user-api', '/api/v1/demo-login', {
      method: 'POST',
      headers: authHeaders(true, null),
      body: JSON.stringify({ role }),
    })
    const token = persistTokens(payload)
    return { token, user: mapUser(payload.user, role) }
  },
  async me() {
    const payload = await request('/user-api', '/api/v1/users/me', {
      headers: authHeaders(),
    })
    return mapUser(payload)
  },
  async getProfile(id) {
    const payload = await request('/user-api', `/api/v1/users/${id}`, {
      headers: authHeaders(),
    })
    return mapUser(payload)
  },
  async updateProfile(id, data) {
    const payload = await request('/user-api', `/api/v1/users/${id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({
        nickname: data.name,
        bio: data.bio,
        username: data.username,
        profile_completed: data.profileCompleted,
      }),
    })
    return mapUser({ ...payload, role: (await userApi.me()).role })
  },
  async uploadAvatar(id, file) {
    const form = new FormData()
    form.append('avatar', file)
    const payload = await apiUpload('/user-api', `/api/v1/users/${id}/avatar`, form)
    if (payload?.user) return mapUser(payload.user)
    if (payload?.id || payload?.username) return mapUser(payload)
    return userApi.me()
  },
  // Follow system
  async follow(targetUserId) {
    const payload = await request('/user-api', '/api/v1/users/me/following', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ target_user_id: targetUserId }),
    })
    return payload
  },
  async unfollow(targetUserId) {
    await request('/user-api', `/api/v1/users/me/following/${targetUserId}`, {
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
}

export const forumApi = {
  async getBoards(includeDisabled = false) {
    const boards = await request('/forum-api', '/api/v1/boards')
    const list = Array.isArray(boards) ? boards.map(mapBoard) : []
    return includeDisabled ? list : list.filter((item) => item.enabled)
  },
  async getPosts({ boardId = '', authorId = '', page = 1, limit = 20, includePending = false, q = '', sort = 'hot' } = {}) {
    const query = new URLSearchParams()
    if (boardId) query.set('board_id', String(boardId))
    if (authorId) query.set('author_id', String(authorId))
    if (q) query.set('q', q)
    if (sort) query.set('sort', sort)
    query.set('page', String(page))
    query.set('limit', String(limit))
    const payload = await viewerRequest('/forum-api', `/api/v1/posts?${query}`)
    let posts = (payload.posts || []).map(mapPostListItem)
    if (!includePending) posts = posts.filter((item) => item.status === 'published')
    return {
      posts,
      total: payload.total ?? posts.length,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },
  async getPost(id) {
    const post = await viewerRequest('/forum-api', `/api/v1/posts/${id}`)
    const commentsPayload = await request('/forum-api', `/api/v1/posts/${id}/comments?limit=100`)
    return mapPostDetail(post, commentsPayload.comments || [])
  },
  async createPost(payload) {
    const body = {
      title: payload.title,
      content: payload.content,
      board_id: Number(payload.boardId),
    }
    if (payload.attachmentIds?.length) body.attachment_ids = payload.attachmentIds.map(Number)
    const post = await request('/forum-api', '/api/v1/posts', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify(body),
    })
    return mapPostListItem(post)
  },
  async updatePost(id, payload) {
    const post = await request('/forum-api', `/api/v1/posts/${id}`, {
      method: 'PUT',
      headers: authHeaders(),
      body: JSON.stringify({ title: payload.title, content: payload.content }),
    })
    return mapPostListItem(post)
  },
  async deletePost(id) {
    await request('/forum-api', `/api/v1/posts/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    })
  },
  async createComment(id, payload) {
    const body = { content: payload.content }
    if (payload.parentId) body.parent_id = Number(payload.parentId)
    await request('/forum-api', `/api/v1/posts/${id}/comments`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify(body),
    })
  },
  async likeComment(id) {
    const resp = await request('/forum-api', `/api/v1/comments/${id}/like`, {
      method: 'POST',
      headers: authHeaders(),
    })
    return { likeCount: resp.like_count ?? 0, liked: resp.liked }
  },
  async dislikeComment(id) {
    const resp = await request('/forum-api', `/api/v1/comments/${id}/dislike`, {
      method: 'POST',
      headers: authHeaders(),
    })
    return { dislikeCount: resp.dislike_count ?? 0, disliked: resp.disliked }
  },
  async checkSensitiveWords(text) {
    const result = await request('/forum-api', '/api/v1/moderation/check', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ text }),
    })
    return result
  },
  async likePost(id) {
    const resp = await request('/forum-api', `/api/v1/posts/${id}/like`, {
      method: 'POST',
      headers: authHeaders(),
    })
    return {
      likeCount: resp.like_count ?? 0,
      dislikeCount: resp.dislike_count ?? 0,
      liked: !!resp.liked,
      disliked: !!resp.disliked,
    }
  },
  async dislikePost(id) {
    const resp = await request('/forum-api', `/api/v1/posts/${id}/dislike`, {
      method: 'POST',
      headers: authHeaders(),
    })
    return {
      likeCount: resp.like_count ?? 0,
      dislikeCount: resp.dislike_count ?? 0,
      liked: !!resp.liked,
      disliked: !!resp.disliked,
    }
  },
  async collectPost(id) {
    return request('/forum-api', `/api/v1/posts/${id}/collect`, {
      method: 'POST',
      headers: authHeaders(),
    })
  },
  async getMyCollections({ page = 1, limit = 20 } = {}) {
    const query = new URLSearchParams({ page: String(page), limit: String(limit) })
    const payload = await request('/forum-api', `/api/v1/me/collections?${query}`, {
      headers: authHeaders(),
    })
    const posts = (payload.posts || []).map(mapPostListItem)
    return {
      posts,
      total: payload.total ?? posts.length,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },
  async getCommunityStats() {
    const payload = await request('/forum-api', '/api/v1/stats/community')
    return payload.data || payload
  },
  async listNotifications(page = 1, limit = 20) {
    const payload = await request('/forum-api', `/api/v1/notifications?page=${page}&limit=${limit}`, {
      headers: authHeaders(),
    })
    const items = (payload.notifications || []).map(mapNotification)
    return {
      items,
      total: payload.total ?? items.length,
    }
  },
  async getUnreadNotificationCount() {
    const payload = await request('/forum-api', '/api/v1/notifications/unread-count', {
      headers: authHeaders(),
    })
    return Number(payload.unread_count ?? payload.unreadCount ?? 0)
  },
  async markNotificationRead(id) {
    await request('/forum-api', `/api/v1/notifications/${id}/read`, {
      method: 'PUT',
      headers: authHeaders(),
    })
  },
  async reportPost(id, reason) {
    return request('/forum-api', `/api/v1/posts/${id}/report`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ reason }),
    })
  },
  async uploadAttachment({ type, file, linkUrl }) {
    const form = new FormData()
    form.append('type', type)
    if (type === 'link') form.append('link_url', linkUrl)
    else if (file) form.append('file', file)
    const payload = await apiUpload('/forum-api', '/api/v1/attachments/upload', form)
    return String(payload.id)
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
      forumApi.getBoards(true),
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
    const boards = await forumApi.getBoards(true)
    return mapConfigFromBackend(payload.configs || {}, boards)
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
    const boards = await forumApi.getBoards(true)
    for (const [boardId, enabled] of Object.entries(config.boardSwitches || {})) {
      const board = boards.find((item) => item.id === boardId)
      if (board) updates.push([`board_${board.slug}_enabled`, enabled ? 'true' : 'false'])
    }
    for (const [key, value] of updates) {
      await request('/admin-api', `/api/v1/admin/config/${key}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ value: String(value) }),
      })
    }
    return adminApi.getConfig()
  },
  async getPendingAudit() {
    const payload = await request('/admin-api', '/api/v1/admin/audit/pending', {
      headers: authHeaders(),
    })
    return (payload.posts || []).map((post) => ({
      id: String(post.id),
      postId: String(post.id),
      title: post.title,
      authorName: post.author_name,
      boardName: post.board_name,
      status: post.status,
      reason: post.matched_words?.length ? `命中敏感词：${post.matched_words.join('、')}` : '',
    }))
  },
  async approveAudit(id) {
    await request('/admin-api', `/api/v1/admin/audit/${id}/approve`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({}),
    })
  },
  async rejectAudit(id, reason = '不符合社区规范') {
    await request('/admin-api', `/api/v1/admin/audit/${id}/reject`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ reason }),
    })
  },
  async batchDeleteAudit(postIds, reason = '') {
    await request('/admin-api', '/api/v1/admin/audit/batch-delete', {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({
        post_ids: postIds.map((id) => Number(id)),
        reason,
      }),
    })
  },
  async getReports(status = 'pending', page = 1, limit = 50) {
    const payload = await request(
      '/admin-api',
      `/api/v1/admin/reports?page=${page}&limit=${limit}&status=${encodeURIComponent(status)}`,
      { headers: authHeaders() },
    )
    return (payload.reports || []).map((row) => ({
      id: String(row.id),
      postId: String(row.post_id),
      postTitle: row.post_title || '（帖子已删除）',
      reporterId: String(row.reporter_id),
      reporterName: row.reporter_name || row.reporter_id,
      reason: row.reason,
      status: row.status,
      adminNote: row.admin_note || '',
      createdAtIso: row.created_at,
      createdAt: formatDisplayDate(row.created_at),
    }))
  },
  async resolveReport(id, { action, delete_post = false, admin_note = '' }) {
    await request('/admin-api', `/api/v1/admin/reports/${id}/resolve`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({
        action,
        delete_post,
        admin_note,
      }),
    })
  },
  async listPosts(page = 1, limit = 20) {
    const payload = await request(
      '/admin-api',
      `/api/v1/admin/posts?page=${page}&limit=${limit}`,
      { headers: authHeaders() },
    )
    return {
      posts: (payload.posts || []).map(mapPostListItem),
      total: payload.total ?? 0,
      page: payload.page ?? page,
      limit: payload.limit ?? limit,
    }
  },
  async deletePost(id) {
    await request('/admin-api', `/api/v1/admin/posts/${id}/delete`, {
      method: 'POST',
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
}
