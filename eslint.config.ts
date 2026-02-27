import js from '@eslint/js';
import stylistic from '@stylistic/eslint-plugin';
import { defineConfig } from 'eslint/config';
import eslintConfigPrettier from 'eslint-config-prettier';
import importPlugin from 'eslint-plugin-import';
import pluginVue from 'eslint-plugin-vue';
<<<<<<< HEAD
=======
import tseslint from 'typescript-eslint';
import stylistic from '@stylistic/eslint-plugin';
>>>>>>> main
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default defineConfig(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/recommended'],

  {
    plugins: {
      '@stylistic': stylistic,
<<<<<<< HEAD
      import: importPlugin,
=======
>>>>>>> main
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
<<<<<<< HEAD
      eqeqeq: 'error', // 必须使用全等
      'no-console': [
        process.env.NODE_ENV === 'production' ? 'error' : 'warn',
        { allow: ['error', 'warn'] },
      ],

      '@stylistic/spaced-comment': 'warn',
      '@stylistic/member-delimiter-style': 'warn',
      '@stylistic/lines-between-class-members': 'warn',
=======
      'eqeqeq': 'error', // 必须使用全等
      'no-console': ['error', { allow: ['error'] }],

      '@stylistic/quotes': ['warn', 'single'], // 使用单引号
      '@stylistic/comma-dangle': ['warn', 'always-multiline'], // 在多行后要加逗号
      '@stylistic/comma-spacing': ['warn', { 'before': false, 'after': true }],
      '@stylistic/comma-style': ['warn', 'last'],
      '@stylistic/semi': 'warn',
      '@stylistic/indent': ['warn', 2],
      '@stylistic/max-len': ['warn', { 'code': 180 }],
      '@stylistic/object-curly-spacing': ['warn', 'always'],
      '@stylistic/space-infix-ops': 'warn',
      '@stylistic/eol-last': 'warn',
>>>>>>> main

      'vue/multi-word-component-names': 'off', // 组件名用两个及以上的单词（关闭检查）

      // import顺序规范
      'import/order': [
        'warn',
        {
          // 定义分组顺序
          groups: [
            'builtin', // node 内置模块 (path, fs 等)
            'external', // 外部库 (vue, element-plus 等)
            'internal', // 内部别名路径 (@/components 等)
            ['parent', 'sibling'], // 父级及同级目录文件
            'index', // 索引文件
            'object', // object 类型的 import
            'type', // 类型导入 (import type)
          ],
          'newlines-between': 'always', // 强制分组之间换行
          // 按字母顺序排序
          alphabetize: {
            order: 'asc',
            caseInsensitive: true,
          },
        },
      ],
    },
  },

  eslintConfigPrettier,
);
