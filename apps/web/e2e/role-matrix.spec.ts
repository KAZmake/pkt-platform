/**
 * Role-matrix access control — automated verification.
 *
 * Frontend (Next.js middleware) rules:
 *   /cabinet   → borrower, employee, expert, admin
 *   /expertise → employee, expert, admin
 *   /analytics → employee, expert, admin
 *   unauthenticated → /login (redirect)
 *   wrong role → /403 (redirect)
 */
import { test, expect } from '@playwright/test';
import path from 'path';

// ─── helpers ──────────────────────────────────────────────────────────────────

function withRole(role: string) {
  return path.join(import.meta.dirname, `.auth/${role}.json`);
}

// ─── unauthenticated (public) ─────────────────────────────────────────────────

test.describe('public — no auth', () => {
  test('/cabinet → redirects to /login', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL(/\/login/);
  });

  test('/expertise → redirects to /login', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/login/);
  });

  test('/analytics → redirects to /login', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/login/);
  });

  test('public pages accessible without auth', async ({ page }) => {
    for (const url of ['/', '/programs', '/calculator', '/faq', '/contacts']) {
      await page.goto(url);
      await expect(page.locator('main')).toBeVisible();
    }
  });
});

// ─── borrower ─────────────────────────────────────────────────────────────────

test.describe('borrower role', () => {
  test.use({ storageState: withRole('borrower') });

  test('/cabinet → accessible', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL('/cabinet');
    await expect(page.locator('main')).toBeVisible();
  });

  test('/expertise → 403 (no access)', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL('/403');
    await expect(page.getByText('403')).toBeVisible();
  });

  test('/analytics → 403 (no access)', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL('/403');
    await expect(page.getByText('403')).toBeVisible();
  });

  test('public pages accessible', async ({ page }) => {
    await page.goto('/programs');
    await expect(page.locator('main')).toBeVisible();
  });
});

// ─── employee ─────────────────────────────────────────────────────────────────

test.describe('employee role', () => {
  test.use({ storageState: withRole('employee') });

  test('/cabinet → accessible (employee can view cabinet)', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL('/cabinet');
    await expect(page.locator('main')).toBeVisible();
  });

  test('/expertise → accessible', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/expertise/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('/analytics → accessible', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/analytics/);
    await expect(page.locator('main')).toBeVisible();
  });
});

// ─── expert ───────────────────────────────────────────────────────────────────

test.describe('expert role', () => {
  test.use({ storageState: withRole('expert') });

  test('/expertise → accessible', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/expertise/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('/analytics → accessible', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/analytics/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('/cabinet → accessible', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL('/cabinet');
    await expect(page.locator('main')).toBeVisible();
  });
});

// ─── admin ────────────────────────────────────────────────────────────────────

test.describe('admin role', () => {
  test.use({ storageState: withRole('admin') });

  test('/cabinet → accessible', async ({ page }) => {
    await page.goto('/cabinet');
    await expect(page).toHaveURL('/cabinet');
    await expect(page.locator('main')).toBeVisible();
  });

  test('/expertise → accessible', async ({ page }) => {
    await page.goto('/expertise');
    await expect(page).toHaveURL(/\/expertise/);
    await expect(page.locator('main')).toBeVisible();
  });

  test('/analytics → accessible', async ({ page }) => {
    await page.goto('/analytics');
    await expect(page).toHaveURL(/\/analytics/);
    await expect(page.locator('main')).toBeVisible();
  });
});

// ─── 403 page itself ──────────────────────────────────────────────────────────

test.describe('403 page', () => {
  test('renders correctly', async ({ page }) => {
    await page.goto('/403');
    await expect(page.getByText('403')).toBeVisible();
    await expect(page.getByRole('heading', { name: /доступ запрещён/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /на главную/i })).toBeVisible();
  });
});
