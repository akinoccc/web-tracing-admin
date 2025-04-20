import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import tailwindcss from '@tailwindcss/vite'


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/trackweb': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/getAllTracingList': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/cleanTracingList': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/getBaseInfo': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
