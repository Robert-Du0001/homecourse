import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import checker from 'vite-plugin-checker';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    checker({
      typescript: true, // 检查 TS
      eslint: {
        lintCommand: 'eslint "./src/**/*.{ts,vue}"', // 检查 ESLint
      },
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  publicDir: './src/static',
  build: {
    outDir: './public', // 确保指向 Goravel 的 public
    emptyOutDir: true,
  },
  server: {
    proxy: {
      // 代理到后端服务，并处理跨域
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
      '/media': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      }
    }
  }
})
