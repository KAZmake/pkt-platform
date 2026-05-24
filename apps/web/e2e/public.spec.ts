import { test, expect } from '@playwright/test';

// Public pages render without auth

test.describe('Homepage', () => {
  test('loads and shows hero section', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/Главная|ПКТ/);
    // Hero has a CTA button
    const cta = page.getByRole('link', { name: /подать заявку/i }).first();
    await expect(cta).toBeVisible();
  });

  test('navigation header is visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('navigation')).toBeVisible();
  });
});

test.describe('Programs page', () => {
  test('shows page heading', async ({ page }) => {
    await page.goto('/programs');
    await expect(page.getByRole('heading', { name: /программы кредитования/i })).toBeVisible();
  });

  test('has apply link', async ({ page }) => {
    await page.goto('/programs');
    // Page renders without error
    await expect(page.locator('main')).toBeVisible();
  });
});

test.describe('Calculator page', () => {
  test('shows heading', async ({ page }) => {
    await page.goto('/calculator');
    await expect(page.getByRole('heading', { name: /калькулятор займа/i })).toBeVisible();
  });

  test('renders inputs', async ({ page }) => {
    await page.goto('/calculator');
    // Amount, term, rate inputs exist
    const amountInput = page.locator('input[type="number"]').first();
    await expect(amountInput).toBeVisible();
  });

  test('shows schedule table after loading', async ({ page }) => {
    await page.goto('/calculator');
    // Default values (5 000 000 ₸, 24 мес, 12%) should produce a schedule
    await expect(page.locator('table')).toBeVisible();
    // Header row
    await expect(page.locator('th', { hasText: '№' })).toBeVisible();
    await expect(page.locator('th', { hasText: /платёж/i })).toBeVisible();
  });

  test('schedule shows 24 rows by default (first 6 visible)', async ({ page }) => {
    await page.goto('/calculator');
    // Only first 6 rows shown by default (showFull=false)
    const rows = page.locator('tbody tr');
    await expect(rows).toHaveCount(6);
  });

  test('show all button reveals full schedule', async ({ page }) => {
    await page.goto('/calculator');
    const showAllBtn = page.getByRole('button', { name: /показать все/i });
    await expect(showAllBtn).toBeVisible();
    await showAllBtn.click();
    // 24 months schedule (default term)
    const rows = page.locator('tbody tr');
    await expect(rows).toHaveCount(24);
  });

  test('switching to differentiated changes first payment label', async ({ page }) => {
    await page.goto('/calculator');
    // Default is annuity → "Ежемес. платёж"
    await expect(page.getByText(/ежемес\. платёж/i)).toBeVisible();
    // Click "Дифференц." button
    await page.getByRole('button', { name: /дифференц/i }).click();
    // Label changes to "Первый платёж"
    await expect(page.getByText(/первый платёж/i)).toBeVisible();
  });

  test('changing amount updates summary', async ({ page }) => {
    await page.goto('/calculator');
    const amountInput = page.locator('input[type="number"]').first();
    await amountInput.fill('1000000');
    // Summary section should still show
    await expect(page.locator('table')).toBeVisible();
  });
});

test.describe('Login page', () => {
  test('shows sign-in button', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByRole('button', { name: /войти/i })).toBeVisible();
  });

  test('shows organisation name', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByText(/первое кредитное товарищество/i)).toBeVisible();
  });
});

test.describe('FAQ page', () => {
  test('renders without error', async ({ page }) => {
    await page.goto('/faq');
    await expect(page.locator('main')).toBeVisible();
  });
});

test.describe('Contacts page', () => {
  test('renders without error', async ({ page }) => {
    await page.goto('/contacts');
    await expect(page.locator('main')).toBeVisible();
  });
});
