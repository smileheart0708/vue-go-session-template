<template>
  <div class="flex flex-col gap-6">
    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">上游服务</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">管理调用地址、超时与健康检查</p>
      </header>

      <div
        class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">请求地址</span>
          <span class="text-[0.85rem] text-text-secondary">用于拉取模型与统计数据的主入口</span>
        </div>
        <input
          v-model="baseUrl"
          type="url"
          class="w-55 rounded-md px-2.5 py-1.5 max-sm:w-full"
          placeholder="https://api.example.com"
        />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">请求超时</span>
          <span class="text-[0.85rem] text-text-secondary">单位秒，建议 10-60 秒</span>
        </div>
        <div class="flex items-center gap-2 max-sm:w-full">
          <input
            v-model.number="timeoutSeconds"
            type="number"
            min="5"
            max="120"
            step="1"
            class="w-24 rounded-md px-2.5 py-1.5 text-center max-sm:w-full max-sm:text-left"
          />
          <span class="text-[0.9rem] text-text-secondary">秒</span>
        </div>
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">健康检查</span>
          <span class="text-[0.85rem] text-text-secondary">后台定时探测上游服务状态</span>
        </div>
        <AppSwitch v-model="healthCheckEnabled" />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">流量模式</span>
          <span class="text-[0.85rem] text-text-secondary">不同模式会影响并发与队列策略</span>
        </div>
        <select
          v-model="trafficMode"
          class="w-55 rounded-md bg-bg-surface px-2.5 py-1.5 max-sm:w-full"
        >
          <option value="balanced">均衡</option>
          <option value="low-latency">低延迟</option>
          <option value="high-throughput">高吞吐</option>
        </select>
      </div>
    </section>

    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">鉴权与重试</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">保护接口并控制失败重试策略</p>
      </header>

      <div
        class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">API Token</span>
          <span class="text-[0.85rem] text-text-secondary">建议使用只读令牌</span>
        </div>
        <input
          v-model="apiToken"
          type="password"
          class="w-55 rounded-md px-2.5 py-1.5 max-sm:w-full"
          placeholder="sk-..."
        />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">签名校验</span>
          <span class="text-[0.85rem] text-text-secondary">对请求体进行签名校验</span>
        </div>
        <AppSwitch v-model="signatureCheckEnabled" />
      </div>

      <div
        class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start"
      >
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">失败重试</span>
          <span class="text-[0.85rem] text-text-secondary">网络异常时自动尝试重新请求</span>
        </div>
        <select
          v-model="retryPolicy"
          class="w-55 rounded-md bg-bg-surface px-2.5 py-1.5 max-sm:w-full"
        >
          <option value="none">关闭</option>
          <option value="conservative">保守</option>
          <option value="aggressive">激进</option>
        </select>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { z } from 'zod'
import { AppSwitch } from '@/components/common'
import { useValidatedLocalStorage } from '@/composables/useValidatedLocalStorage'

const trafficModeSchema = z.enum(['balanced', 'low-latency', 'high-throughput'])
const retryPolicySchema = z.enum(['none', 'conservative', 'aggressive'])

const upstreamBaseUrlSchema = z.string().trim().min(1).max(2048)

const timeoutSecondsSchema = z.number().finite().transform((value) => {
  const rounded = Math.round(value)
  if (rounded < 5) return 5
  if (rounded > 120) return 120
  return rounded
})

type TrafficMode = z.infer<typeof trafficModeSchema>
type RetryPolicy = z.infer<typeof retryPolicySchema>

const baseUrl = useValidatedLocalStorage(
  'settings.upstream_base_url',
  upstreamBaseUrlSchema,
  'https://api.example.com',
)
const timeoutSeconds = useValidatedLocalStorage(
  'settings.upstream_timeout_seconds',
  timeoutSecondsSchema,
  20,
)
const healthCheckEnabled = useValidatedLocalStorage(
  'settings.upstream_health_check',
  z.boolean(),
  true,
)
const trafficMode = useValidatedLocalStorage<TrafficMode>(
  'settings.upstream_traffic_mode',
  trafficModeSchema,
  'balanced',
)
const apiToken = useValidatedLocalStorage('settings.upstream_api_token', z.string(), '')
const signatureCheckEnabled = useValidatedLocalStorage(
  'settings.upstream_signature_check',
  z.boolean(),
  false,
)
const retryPolicy = useValidatedLocalStorage<RetryPolicy>(
  'settings.upstream_retry_policy',
  retryPolicySchema,
  'conservative',
)
</script>
