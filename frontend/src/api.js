import { apiRequest, LOGIN_TIMEOUT_MS } from './api/http'

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
  async setUserStatus(id, status) {
    if (status === 'banned') {
      await request('/admin-api', `/api/v1/admin/users/${id}/ban`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ reason: '婕旂ず灏佺' }),
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

}
