import { expect, test } from '@playwright/test'

test('login page renders on desktop and mobile', async ({ page }) => {
  await page.goto('/')

  await expect(page).toHaveTitle(/AI/)
  await expect(page.locator('#username')).toBeVisible()
  await expect(page.locator('#password')).toBeVisible()
  await expect(page.locator('button[type="submit"]')).toBeVisible()
})
