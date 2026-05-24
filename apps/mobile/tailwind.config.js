/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./app/**/*.{js,jsx,ts,tsx}', './components/**/*.{js,jsx,ts,tsx}'],
  presets: [require('nativewind/preset')],
  theme: {
    extend: {
      colors: {
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
    },
  },
  plugins: [],
};
