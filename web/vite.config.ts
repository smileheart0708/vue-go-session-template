import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [
      vue(), // 如果想全局开启 Vapor，可配置 vue({ template: { compilerOptions: { vapor: true } } })
      vueDevTools(),
    ],
    resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
    server: {
      port: 5173,
      host: true,
      proxy: {
        '/api': { target: env.VITE_API_BASE_URL || 'http://localhost:8080', changeOrigin: true },
      },
    },
    build: {
      chunkSizeWarningLimit: 1000,
      rollupOptions: {
        output: {
          manualChunks: {
            // 保持现状很好，但要确保代码里是按需引入 echarts
            vendor: ['vue', 'vue-router', 'pinia', '@vueuse/core'],
          },
        },
      },
    },
  }
})
