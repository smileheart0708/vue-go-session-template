import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { compression } from 'vite-plugin-compression2'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [
      vue(), // 如果想全局开启 Vapor，可配置 vue({ template: { compilerOptions: { vapor: true } } })
      vueDevTools(),
      // 生成 .zst 文件 (Zstd) - 2026 优先级最高
      compression({
        algorithms: ['zstd'],
        exclude: [/\.(zst)$/, /\.(br)$/, /\.(gz)$/],
        deleteOriginalAssets: false, // 必须保留原文件，作为不支持压缩的浏览器的兜底
      }),
      // 生成 .br 文件 (Brotli) - 兼容性兜底
      compression({
        algorithms: ['brotliCompress'],
        exclude: [/\.(zst)$/, /\.(br)$/, /\.(gz)$/],
        deleteOriginalAssets: false,
      }),
      // 生成 .gz 文件 (Gzip) - 兼容性兜底
      compression({ algorithms: ['gzip'], exclude: [/\.(zst)$/, /\.(br)$/, /\.(gz)$/] }),
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
