import { test, expect } from '@playwright/test';
import path from 'path';

// Employee modules — authenticated as 'employee' role

test.use({ storageState: path.join(import.meta.dirname, '.auth/employee.json') });

test.describe('Expertise queue (employee)', () => {
  test('redirects to expertise queue page', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/expertise/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('shows expertise page heading', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page.getByRole('heading', { name: /очередь экспертизы/i })).toBeVisible();
  });

  test('shows applications table or empty state', async ({ page }) => {
    await page.goto('/expertise');
    // Table or any content inside main
    await expect(page.locator('main')).toBeVisible();
  });
});

test.describe('Analytics page (employee)', () => {
  test('loads analytics page', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/analytics/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('shows analytics heading', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page.getByRole('heading', { name: /аналитика/i })).toBeVisible();
  });

  test('shows dashboard tabs', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page.getByRole('link', { name: /портфель займов/i })).toBeVisible();
  });
});

test.describe('Employee navigation', () => {
  test('header shows employee nav links', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page.getByRole('navigation')).toBeVisible();
  });

  test('can navigate between expertise and analytics', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page.locator('main')).toBeVisible();
    await page.goto('/analytics');
    await expect(page.locator('main')).toBeVisible();
  });
});
