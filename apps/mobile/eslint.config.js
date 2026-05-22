// @ts-check
import rootConfig from '../../eslint.config.js';

/** @type {import('typescript-eslint').Config} */
export default [
  ...rootConfig,
  {
    rules: {
      // React Native специфичные послабления
      'no-console': 'off', // логи в RN нужны для отладки
    },
  },
];
