import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    passWithNoTests: true,
  },
  server: {
    port: 3001,
    proxy: {
      // /api is the dev proxy prefix; rewrite strips it so Traefik receives /v1/... (e.g. /api/v1/auth/login → /v1/auth/login).
      // Traefik then applies stripPrefix per service (e.g. strip-v1-auth → auth service receives /login). Both steps are intentional.
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        // Traefik routes by Host(api.vokzal.tech); without this, localhost:8000 would not match and API calls would fail in local dev
        headers: {
          'Host': 'api.vokzal.tech',
        },
      },
    },
  },
})
