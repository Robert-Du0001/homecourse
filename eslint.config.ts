import { defineConfig } from 'eslint/config';
import js from '@eslint/js';
import pluginVue from 'eslint-plugin-vue';
import tseslint from 'typescript-eslint';
import globals from 'globals';

export default defineConfig(
  js.configs.recommended,
  tseslint.configs.recommended,
  pluginVue.configs['flat/recommended'],

  {
    files: ['src/**/*.{ts,vue}', '*.ts'],

    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        parser: tseslint.parser,
      },
      globals: {
        ...globals.browser, // 允许使用 window, document 等
        ...globals.node,     // 允许使用 process 等
      },
    },

    rules: {
      'quotes': ['warn', 'single'], // 使用单引号
      'comma-dangle': ['warn', 'always-multiline'], // 在多行后要加逗号
      'semi': 'warn',
      'indent': ['warn', 2],
      'max-len': ['warn', { 'code': 180 }],

      'no-console': ['error', {allow: ['error']}],
      'no-unused-vars': 'error',
      'no-var': 'error', // 禁止使用 var
      'eqeqeq': 'error',
      'eol-last': 'error',

      'vue/multi-word-component-names': 'off', // 组件名用两个及以上的单词（关闭检查）
    },
  },
);
