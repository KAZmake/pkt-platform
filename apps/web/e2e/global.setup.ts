import { test as setup } from '@playwright/test';
import path from 'path';
import fs from 'fs/promises';

const AUTH_DIR = path.join(import.meta.dirname, '.auth');

setup.beforeAll(async () => {
  await fs.mkdir(AUTH_DIR, { recursive: true });
});

setup('create borrower auth state', async ({ page }) => {
  await page.goto('/api/e2e-auth?role=borrower');
  await page.waitForSelector('text=true'); // JSON { ok: true }
  await page.context().storageState({ path: path.join(AUTH_DIR, 'borrower.json') });
});

setup('create employee auth state', async ({ page }) => {
  await page.goto('/api/e2e-auth?role=employee');
  await page.waitForSelector('text=true');
  await page.context().storageState({ path: path.join(AUTH_DIR, 'employee.json') });
});

setup('create expert auth state', async ({ page }) => {
  await page.goto('/api/e2e-auth?role=expert');
  await page.waitForSelector('text=true');
  await page.context().storageState({ path: path.join(AUTH_DIR, 'expert.json') });
});

setup('create admin auth state', async ({ page }) => {
  await page.goto('/api/e2e-auth?role=admin');
  await page.waitForSelector('text=true');
  await page.context().storageState({ path: path.join(AUTH_DIR, 'admin.json') });
});
