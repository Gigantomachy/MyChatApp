import { test as setup, expect } from '@playwright/test'

const AUTH_DIR = 'tests/.auth'

setup('login as steven', async ({ page }) => {
  await page.goto('/')
  await page.getByPlaceholder('jdoe').fill('steven')
  await page.getByPlaceholder('••••••••').fill('password123')
  await page.getByRole('button', { name: 'Sign In' }).click()
  await expect(page.locator('.chat-layout')).toBeVisible({ timeout: 10000 })
  await page.context().storageState({ path: `${AUTH_DIR}/steven.json` })
})

setup('login as jessica', async ({ page }) => {
  await page.goto('/')
  await page.getByPlaceholder('jdoe').fill('jessica')
  await page.getByPlaceholder('••••••••').fill('password123')
  await page.getByRole('button', { name: 'Sign In' }).click()
  await expect(page.locator('.chat-layout')).toBeVisible({ timeout: 10000 })
  await page.context().storageState({ path: `${AUTH_DIR}/jessica.json` })
})

setup('login as bob', async ({ page }) => {
  await page.goto('/')
  await page.getByPlaceholder('jdoe').fill('bob')
  await page.getByPlaceholder('••••••••').fill('password123')
  await page.getByRole('button', { name: 'Sign In' }).click()
  await expect(page.locator('.chat-layout')).toBeVisible({ timeout: 10000 })
  await page.context().storageState({ path: `${AUTH_DIR}/bob.json` })
})
