import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const { userApi } = vi.hoisted(() => ({
  userApi: {
    me: vi.fn(),
    logout: vi.fn(),
  },
}))

vi.mock('../../src/api', () => ({
  userApi,
}))

describe('session store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    userApi.me.mockReset()
    userApi.logout.mockReset()
    localStorage.clear()
    document.cookie = 'ai-forum-token=;path=/;max-age=0'
  })

  it('validates stored tokens and restores the SSO cookie', async () => {
    localStorage.setItem('ai-forum-token', 'valid-token')
    localStorage.setItem(
      'ai-forum-user',
      JSON.stringify({ id: '7', username: 'tester', role: 'student' }),
    )
    userApi.me.mockResolvedValueOnce({ id: '7', username: 'tester', role: 'student' })

    const { useSessionStore } = await import('../../src/stores/session')
    const session = useSessionStore()

    await expect(session.ensureValidSession()).resolves.toBe(true)
    expect(userApi.me).toHaveBeenCalledTimes(1)
    expect(document.cookie).toContain('ai-forum-token=valid-token')
    expect(localStorage.getItem('ai-forum-token')).toBe('valid-token')
  })

  it('clears stale stored tokens instead of redirecting to SSO', async () => {
    localStorage.setItem('ai-forum-token', 'stale-token')
    localStorage.setItem(
      'ai-forum-user',
      JSON.stringify({ id: '7', username: 'tester', role: 'student' }),
    )
    userApi.me.mockRejectedValueOnce(new Error('unauthorized'))

    const { useSessionStore } = await import('../../src/stores/session')
    const session = useSessionStore()

    await expect(session.ensureValidSession()).resolves.toBe(false)
    expect(localStorage.getItem('ai-forum-token')).toBeNull()
    expect(session.token).toBe('')
    expect(session.currentUser).toBeNull()
  })
})
