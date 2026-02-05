<template>
  <div class="settings-section">
    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">下游代理</h2>
        <p class="settings-card__subtitle">控制转发端口、协议模式与来源限制</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">启用代理</span>
          <span class="settings-desc">开启后将对外提供统一入口</span>
        </div>
        <AppSwitch v-model="proxyEnabled" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">监听端口</span>
          <span class="settings-desc">建议使用 1024-65535 范围</span>
        </div>
        <div class="settings-control">
          <input
            v-model.number="listenPort"
            type="number"
            min="1024"
            max="65535"
            step="1"
            class="settings-input settings-input--short"
          />
          <span class="settings-unit">端口</span>
        </div>
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">代理模式</span>
          <span class="settings-desc">决定请求是否改写或镜像</span>
        </div>
        <select v-model="proxyMode" class="settings-select">
          <option value="transparent">透明转发</option>
          <option value="rewrite">路径改写</option>
          <option value="mirror">流量镜像</option>
        </select>
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">允许来源</span>
          <span class="settings-desc">用英文逗号分隔多个域名</span>
        </div>
        <textarea
          v-model="allowedOrigins"
          class="settings-textarea"
          placeholder="https://console.example.com, https://ops.example.com"
        />
      </div>
    </section>

    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">缓存与限速</h2>
        <p class="settings-card__subtitle">提升稳定性并避免突发流量</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">启用缓存</span>
          <span class="settings-desc">短时间内复用上游响应结果</span>
        </div>
        <AppSwitch v-model="cacheEnabled" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">缓存时长</span>
          <span class="settings-desc">缓存命中后重复请求将直接返回</span>
        </div>
        <div class="settings-control">
          <input
            v-model.number="cacheTtlSeconds"
            type="number"
            min="5"
            max="600"
            step="5"
            class="settings-input settings-input--short"
          />
          <span class="settings-unit">秒</span>
        </div>
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">限速等级</span>
          <span class="settings-desc">限制单账号的最大并发</span>
        </div>
        <select v-model="rateLimitLevel" class="settings-select">
          <option value="off">关闭</option>
          <option value="soft">温和</option>
          <option value="strict">严格</option>
        </select>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import AppSwitch from '@/components/common/AppSwitch.vue'

type ProxyMode = 'transparent' | 'rewrite' | 'mirror'
type RateLimitLevel = 'off' | 'soft' | 'strict'

const proxyEnabled = useLocalStorage('settings.proxy_enabled', true)
const listenPort = useLocalStorage('settings.proxy_listen_port', 8080)
const proxyMode = useLocalStorage<ProxyMode>('settings.proxy_mode', 'transparent')
const allowedOrigins = useLocalStorage('settings.proxy_allowed_origins', '')
const cacheEnabled = useLocalStorage('settings.proxy_cache_enabled', true)
const cacheTtlSeconds = useLocalStorage('settings.proxy_cache_ttl_seconds', 60)
const rateLimitLevel = useLocalStorage<RateLimitLevel>('settings.proxy_rate_limit', 'soft')
</script>
