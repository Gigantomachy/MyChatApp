import { test, expect, type Browser } from '@playwright/test'

test('Steven sends friend request to Bob, Bob accepts', async ({ browser }: { browser: Browser }) => {
  const stevenCtx = await browser.newContext({ storageState: 'tests/.auth/steven.json' })
  const bobCtx = await browser.newContext({ storageState: 'tests/.auth/bob.json' })

  const stevenPage = await stevenCtx.newPage()
  const bobPage = await bobCtx.newPage()

  await stevenPage.goto('/')
  await bobPage.goto('/')

  // Steven searches for Bob
  await stevenPage.getByRole('button', { name: 'Search' }).click()
  await stevenPage.getByPlaceholder('Search channels and people...').fill('bob')
  await stevenPage.getByRole('tab', { name: 'People' }).click()

  // Click "Add Friend"
  const addFriendBtn = stevenPage.getByRole('button', { name: 'Add Friend' })
  await expect(addFriendBtn).toBeVisible()
  await addFriendBtn.click()

  // Button should change to "Cancel"
  await expect(stevenPage.getByRole('button', { name: 'Cancel' })).toBeVisible()

  // Bob opens his profile panel
  await bobPage.locator('.top-bar-avatar-btn').click()

  // Bob sees the friend request from Steven
  await expect(bobPage.getByText('Steven Miller')).toBeVisible({ timeout: 10000 })

  // Bob accepts
  await bobPage.getByRole('button', { name: 'Accept' }).click()

  // Bob sees Steven in his friends list
  await expect(bobPage.getByText('Steven Miller')).toHaveCount(1, { timeout: 5000 })

  // Steven opens his profile panel and sees Bob in friends list
  await stevenPage.locator('.top-bar-avatar-btn').click()
  await expect(stevenPage.getByText('Bob Johnson')).toBeVisible({ timeout: 10000 })

  await stevenCtx.close()
  await bobCtx.close()
})
