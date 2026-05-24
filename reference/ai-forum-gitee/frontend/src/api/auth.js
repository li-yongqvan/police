import api from './index'

export function register(data) {
  return api.post('/register', data)
}

export function login(data) {
  return api.post('/login', data)
}

export function refreshToken(refreshToken) {
  return api.post('/auth/refresh', { refresh_token: refreshToken })
}

export function getProfile(userId) {
  return api.get(`/users/${userId}`)
}

export function updateProfile(userId, data) {
  return api.put(`/users/${userId}`, data)
}

export function uploadAvatar(userId, formData) {
  return api.post(`/users/${userId}/avatar`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
