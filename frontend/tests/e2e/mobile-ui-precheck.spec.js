import { expect, test } from '@playwright/test'

const adminUser = {
  id: 1,
  username: 'demo_admin',
  name: '演示管理员',
  role: 'platform_admin',
  level: 5,
  status: 'active',
}

const sampleUsers = Array.from({ length: 6 }, (_, index) => ({
  id: index + 1,
  username: `qq_user_${index + 1}`,
  name: `用户 ${index + 1}`,
  role: index === 0 ? 'admin' : 'student',
  level: index % 5,
  status: index === 4 ? 'banned' : 'active',
}))

const inviteCodes = [
  { code: 'DEMO2026', status: 'active' },
  { code: 'USED2026', status: 'active', used_by: 'qq_user_1' },
]

async function mockAdminSession(page) {
  await page.addInitScript((user) => {
    localStorage.setItem('ai-forum-token', 'local-precheck-token')
    localStorage.setItem('ai-forum-user', JSON.stringify(user))
  }, adminUser)

  await page.route('**/user-api/api/v1/users/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(adminUser),
    })
  })

  await page.route('**/user-api/api/v1/admin/users**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ users: sampleUsers, total: 46 }),
    })
  })

  await page.route('**/admin-api/api/v1/admin/invite-codes**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ codes: inviteCodes }),
    })
  })

  await page.route('**/admin-api/**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ ok: true }),
    })
  })
}

test.describe('mobile UI precheck', () => {
  test('login inputs are visible on a short landscape viewport', async ({ page }) => {
    await page.setViewportSize({ width: 667, height: 375 })
    await page.goto('/')

    const username = await page.locator('#username').boundingBox()
    const password = await page.locator('#password').boundingBox()
    const viewport = page.viewportSize()

    expect(username).not.toBeNull()
    expect(password).not.toBeNull()
    expect(username.y).toBeGreaterThanOrEqual(0)
    expect(password.y + password.height).toBeLessThanOrEqual(viewport.height)
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
  })

  for (const route of ['/admin', '/admin/users', '/admin/invites']) {
    test(`admin drawer and topbar pass mobile checks on ${route}`, async ({ page }) => {
      await page.setViewportSize({ width: 390, height: 844 })
      await mockAdminSession(page)
      await page.goto(route)

      const menu = await page.locator('.gx-admin-topbar__menu').boundingBox()
      const logout = await page.locator('.gx-admin-topbar__logout').boundingBox()

      expect(menu.width).toBeGreaterThanOrEqual(44)
      expect(menu.height).toBeGreaterThanOrEqual(44)
      expect(logout.width).toBeGreaterThanOrEqual(44)
      expect(logout.height).toBeGreaterThanOrEqual(44)

      await page.locator('.gx-admin-topbar__menu').click()
      await expect(page.locator('.gx-admin-sidebar')).toHaveClass(/is-open/)
      await page.waitForTimeout(180)

      const drawer = await page.locator('.gx-admin-sidebar').boundingBox()
      const drawerLinks = await page.locator('.gx-admin-sidebar .gx-sidebar-nav__link').evaluateAll((links) =>
        links.map((link) => {
          const rect = link.getBoundingClientRect()
          return { left: rect.left, height: rect.height }
        }),
      )

      expect(drawer.x).toBeGreaterThanOrEqual(0)
      expect(drawer.x).toBeLessThan(1)
      expect(await page.evaluate(() => document.body.classList.contains('mw-drawer-open'))).toBe(true)
      expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)

      for (const link of drawerLinks) {
        expect(link.left).toBeGreaterThanOrEqual(0)
        expect(link.height).toBeGreaterThanOrEqual(44)
      }
    })
  }
})
