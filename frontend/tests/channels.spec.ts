import { test, expect, type Browser } from '@playwright/test'

test('Steven creates a public channel, Jessica joins it', async ({ browser }: { browser: Browser }) => {
  const stevenCtx = await browser.newContext({ storageState: 'tests/.auth/steven.json' })
  const jessicaCtx = await browser.newContext({ storageState: 'tests/.auth/jessica.json' })

  const stevenPage = await stevenCtx.newPage()
  const jessicaPage = await jessicaCtx.newPage()

  await stevenPage.goto('/')
  await jessicaPage.goto('/')

  // Steven creates a public channel
  await stevenPage.getByRole('button', { name: 'Create channel' }).click()
  await stevenPage.getByPlaceholder('new-channel').fill('test-channel')
  await stevenPage.getByRole('button', { name: 'Create Channel' }).click()

  // Steven sees the channel in his sidebar
  await expect(stevenPage.locator('.sidebar-channel-name', { hasText: '#test-channel' })).toBeVisible({ timeout: 5000 })

  // Jessica opens search and finds the channel
  await jessicaPage.getByRole('button', { name: 'Search' }).click()
  await jessicaPage.getByPlaceholder('Search channels and people...').fill('test')
  await jessicaPage.getByRole('tab', { name: 'Channels' }).click()

  // Jessica sees the channel with a Join button
  await expect(jessicaPage.getByText('#test-channel')).toBeVisible({ timeout: 5000 })
  await jessicaPage.getByRole('button', { name: 'Join' }).click()

  // Jessica sees the channel in her sidebar
  await expect(jessicaPage.locator('.sidebar-channel-name', { hasText: '#test-channel' })).toBeVisible({ timeout: 5000 })

  // Steven sends a message in the channel
  await stevenPage.locator('.sidebar-channel-name', { hasText: '#test-channel' }).click()
  const stevenInput = stevenPage.locator('.chat-input')
  await stevenInput.fill('Welcome to the channel!')
  await stevenInput.press('Enter')

  // Jessica sees the message in real-time
  await expect(jessicaPage.getByText('Welcome to the channel!')).toBeVisible({ timeout: 10000 })

  await stevenCtx.close()
  await jessicaCtx.close()
})
