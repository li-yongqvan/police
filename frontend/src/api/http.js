import { formatApiError } from './errors'

const REFRESH_KEY = 'ai-forum-refresh-token'

let unauthorizedHandler = () => {
  localStorage.removeItem('ai-forum-token')
  localStorage.removeItem('ai-forum-user')
  localStorage.removeItem(REFRESH_KEY)
  if (!window.location.pathname.startsWith('/')) {
    window.location.href = '/'
  } else if (window.location.pathname !== '/') {
    window.location.href = '/'
  }
}

export function setUnauthorizedHandler(handler) {
  unauthorizedHandler = handler
}

export function getRefreshToken() {
  return localStorage.getItem(REFRESH_KEY) || ''
}

export function setRefreshToken(token) {
  if (token) {
    localStorage.setItem(REFRESH_KEY, token)
  } else {
    localStorage.removeItem(REFRESH_KEY)
  }
}

async function tryRefreshAccessToken() {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return false

  try {
    const response = await fetch('/user-api/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
    const payload = await response.json().catch(() => ({}))
    if (!response.ok || !payload.access_token) return false
    localStorage.setItem('ai-forum-token', payload.access_token)
    return true
  } catch {
    return false
  }
}

export async function apiRequest(base, path, options = {}) {
  const response = await fetch(`${base}${path}`, options)
  const payload = await response.json().catch(() => ({}))

  if (response.status === 401 && !options._retried) {
    const refreshed = await tryRefreshAccessToken()
    if (refreshed) {
      const headers = new Headers(options.headers || {})
      const token = localStorage.getItem('ai-forum-token')
      if (token) headers.set('Authorization', `Bearer ${token}`)
      return apiRequest(base, path, {
        ...options,
        headers,
        _retried: true,
      })
    }
    unauthorizedHandler()
    throw new Error(formatApiError({ message: payload.error || payload.message }, '登录已过期，请重新登录'))
  }

  if (!response.ok) {
    throw new Error(formatApiError({ message: payload.error || payload.message }))
  }
  return payload
}

export async function apiUpload(base, path, formData) {
  const headers = {}
  const token = localStorage.getItem('ai-forum-token')
  if (token) headers.Authorization = `Bearer ${token}`

  const response = await fetch(`${base}${path}`, {
    method: 'POST',
    headers,
    body: formData,
  })
  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(formatApiError({ message: payload.error || payload.message }, '上传失败'))
  }
  return payload
}
