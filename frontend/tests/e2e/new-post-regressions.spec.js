import { expect, test } from '@playwright/test'

const boards = [
  {
    id: 1,
    name: 'Study',
    slug: 'study',
    description: 'Study board',
    enabled: true,
    sort_order: 1,
    post_count: 0,
  },
  {
    id: 2,
    name: 'Training',
    slug: 'training',
    description: 'Training board',
    enabled: true,
    sort_order: 2,
    post_count: 0,
  },
]

async function seedSession(page) {
  await page.addInitScript(() => {
    localStorage.setItem('ai-forum-token', 'test-token')
    localStorage.setItem(
      'ai-forum-user',
      JSON.stringify({
        id: '7',
        username: 'tester',
        name: 'Tester',
        role: 'student',
        level: 3,
      }),
    )
  })
}

async function mockForumApi(page) {
  await page.route('**/user-api/api/v1/users/me', async (route) => {
    await route.fulfill({
      json: {
        id: 7,
        username: 'tester',
        nickname: 'Tester',
        role: 'student',
        level: 3,
      },
    })
  })
  await page.route('**/forum-api/api/v1/boards', async (route) => {
    await route.fulfill({ json: boards })
  })
  await page.route('**/forum-api/api/v1/stats/community', async (route) => {
    await route.fulfill({
      json: {
        total_users: 3,
        total_posts: 0,
        online_users: 1,
        posts_today: 0,
      },
    })
  })
  await page.route('**/forum-api/api/v1/posts**', async (route) => {
    await route.fulfill({ json: { posts: [], total: 0, page: 1, limit: 10 } })
  })
  await page.route('**/forum-api/api/v1/notifications/unread-count', async (route) => {
    await route.fulfill({ json: { unread_count: 0 } })
  })
}

test.beforeEach(async ({ page }) => {
  await seedSession(page)
  await mockForumApi(page)
})

test('post body accepts lowercase k', async ({ page }) => {
  await page.goto('/community/posts/new')

  const content = page.locator('#content')
  await expect(content).toBeVisible()
  await content.pressSequentially('k')

  await expect(content).toHaveValue('k')
})

test('switching boards clears drafted file attachments', async ({ page }) => {
  await page.goto('/community/posts/new')

  await page.locator('[data-attachment-type="image"]').click()
  const fileInput = page.locator('#attach-file')
  await fileInput.setInputFiles({
    name: 'draft.png',
    mimeType: 'image/png',
    buffer: Buffer.from('image-content'),
  })

  await expect(page.locator('.gx-compose-file__name')).toHaveText('draft.png')
  await page.locator('.gx-compose-select-board__select').evaluate((select) => {
    select.value = '2'
    select.dispatchEvent(new Event('change', { bubbles: true }))
  })

  await expect(page.locator('.gx-compose-file__name')).not.toHaveText('draft.png')
  await expect.poll(() => fileInput.evaluate((input) => input.files.length)).toBe(0)
})

test('more rules link opens the rules section', async ({ page }, testInfo) => {
  test.skip(testInfo.project.name.includes('mobile'), 'board aside is desktop-only')

  await page.goto('/community/boards/study')

  await page.locator('a[href="/community/about#rules"]').click()

  await expect(page).toHaveURL(/\/community\/about#rules$/)
  await expect(page.locator('#rules')).toBeVisible()
})
