import { defineStore } from 'pinia'
import { register, login as loginApi, getProfile, updateProfile, uploadAvatar } from '../api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    accessToken: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
  }),
  getters: {
    isLoggedIn: (state) => !!state.accessToken,
    userLevel: (state) => state.user?.level || 0,
  },
  actions: {
    async login(username, password) {
      const { data } = await loginApi({ username, password })
      this.accessToken = data.access_token
      this.refreshToken = data.refresh_token
      this.user = data.user
      localStorage.setItem('access_token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)
      localStorage.setItem('user', JSON.stringify(data.user))
      return data
    },
    async register(username, password, invitationCode) {
      const { data } = await register({ username, password, invitation_code: invitationCode })
      return data
    },
    logout() {
      this.user = null
      this.accessToken = ''
      this.refreshToken = ''
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
    },
    async loadUser() {
      if (this.user?.id) {
        try {
          const { data } = await getProfile(this.user.id)
          this.user = data
          localStorage.setItem('user', JSON.stringify(data))
        } catch (e) {
          // ignore
        }
      }
    },
    async updateProfile(data) {
      if (!this.user) return
      const { data: updated } = await updateProfile(this.user.id, data)
      this.user = updated
      localStorage.setItem('user', JSON.stringify(updated))
    },
    async uploadAvatar(file) {
      if (!this.user) return
      const formData = new FormData()
      formData.append('avatar', file)
      const { data } = await uploadAvatar(this.user.id, formData)
      this.user.avatar = data.avatar
      localStorage.setItem('user', JSON.stringify(this.user))
    },
  },
})
