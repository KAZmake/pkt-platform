import { test, expect } from '@playwright/test';
import path from 'path';

// Borrower cabinet — authenticated as 'borrower' role

test.use({ storageState: path.join(import.meta.dirname, '.auth/borrower.json') });

test.describe('Cabinet dashboard (borrower)', () => {
  test('shows dashboard heading', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page.getByRole('heading', { name: /добро пожаловать/i })).toBeVisible();
  });

  test('shows "Личный кабинет заёмщика" subtitle', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page.getByText(/личный кабинет заёмщика/i)).toBeVisible();
  });

  test('shows loan summary stat cards', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page.getByText(/остаток долга/i)).toBeVisible();
    await expect(page.getByText(/следующий платёж/i)).toBeVisible();
    await expect(page.getByText(/документов/i)).toBeVisible();
    await expect(page.getByText(/уведомлений/i)).toBeVisible();
  });

  test('shows loan card with program name', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page.getByText(/агробизнес 2025/i)).toBeVisible();
  });

  test('shows progress bar', async ({ page }) => {
    await page.goto('/cabinet');
    // Progress bar div exists
    const bar = page.locator('.bg-brand-green.rounded-full').first();
    await expect(bar).toBeVisible();
  });

  test('quick actions are clickable links', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page.getByRole('link', { name: /создать обращение/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /уведомления/i }).first()).toBeVisible();
  });

  test('sidebar is rendered', async ({ page }) => {
    await page.goto('/cabinet');
    // Sidebar navigation present
    await expect(page.locator('nav, aside, [class*="sidebar"]').first()).toBeVisible();
  });
});

test.describe('Cabinet navigation (borrower)', () => {
  test('can navigate to schedule page', async ({ page }) => {
    await page.goto('/cabinet/schedule');
    await expect(page.locator('main')).toBeVisible();
    await expect(page).toHaveURL('/cabinet/schedule');
  });

  test('can navigate to notifications page', async ({ page }) => {
    await page.goto('/cabinet/notifications');
    await expect(page.locator('main')).toBeVisible();
    await expect(page).toHaveURL('/cabinet/notifications');
  });

  test('can navigate to documents page', async ({ page }) => {
    await page.goto('/cabinet/documents');
    await expect(page.locator('main')).toBeVisible();
    await expect(page).toHaveURL('/cabinet/documents');
  });

  test('can navigate to tickets page', async ({ page }) => {
    await page.goto('/cabinet/tickets');
    await expect(page.locator('main')).toBeVisible();
    await expect(page).toHaveURL('/cabinet/tickets');
  });
});
