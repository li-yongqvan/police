const defaultHeaders = (token) => ({
  'Content-Type': 'application/json',
  ...(token ? { Authorization: `Bearer ${token}` } : {}),
})

async function request(base, path, options = {}) {
  const response = await fetch(`${base}${path}`, options)
  const payload = await response.json()
  if (!response.ok || payload.success === false) {
    throw new Error(payload.message || '请求失败')
  }
  return payload.data
}

export const userApi = {
  demoLogin(role) {
    return request('/user-api', '/api/v1/demo-login', {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify({ role }),
    })
  },
  me(token) {
    return request('/user-api', '/api/v1/users/me', {
      headers: defaultHeaders(token),
    })
  },
  listUsers() {
    return request('/user-api', '/api/v1/admin/users', {
      headers: defaultHeaders(),
    })
  },
}

export const forumApi = {
  getBoards(includeDisabled = false) {
    return request('/forum-api', `/api/v1/boards${includeDisabled ? '?includeDisabled=1' : ''}`)
  },
  getPosts(boardId = '', includePending = false) {
    const query = new URLSearchParams()
    if (boardId) query.set('boardId', boardId)
    if (includePending) query.set('includePending', '1')
    return request('/forum-api', `/api/v1/posts${query.toString() ? `?${query}` : ''}`)
  },
  getPost(id, includeHidden = false) {
    return request('/forum-api', `/api/v1/posts/${id}${includeHidden ? '?includeHidden=1' : ''}`)
  },
  createPost(payload) {
    return request('/forum-api', '/api/v1/posts', {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify(payload),
    })
  },
  createComment(id, payload) {
    return request('/forum-api', `/api/v1/posts/${id}/comments`, {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify(payload),
    })
  },
  likePost(id) {
    return request('/forum-api', `/api/v1/posts/${id}/like`, {
      method: 'POST',
      headers: defaultHeaders(),
    })
  },
  toggleFeature(id) {
    return request('/forum-api', `/api/v1/posts/${id}/feature`, {
      method: 'POST',
      headers: defaultHeaders(),
    })
  },
}

export const adminApi = {
  getOverview() {
    return request('/admin-api', '/api/v1/admin/stats/overview')
  },
  getConfig() {
    return request('/admin-api', '/api/v1/admin/config')
  },
  updateConfig(payload) {
    return request('/admin-api', '/api/v1/admin/config', {
      method: 'PUT',
      headers: defaultHeaders(),
      body: JSON.stringify(payload),
    })
  },
  getPendingAudit() {
    return request('/admin-api', '/api/v1/admin/audit/pending')
  },
  approveAudit(id, reviewerId) {
    return request('/admin-api', `/api/v1/admin/audit/${id}/approve`, {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify({ reviewerId }),
    })
  },
  rejectAudit(id, reviewerId) {
    return request('/admin-api', `/api/v1/admin/audit/${id}/reject`, {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify({ reviewerId }),
    })
  },
  deletePost(id) {
    return request('/admin-api', `/api/v1/admin/posts/${id}/delete`, {
      method: 'POST',
      headers: defaultHeaders(),
    })
  },
  setUserStatus(id, status) {
    return request('/admin-api', `/api/v1/admin/users/${id}/ban`, {
      method: 'POST',
      headers: defaultHeaders(),
      body: JSON.stringify({ status }),
    })
  },
}
