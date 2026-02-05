<template>
  <div class="settings-section">
    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">自动刷新</h2>
        <p class="settings-card__subtitle">页面可见且网络正常时自动更新数据</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">启用自动刷新</span>
          <span class="settings-desc">关闭后将暂停所有自动刷新任务</span>
        </div>
        <AppSwitch v-model="refreshEnabled" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">刷新间隔</span>
          <span class="settings-desc">
            范围 {{ refreshStore.minIntervalSeconds }}-{{ refreshStore.maxIntervalSeconds }} 秒
          </span>
        </div>
        <div class="settings-control">
          <input
            v-model.number="intervalSeconds"
            type="number"
            class="settings-input settings-input--short"
            :min="refreshStore.minIntervalSeconds"
            :max="refreshStore.maxIntervalSeconds"
            step="1"
            :disabled="!refreshEnabled"
          />
          <span class="settings-unit">秒</span>
        </div>
      </div>
    </section>

    <section class="settings-card">
      <header class="settings-card__header">
        <h2 class="settings-card__title">界面偏好</h2>
        <p class="settings-card__subtitle">控制仪表板与列表的展示方式</p>
      </header>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">紧凑布局</span>
          <span class="settings-desc">减少卡片和列表的行间距</span>
        </div>
        <AppSwitch v-model="compactLayout" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">显示快捷提示</span>
          <span class="settings-desc">在关键按钮旁展示提示引导</span>
        </div>
        <AppSwitch v-model="showTips" />
      </div>

      <div class="settings-row">
        <div class="settings-info">
          <span class="settings-label">默认入口</span>
          <span class="settings-desc">登录后默认进入的页面</span>
        </div>
        <select v-model="defaultLanding" class="settings-select">
          <option value="dashboard">仪表板</option>
          <option value="logs">日志</option>
          <option value="settings">设置</option>
        </select>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import AppSwitch from '@/components/common/AppSwitch.vue'
import { useRefreshStore } from '@/stores/refresh'

type LandingPage = 'dashboard' | 'logs' | 'settings'

const refreshStore = useRefreshStore()

const refreshEnabled = computed({
  get: () => refreshStore.isEnabled,
  set: (value: boolean) => {
    refreshStore.isEnabled = value
  },
})

const intervalSeconds = computed({
  get: () => refreshStore.intervalSeconds,
  set: (value: number | null) => {
    refreshStore.setIntervalSeconds(value ?? refreshStore.intervalSeconds)
  },
})

const compactLayout = useLocalStorage('settings.ui_compact', false)
const showTips = useLocalStorage('settings.ui_show_tips', true)
const defaultLanding = useLocalStorage<LandingPage>('settings.default_landing', 'dashboard')
</script>
