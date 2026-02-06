/** @type {import('stylelint').Config} */
export default {
  extends: [
    'stylelint-config-standard',
    'stylelint-config-standard-scss',
    'stylelint-config-recess-order',
    'stylelint-config-recommended-vue',
  ],
  ignoreFiles: ['**/*', '!src/**/*'],
  rules: {
    "unit-no-unknown": true,
    'color-no-hex': true,
    'at-rule-no-unknown': null,
    'scss/at-rule-no-unknown': true,
    'font-family-no-missing-generic-family-keyword': true, 
    'selector-class-pattern': [
      // 兼容element-plus组件的类名（el- 开头的）
      '^([a-z][a-z0-9]*)(-[a-z0-9]+)*$|^el-([a-z0-9_,-]+)*$',
      {
        message: (selector) => `Expected class selector "${selector}" to be kebab-case or Element Plus BEM style`,
        severity: 'error',
      },
    ],
  }
}