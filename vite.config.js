import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
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
