import { test, expect } from '@playwright/test';

// Unauthenticated access to protected pages must redirect to /login

test.describe('Auth redirects (unauthenticated)', () => {
  test('/cabinet redirects to /login', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL(/\/login/);
  });

  test('/cabinet redirect preserves callbackUrl', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL(/callbackUrl/);
  });

  test('/expertise redirects to /login', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/login/);
  });

  test('/analytics redirects to /login', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/login/);
  });
});
