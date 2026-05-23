import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}', './lib/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        // ─── ПКТ brand palette ────────────────────────────────────────────
        brand: {
          green: {
            DEFAULT: '#1a5c36',
            dark: '#0f3b22',
            light: '#2d7a50',
            50: '#edf7f1',
            100: '#c8e9d4',
            500: '#1a5c36',
            600: '#155030',
            700: '#0f3b22',
          },
          gold: {
            DEFAULT: '#c8921a',
            dark: '#a07215',
            light: '#e8b040',
            50: '#fdf6e8',
            100: '#f5e0b0',
            500: '#c8921a',
          },
        },
        // ─── Status colors (application FSM) ─────────────────────────────
        status: {
          received: '#6b7280',
          scoring: '#3b82f6',
          security: '#8b5cf6',
          collateral: '#f59e0b',
          legal: '#ec4899',
          analysis: '#06b6d4',
          committee: '#f97316',
          approved: '#10b981',
          rejected: '#ef4444',
          revision: '#eab308',
          documentation: '#6366f1',
          issued: '#059669',
        },
      },
      fontFamily: {
        sans: ['var(--font-inter)', 'system-ui', 'sans-serif'],
      },
      screens: {
        xs: '480px',
      },
      boxShadow: {
        card: '0 1px 3px 0 rgb(0 0 0 / 0.07), 0 1px 2px -1px rgb(0 0 0 / 0.07)',
        'card-hover': '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
      },
    },
  },
  plugins: [],
};

export default config;
