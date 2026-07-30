import { errorMessageFrom, formatApiError } from './errors'

function apiErrorFromPayload(payload, fallback) {
  const msg = payload?.error ?? payload?.message ?? payload?.detail
  return formatApiError({ message: errorMessageFrom(msg, fallback) }, fallback)
}

const REFRESH_KEY = 'ai-forum-refresh-token'

export const DEFAULT_TIMEOUT_MS = 15000
export const LOGIN_TIMEOUT_MS = 30000

/**
 * Write token to document.cookie so Discourse SSO (cross-port 80->8080)
 * can read it via c.Request.Cookie("ai-forum-token").
 */
export function setTokenCookie(token) {
  if (token) {
    document.cookie = 'ai-forum-token=' + encodeURIComponent(token) + ';path=/;SameSite=Lax'
  } else {
    document.cookie = 'ai-forum-token=;path=/;SameSite=Lax;max-age=0'
  }
}

let unauthorizedHandler = () => {
  localStorage.removeItem('ai-forum-token')
  localStorage.removeItem('ai-forum-user')
  localStorage.removeItem(REFRESH_KEY)
  setTokenCookie('')
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

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function isRetryableStatus(status) {
  return status === 502 || status === 503 || status === 504
}

function isNetworkError(error) {
  return error?.name === 'AbortError' || error?.name === 'TypeError'
}

export function networkErrorMessage(error) {
  if (error?.name === 'AbortError') {
    return '网络响应超时，请切换 Wi-Fi/4G 后重试'
  }
  if (error?.name === 'TypeError') {
    return '网络不稳定，无法连接服务器，请稍后重试'
  }
  return null
}

/**
 * fetch with timeout and limited retries (GET or explicit retryOnPost).
 */
export async function fetchWithRetry(url, options = {}) {
  const timeoutMs = options.timeoutMs ?? DEFAULT_TIMEOUT_MS
  const retries = options.retries ?? 0
  const method = (options.method || 'GET').toUpperCase()
  const allowRetry = method === 'GET' || options.retryOnPost === true

  const {
    timeoutMs: _t,
    retries: _r,
    retryOnPost: _p,
    _retried: _rt,
    ...fetchOpts
  } = options

  let lastError
  const attempts = allowRetry ? retries + 1 : 1

  for (let attempt = 0; attempt < attempts; attempt += 1) {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), timeoutMs)
    try {
      const response = await fetch(url, { ...fetchOpts, signal: controller.signal })
      clearTimeout(timer)
      if (allowRetry && attempt < attempts - 1 && isRetryableStatus(response.status)) {
        await sleep(400 * (attempt + 1))
        continue
      }
      return response
    } catch (error) {
      clearTimeout(timer)
      lastError = error
      if (allowRetry && attempt < attempts - 1 && isNetworkError(error)) {
        await sleep(400 * (attempt + 1))
        continue
      }
      throw error
    }
  }
  throw lastError
}

async function tryRefreshAccessToken() {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return false

  try {
    const response = await fetchWithRetry('/user-api/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
      timeoutMs: LOGIN_TIMEOUT_MS,
      retries: 1,
    })
    const payload = await response.json().catch(() => ({}))
    if (!response.ok || !payload.access_token) return false
    localStorage.setItem('ai-forum-token', payload.access_token)
    setTokenCookie(payload.access_token)
    return true
  } catch {
    return false
  }
}

export async function apiRequest(base, path, options = {}) {
  const url = `${base}${path}`
  let response
  try {
    response = await fetchWithRetry(url, options)
  } catch (error) {
    throw new Error(networkErrorMessage(error) || formatApiError(error, '请求失败'))
  }

  const payload = await response.json().catch(() => ({}))

  if (response.status === 401 && !options._retried) {
    const isAuthEndpoint =
      path.includes('/api/v1/login') ||
      path.includes('/api/v1/register') ||
      path.includes('/api/v1/demo-login')
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
    if (!isAuthEndpoint) {
      unauthorizedHandler()
    }
    throw new Error(apiErrorFromPayload(payload, '登录已过期，请重新登录'))
  }

  if (!response.ok) {
    throw new Error(apiErrorFromPayload(payload, '请求失败'))
  }
  return payload
}

export async function apiUpload(base, path, formData, options = {}) {
  const headers = {}
  const token = localStorage.getItem('ai-forum-token')
  if (token) headers.Authorization = `Bearer ${token}`

  let response
  try {
    response = await fetchWithRetry(`${base}${path}`, {
      method: 'POST',
      headers,
      body: formData,
      timeoutMs: options.timeoutMs ?? 60000,
      retries: options.retries ?? 1,
      retryOnPost: true,
    })
  } catch (error) {
    throw new Error(networkErrorMessage(error) || formatApiError(error, '上传失败'))
  }

  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new Error(apiErrorFromPayload(payload, '上传失败'))
  }
  return payload
}