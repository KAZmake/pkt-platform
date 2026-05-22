// @ts-check
import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import prettierConfig from 'eslint-config-prettier';
import globals from 'globals';

export default tseslint.config(
  // Базовые правила JS
  js.configs.recommended,

  // TypeScript
  ...tseslint.configs.recommended,

  // Prettier (отключает конфликтующие правила ESLint)
  prettierConfig,

  // Глобальные настройки
  {
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2022,
      },
    },
    rules: {
      // TypeScript
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/consistent-type-imports': ['error', { prefer: 'type-imports' }],

      // Общие
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'prefer-const': 'error',
      'no-var': 'error',
    },
  },

  // Игнорируемые файлы
  {
    ignores: [
      '**/node_modules/**',
      '**/.next/**',
      '**/.expo/**',
      '**/dist/**',
      '**/out/**',
      '**/build/**',
      '**/.turbo/**',
      'apps/api/**',         // Go — не нужен ESLint
      '_docs/**',
    ],
  },
);
