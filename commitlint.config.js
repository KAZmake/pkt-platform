/** @type {import('@commitlint/types').UserConfig} */
export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    // Типы коммитов разрешённые в проекте
    'type-enum': [
      2,
      'always',
      [
        'feat',     // новая функциональность
        'fix',      // исправление бага
        'docs',     // документация
        'style',    // форматирование (не логика)
        'refactor', // рефакторинг
        'test',     // тесты
        'chore',    // инфраструктура, конфиги
        'ci',       // CI/CD
        'perf',     // производительность
        'revert',   // откат коммита
      ],
    ],
    'subject-case': [2, 'never', ['upper-case', 'pascal-case', 'snake-case']],
    'header-max-length': [2, 'always', 100],
    'body-max-line-length': [2, 'always', 200],
  },
};
