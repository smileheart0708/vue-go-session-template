<template>
  <div class="flex flex-col gap-6">
    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">下游代理</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">控制转发端口、协议模式与来源限制</p>
      </header>

      <div
        class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">启用代理</span>
          <span class="text-[0.85rem] text-text-secondary">开启后将对外提供统一入口</span>
        </div>
        <AppSwitch v-model="proxyEnabled" />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">监听端口</span>
          <span class="text-[0.85rem] text-text-secondary">建议使用 1024-65535 范围</span>
        </div>
        <div class="flex items-center gap-2 max-sm:w-full">
          <input
            v-model.number="listenPort"
            type="number"
            min="1024"
            max="65535"
            step="1"
            class="w-24 rounded-md px-2.5 py-1.5 text-center max-sm:w-full max-sm:text-left"
          />
          <span class="text-[0.9rem] text-text-secondary">端口</span>
        </div>
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">代理模式</span>
          <span class="text-[0.85rem] text-text-secondary">决定请求是否改写或镜像</span>
        </div>
        <select
          v-model="proxyMode"
          class="w-55 rounded-md bg-bg-surface px-2.5 py-1.5 max-sm:w-full"
        >
          <option value="transparent">透明转发</option>
          <option value="rewrite">路径改写</option>
          <option value="mirror">流量镜像</option>
        </select>
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">允许来源</span>
          <span class="text-[0.85rem] text-text-secondary">用英文逗号分隔多个域名</span>
        </div>
        <textarea
          v-model="allowedOrigins"
          class="min-h-16 w-65 resize-y rounded-md px-2.5 py-1.5 max-sm:w-full"
          placeholder="https://console.example.com, https://ops.example.com"
        />
      </div>
    </section>

    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">缓存与限速</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">提升稳定性并避免突发流量</p>
      </header>

      <div
        class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">启用缓存</span>
          <span class="text-[0.85rem] text-text-secondary">短时间内复用上游响应结果</span>
        </div>
        <AppSwitch v-model="cacheEnabled" />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">缓存时长</span>
          <span class="text-[0.85rem] text-text-secondary">缓存命中后重复请求将直接返回</span>
        </div>
        <div class="flex items-center gap-2 max-sm:w-full">
          <input
            v-model.number="cacheTtlSeconds"
            type="number"
            min="5"
            max="600"
            step="5"
            class="w-24 rounded-md px-2.5 py-1.5 text-center max-sm:w-full max-sm:text-left"
          />
          <span class="text-[0.9rem] text-text-secondary">秒</span>
        </div>
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">限速等级</span>
          <span class="text-[0.85rem] text-text-secondary">限制单账号的最大并发</span>
        </div>
        <select
          v-model="rateLimitLevel"
          class="w-55 rounded-md bg-bg-surface px-2.5 py-1.5 max-sm:w-full"
        >
          <option value="off">关闭</option>
          <option value="soft">温和</option>
          <option value="strict">严格</option>
        </select>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { z } from 'zod'
import { AppSwitch } from '@/components/common'
import { useValidatedLocalStorage } from '@/composables/useValidatedLocalStorage'

const proxyModeSchema = z.enum(['transparent', 'rewrite', 'mirror'])
const rateLimitLevelSchema = z.enum(['off', 'soft', 'strict'])

const listenPortSchema = z.number().finite().transform((value) => {
  const rounded = Math.round(value)
  if (rounded < 1024) return 1024
  if (rounded > 65535) return 65535
  return rounded
})

const cacheTtlSecondsSchema = z.number().finite().transform((value) => {
  const roundedToStep = Math.round(value / 5) * 5
  if (roundedToStep < 5) return 5
  if (roundedToStep > 600) return 600
  return roundedToStep
})

type ProxyMode = z.infer<typeof proxyModeSchema>
type RateLimitLevel = z.infer<typeof rateLimitLevelSchema>

const proxyEnabled = useValidatedLocalStorage('settings.proxy_enabled', z.boolean(), true)
const listenPort = useValidatedLocalStorage('settings.proxy_listen_port', listenPortSchema, 8080)
const proxyMode = useValidatedLocalStorage<ProxyMode>(
  'settings.proxy_mode',
  proxyModeSchema,
  'transparent',
)
const allowedOrigins = useValidatedLocalStorage(
  'settings.proxy_allowed_origins',
  z.string(),
  '',
)
const cacheEnabled = useValidatedLocalStorage('settings.proxy_cache_enabled', z.boolean(), true)
const cacheTtlSeconds = useValidatedLocalStorage(
  'settings.proxy_cache_ttl_seconds',
  cacheTtlSecondsSchema,
  60,
)
const rateLimitLevel = useValidatedLocalStorage<RateLimitLevel>(
  'settings.proxy_rate_limit',
  rateLimitLevelSchema,
  'soft',
)
</script>
