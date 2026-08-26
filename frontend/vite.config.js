import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [react(), tailwindcss()],

    server: {
      proxy: {
        '/employees': {
          target: env.VITE_API_URL,
          changeOrigin: true,
        },

        '/employee': {
          target: env.VITE_API_URL,
          changeOrigin: true,
        },
      },
    },
  }
})