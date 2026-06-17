import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '127.0.0.1',
    port: 4173,
    proxy: {
      '/user-api': {
        target: 'http://127.0.0.1:8001',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/user-api/, ''),
      },
      '/forum-api': {
        target: 'http://127.0.0.1:8002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/forum-api/, ''),
      },
      '/admin-api': {
        target: 'http://127.0.0.1:8003',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/admin-api/, ''),
      },
    },
  },
})
