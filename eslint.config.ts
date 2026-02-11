import { defineConfig } from 'eslint/config';
import js from '@eslint/js';
import pluginVue from 'eslint-plugin-vue';
import tseslint from 'typescript-eslint';
import stylistic from '@stylistic/eslint-plugin';
import eslintConfigPrettier from 'eslint-config-prettier';
import globals from 'globals';

export default defineConfig(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/recommended'],

  {
    plugins: {
      '@stylistic': stylistic,
    },

    files: ['src/**/*.{ts,vue}', '*.ts'],

    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        parser: tseslint.parser,
      },
      globals: {
        ...globals.browser, // 允许使用 window, document 等
        ...globals.node, // 允许使用 process 等
      },
    },

    rules: {
      'no-var': 'error', // 禁止使用 var
      'no-eval': 'error',
      eqeqeq: 'error', // 必须使用全等
      'no-console': ['error', { allow: ['error'] }],

      '@stylistic/spaced-comment': 'warn',
      '@stylistic/member-delimiter-style': 'warn',
      '@stylistic/lines-between-class-members': 'warn',

      'vue/multi-word-component-names': 'off', // 组件名用两个及以上的单词（关闭检查）
    },
  },

  eslintConfigPrettier,
);
