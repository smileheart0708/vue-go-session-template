<template>
  <div class="flex flex-col gap-6">
    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">自动刷新</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">页面可见且网络正常时自动更新数据</p>
      </header>

      <div class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start">
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">启用自动刷新</span>
          <span class="text-[0.85rem] text-text-secondary">关闭后将暂停所有自动刷新任务</span>
        </div>
        <AppSwitch v-model="refreshEnabled" />
      </div>

      <div class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start">
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">刷新间隔</span>
          <span class="text-[0.85rem] text-text-secondary">
            范围 {{ refreshStore.minIntervalSeconds }}-{{ refreshStore.maxIntervalSeconds }} 秒
          </span>
        </div>
        <div class="flex items-center gap-2 max-sm:w-full">
          <input
            v-model.number="intervalSeconds"
            type="number"
            class="w-24 rounded-md px-2.5 py-1.5 text-center max-sm:w-full max-sm:text-left"
            :min="refreshStore.minIntervalSeconds"
            :max="refreshStore.maxIntervalSeconds"
            step="1"
            :disabled="!refreshEnabled"
          />
          <span class="text-[0.9rem] text-text-secondary">秒</span>
        </div>
      </div>
    </section>

    <section class="flex flex-col gap-5 rounded-xl border border-border bg-bg-surface p-6">
      <header class="flex flex-col gap-2">
        <h2 class="m-0 text-[1.1rem] font-semibold text-text-primary">界面偏好</h2>
        <p class="m-0 text-[0.9rem] text-text-secondary">控制仪表板与列表的展示方式</p>
      </header>

      <div class="flex items-center justify-between gap-4 border-t-0 border-border pb-3 pt-0 max-sm:flex-col max-sm:items-start">
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">紧凑布局</span>
          <span class="text-[0.85rem] text-text-secondary">减少卡片和列表的行间距</span>
        </div>
        <AppSwitch v-model="compactLayout" />
      </div>

      <div class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start">
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">显示快捷提示</span>
          <span class="text-[0.85rem] text-text-secondary">在关键按钮旁展示提示引导</span>
        </div>
        <AppSwitch v-model="showTips" />
      </div>

      <div class="flex items-center justify-between gap-4 border-t border-dashed border-border py-3 max-sm:flex-col max-sm:items-start">
        <div class="flex min-w-0 flex-col gap-1.5">
          <span class="font-medium text-text-primary">默认入口</span>
          <span class="text-[0.85rem] text-text-secondary">登录后默认进入的页面</span>
        </div>
        <select v-model="defaultLanding" class="w-[220px] rounded-md bg-bg-surface px-2.5 py-1.5 max-sm:w-full">
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
import { AppSwitch } from '@/components/common'
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
