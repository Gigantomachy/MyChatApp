import { test, expect, type Browser } from '@playwright/test'

test('Steven starts DM with Jessica, real-time messages work', async ({ browser }: { browser: Browser }) => {
  const stevenCtx = await browser.newContext({ storageState: 'tests/.auth/steven.json' })
  const jessicaCtx = await browser.newContext({ storageState: 'tests/.auth/jessica.json' })

  const stevenPage = await stevenCtx.newPage()
  const jessicaPage = await jessicaCtx.newPage()

  await stevenPage.goto('/')
  await jessicaPage.goto('/')

  // Steven starts a DM with Jessica
  await stevenPage.getByRole('button', { name: 'Start new chat' }).click()
  await stevenPage.getByPlaceholder('Search friends...').fill('jessica')
  await stevenPage.getByText('DM').click()

  // Steven sees the DM in his sidebar and the chat area opens
  await expect(stevenPage.locator('.sidebar-channel-name', { hasText: 'Jessica' })).toBeVisible({ timeout: 5000 })

  // Jessica sees the DM appear in her sidebar (via WebSocket channel.new event)
  await expect(jessicaPage.locator('.sidebar-channel-name', { hasText: 'Steven' })).toBeVisible({ timeout: 10000 })

  // Jessica clicks on the DM to open it
  await jessicaPage.locator('.sidebar-channel-name', { hasText: 'Steven' }).click()

  // Steven sends a message
  const stevenInput = stevenPage.locator('.chat-input')
  await stevenInput.fill('Hello Jessica!')
  await stevenInput.press('Enter')

  // Steven sees his own message
  await expect(stevenPage.getByText('Hello Jessica!')).toBeVisible({ timeout: 5000 })

  // Jessica sees the message in real-time (via WebSocket message.new)
  await expect(jessicaPage.getByText('Hello Jessica!')).toBeVisible({ timeout: 10000 })

  // Jessica sends a reply
  const jessicaInput = jessicaPage.locator('.chat-input')
  await jessicaInput.fill('Hi Steven!')
  await jessicaInput.press('Enter')

  // Jessica sees her own message
  await expect(jessicaPage.getByText('Hi Steven!')).toBeVisible({ timeout: 5000 })

  // Steven sees Jessica's reply in real-time
  await expect(stevenPage.getByText('Hi Steven!')).toBeVisible({ timeout: 10000 })

  await stevenCtx.close()
  await jessicaCtx.close()
})
