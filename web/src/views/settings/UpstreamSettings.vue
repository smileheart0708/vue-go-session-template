<template>
  <div class="settings-section">
    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">上游服务</h2>
        <p class="settings-card__subtitle">管理调用地址、超时与健康检查</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">请求地址</span>
          <span class="settings-desc">用于拉取模型与统计数据的主入口</span>
        </div>
        <input
          v-model="baseUrl"
          type="url"
          class="settings-input"
          placeholder="https://api.example.com"
        />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">请求超时</span>
          <span class="settings-desc">单位秒，建议 10-60 秒</span>
        </div>
        <div class="settings-control">
          <input
            v-model.number="timeoutSeconds"
            type="number"
            min="5"
            max="120"
            step="1"
            class="settings-input settings-input--short"
          />
          <span class="settings-unit">秒</span>
        </div>
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">健康检查</span>
          <span class="settings-desc">后台定时探测上游服务状态</span>
        </div>
        <AppSwitch v-model="healthCheckEnabled" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">流量模式</span>
          <span class="settings-desc">不同模式会影响并发与队列策略</span>
        </div>
        <select v-model="trafficMode" class="settings-select">
          <option value="balanced">均衡</option>
          <option value="low-latency">低延迟</option>
          <option value="high-throughput">高吞吐</option>
        </select>
      </div>
    </section>

    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">鉴权与重试</h2>
        <p class="settings-card__subtitle">保护接口并控制失败重试策略</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">API Token</span>
          <span class="settings-desc">建议使用只读令牌</span>
        </div>
        <input v-model="apiToken" type="password" class="settings-input" placeholder="sk-..." />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">签名校验</span>
          <span class="settings-desc">对请求体进行签名校验</span>
        </div>
        <AppSwitch v-model="signatureCheckEnabled" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">失败重试</span>
          <span class="settings-desc">网络异常时自动尝试重新请求</span>
        </div>
        <select v-model="retryPolicy" class="settings-select">
          <option value="none">关闭</option>
          <option value="conservative">保守</option>
          <option value="aggressive">激进</option>
        </select>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import AppSwitch from '@/components/common/AppSwitch.vue'

type TrafficMode = 'balanced' | 'low-latency' | 'high-throughput'
type RetryPolicy = 'none' | 'conservative' | 'aggressive'

const baseUrl = useLocalStorage('settings.upstream_base_url', 'https://api.example.com')
const timeoutSeconds = useLocalStorage('settings.upstream_timeout_seconds', 20)
const healthCheckEnabled = useLocalStorage('settings.upstream_health_check', true)
const trafficMode = useLocalStorage<TrafficMode>('settings.upstream_traffic_mode', 'balanced')
const apiToken = useLocalStorage('settings.upstream_api_token', '')
const signatureCheckEnabled = useLocalStorage('settings.upstream_signature_check', false)
const retryPolicy = useLocalStorage<RetryPolicy>('settings.upstream_retry_policy', 'conservative')
</script>
