import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1/admin',
  timeout: 15000,
})

// Request interceptor: attach JWT token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: handle 401 by redirecting to login
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('access_token')
      window.location.href = '/admin/login'
    }
    return Promise.reject(error)
  }
)

// Audit API
export function getPendingAudits(page, limit) {
  return api.get(`/audit/pending`, { params: { page, limit } })
}
export function approvePost(postId) {
  return api.post(`/audit/${postId}/approve`)
}
export function rejectPost(postId, reason) {
  return api.post(`/audit/${postId}/reject`, { reason })
}
export function batchDeletePosts(postIds, reason) {
  return api.post(`/audit/batch-delete`, { post_ids: postIds, reason })
}

// Post management API
export function listPosts(page, limit) {
  return api.get(`/posts`, { params: { page, limit } })
}
export function deletePost(postId) {
  return api.post(`/posts/${postId}/delete`)
}
export function setPostFeatured(postId, featured) {
  return api.post(`/posts/${postId}/featured`, { featured })
}
export function setPostPinned(postId, pinned) {
  return api.post(`/posts/${postId}/pinned`, { pinned })
}

// User management API
export function listUsers(page, limit, status) {
  return api.get('/users', { params: { page, limit, status } })
}
export function banUser(userId, reason) {
  return api.post(`/users/${userId}/ban`, { reason })
}
export function updateUserLevel(userId, level) {
  return api.put(`/users/${userId}/level`, { level })
}
export function getUserLogs(userId, page, limit) {
  return api.get(`/users/${userId}/logs`, { params: { page, limit } })
}

// Invite code management API
export function listInviteCodes(page, limit) {
  return api.get('/invite-codes', { params: { page, limit } })
}
export function getInviteCodeStatus(code) {
  return api.get(`/invite-codes/${code}/status`)
}
export function voidInviteCode(code) {
  return api.put(`/invite-codes/${code}/void`)
}

// System config API
export function getConfig() {
  return api.get('/config')
}
export function updateConfig(key, value) {
  return api.put(`/config/${key}`, { value })
}

// Board management API
export function listBoards() {
  return api.get('/boards')
}
export function createBoard(data) {
  return api.post('/boards', data)
}
export function updateBoard(id, data) {
  return api.put(`/boards/${id}`, data)
}
export function deleteBoard(id) {
  return api.delete(`/boards/${id}`)
}

// Role management API
export function listRoles() {
  return api.get('/roles')
}
export function assignRole(userId, roleId) {
  return api.post(`/users/${userId}/roles`, { role_id: roleId })
}

// Sensitive words API
export function listSensitiveWords() {
  return api.get('/sensitive-words')
}
export function addSensitiveWord(word, category) {
  return api.post('/sensitive-words', { word, category })
}
export function deleteSensitiveWord(id) {
  return api.delete(`/sensitive-words/${id}`)
}

// Stats API
export function getStatsOverview() {
  return api.get('/stats/overview')
}
export function getDailyStats(days) {
  return api.get('/stats/daily', { params: { days } })
}

export default api
