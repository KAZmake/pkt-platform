// @ts-check
import rootConfig from '../../eslint.config.js';

/** @type {import('typescript-eslint').Config} */
export default [
  ...rootConfig,
  {
    rules: {
      // Next.js специфичные послабления
      '@typescript-eslint/no-explicit-any': 'warn',
    },
  },
];
