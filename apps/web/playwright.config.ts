import { defineConfig, devices } from '@playwright/test';

const BASE_URL = process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:3100';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: process.env.CI ? 'github' : 'list',

  use: {
    baseURL: BASE_URL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    // Longer timeout for Next.js SSR pages
    navigationTimeout: 30_000,
    actionTimeout: 10_000,
  },

  projects: [
    // Setup project: create auth storage states
    {
      name: 'setup',
      testMatch: /global\.setup\.ts/,
    },
    // Main project depends on setup
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      dependencies: ['setup'],
    },
  ],

  webServer: {
    command: 'npx next dev -p 3100',
    url: BASE_URL,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
    env: {
      E2E_TEST: 'true',
      NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET ?? 'change-me-in-production',
      NEXTAUTH_URL: process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:3100',
    },
  },
});
