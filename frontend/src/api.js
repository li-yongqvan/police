const defaultHeaders = (token) => ({
  'Content-Type': 'application/json',
  ...(token ? { Authorization: `Bearer ${token}` } : {}),
})

function getToken(explicit) {
  if (explicit) return explicit
  return localStorage.getItem('ai-forum-token') || ''
}

async function request(base, path, options = {}) {
  const response = await fetch(`${base}${path}`, options)
  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(payload.error || payload.message || '请求失败')
  }
  return payload
}

function mapBoard(board) {
  return {
    id: String(board.id),
    name: board.name,
    slug: board.slug,
    description: board.description,
    enabled: board.enabled,
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
    authorId: String(post.author_id),
    authorName: post.author_name,
    status: post.status,
    isFeatured: post.is_featured,
    isPinned: post.is_pinned,
    likeCount: post.like_count ?? 0,
    commentCount: post.comment_count ?? 0,
    createdAt: post.created_at,
    tags: [],
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
        title: item.filename,
        url: item.file_path,
      })),
      tags: [],
    },
    comments: comments.map((item) => ({
      id: String(item.id),
      postId: String(item.post_id),
      authorId: String(item.author_id),
      authorName: item.author_name || '',
      content: item.content,
      createdAt: item.created_at,
    })),
  }
}

function mapUser(user) {
  return {
    id: String(user.id),
    name: user.name || user.nickname || user.username,
    avatar: user.avatar || '',
    role: user.role || 'student',
    department: user.department || '',
    bio: user.bio || '',
    status: user.status || 'active',
  }
}

function mapConfigFromBackend(configs = {}) {
  const boardSwitches = {}
  Object.entries(configs).forEach(([key, value]) => {
    const match = key.match(/^board_(.+)_enabled$/)
    if (match) {
      boardSwitches[match[1]] = value === 'true'
    }
  })
  return {
    postingEnabled: configs.post_requires_level !== '99',
    moderationMode: configs.sensitive_word_action === 'pending_review' ? 'manual' : 'auto',
    boardSwitches,
    _raw: configs,
  }
}

export const userApi = {
  async demoLogin(role) {
    const payload = await request('/user-api', '/api/v1/demo-login', {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify({ role }),
    })
    return {
      token: payload.token || payload.access_token,
      user: mapUser(payload.user),
    }
  },
  async me(token) {
    const payload = await request('/user-api', '/api/v1/users/me', {
      headers: defaultHeaders(getToken(token)),
    })
    return mapUser(payload)
  },
  async listUsers(token) {
    const payload = await request('/admin-api', '/api/v1/admin/users?limit=100', {
      headers: defaultHeaders(getToken(token)),
    })
    return (payload.users || []).map((user) => ({
      id: String(user.id ?? user.user_id),
      name: user.nickname || user.username,
      avatar: user.avatar || '',
      role: user.role || 'student',
      department: '',
      bio: user.bio || '',
      status: user.status || 'active',
    }))
  },
}

export const forumApi = {
  async getBoards(includeDisabled = false) {
    const boards = await request('/forum-api', '/api/v1/boards')
    const list = Array.isArray(boards) ? boards.map(mapBoard) : []
    if (!includeDisabled) {
      return list.filter((item) => item.enabled)
    }
    return list
  },
  async getPosts(boardId = '', includePending = false) {
    const query = new URLSearchParams()
    if (boardId) query.set('board_id', boardId)
    query.set('limit', '100')
    const payload = await request('/forum-api', `/api/v1/posts?${query}`)
    let posts = (payload.posts || []).map(mapPostListItem)
    if (!includePending) {
      posts = posts.filter((item) => item.status === 'published')
    }
    return posts
  },
  async getPost(id) {
    const post = await request('/forum-api', `/api/v1/posts/${id}`)
    const commentsPayload = await request('/forum-api', `/api/v1/posts/${id}/comments?limit=100`)
    return mapPostDetail(post, commentsPayload.comments || [])
  },
  async createPost(payload, token) {
    const post = await request('/forum-api', '/api/v1/posts', {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
      body: JSON.stringify({
        title: payload.title,
        content: payload.content,
        board_id: Number(payload.boardId),
      }),
    })
    return mapPostListItem(post)
  },
  async createComment(id, payload, token) {
    await request('/forum-api', `/api/v1/posts/${id}/comments`, {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
      body: JSON.stringify({ content: payload.content }),
    })
  },
  async likePost(id, token) {
    await request('/forum-api', `/api/v1/posts/${id}/like`, {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
    })
  },
  async toggleFeature() {
    throw new Error('精华帖操作请在中台完成')
  },
}

export const adminApi = {
  async getOverview(token) {
    const payload = await request('/admin-api', '/api/v1/admin/stats/overview', {
      headers: defaultHeaders(getToken(token)),
    })
    const data = payload.data || payload
    return {
      userCount: data.total_users ?? 0,
      postCount: data.total_posts ?? 0,
      commentCount: data.total_comments ?? 0,
      pendingAuditCount: 0,
      boardActivity: [],
    }
  },
  async getConfig(token) {
    const payload = await request('/admin-api', '/api/v1/admin/config', {
      headers: defaultHeaders(getToken(token)),
    })
    const configs = payload.configs || {}
    const mapped = mapConfigFromBackend(configs)
    const boards = await forumApi.getBoards(true)
    const boardSwitches = {}
    boards.forEach((board) => {
      boardSwitches[board.id] = configs[`board_${board.slug}_enabled`] !== 'false'
    })
    mapped.boardSwitches = boardSwitches
    return mapped
  },
  async updateConfig(config, token) {
    const updates = []
    if (config.moderationMode) {
      updates.push(['sensitive_word_action', config.moderationMode === 'manual' ? 'pending_review' : 'reject'])
    }
    for (const [boardId, enabled] of Object.entries(config.boardSwitches || {})) {
      const boards = await forumApi.getBoards(true)
      const board = boards.find((item) => item.id === boardId)
      if (board) {
        updates.push([`board_${board.slug}_enabled`, enabled ? 'true' : 'false'])
      }
    }
    for (const [key, value] of updates) {
      await request('/admin-api', `/api/v1/admin/config/${key}`, {
        method: 'PUT',
        headers: defaultHeaders(getToken(token)),
        body: JSON.stringify({ value: String(value) }),
      })
    }
    return this.getConfig(token)
  },
  async getPendingAudit(token) {
    const payload = await request('/admin-api', '/api/v1/admin/audit/pending', {
      headers: defaultHeaders(getToken(token)),
    })
    return (payload.posts || []).map((post) => ({
      id: String(post.id),
      postId: String(post.id),
      title: post.title,
      authorName: post.author_name,
      boardName: post.board_name,
      status: post.status,
    }))
  },
  async approveAudit(id, token) {
    await request('/admin-api', `/api/v1/admin/audit/${id}/approve`, {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
      body: JSON.stringify({}),
    })
  },
  async rejectAudit(id, token) {
    await request('/admin-api', `/api/v1/admin/audit/${id}/reject`, {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
      body: JSON.stringify({ reason: '演示驳回' }),
    })
  },
  async deletePost(id, token) {
    await request('/admin-api', `/api/v1/admin/posts/${id}/delete`, {
      method: 'POST',
      headers: defaultHeaders(getToken(token)),
    })
  },
  async setUserStatus(id, status, token) {
    if (status === 'banned') {
      await request('/admin-api', `/api/v1/admin/users/${id}/ban`, {
        method: 'POST',
        headers: defaultHeaders(getToken(token)),
        body: JSON.stringify({ reason: '演示封禁' }),
      })
    }
  },
}
