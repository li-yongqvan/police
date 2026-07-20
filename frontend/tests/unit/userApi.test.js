import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiRequest, apiUpload, setRefreshToken } = vi.hoisted(() => ({
  apiRequest: vi.fn(),
  apiUpload: vi.fn(),
  setRefreshToken: vi.fn(),
}))

vi.mock('../../src/api/http', () => ({
  apiRequest,
  apiUpload,
  LOGIN_TIMEOUT_MS: 30000,
  setRefreshToken,
}))

describe('userApi.register', () => {
  beforeEach(() => {
    apiRequest.mockReset()
    apiUpload.mockReset()
    setRefreshToken.mockReset()
    localStorage.clear()
  })

  it('uppercases invitation codes and maps the returned user payload', async () => {
    apiRequest.mockResolvedValueOnce({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      user: {
        id: 18,
        username: 'alice',
        nickname: 'Alice',
        avatar: '/uploads/avatar.png',
        role: 'student',
        level: 3,
        profile_completed: true,
      },
    })

    const { userApi } = await import('../../src/api')
    const result = await userApi.register({
      username: 'alice',
      password: 'secret',
      invitationCode: 'demo2026',
    })

    expect(apiRequest).toHaveBeenCalledWith(
      '/user-api',
      '/api/v1/register',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          username: 'alice',
          password: 'secret',
          invitation_code: 'DEMO2026',
        }),
      }),
    )
    expect(localStorage.getItem('ai-forum-token')).toBe('access-token')
    expect(setRefreshToken).toHaveBeenCalledWith('refresh-token')
    expect(result).toMatchObject({
      token: 'access-token',
      refreshToken: 'refresh-token',
      user: {
        id: '18',
        username: 'alice',
        name: 'Alice',
        role: 'student',
        level: 3,
      },
    })
  })
})
