<template>
  <div class="settings">
    <h1 class="page-title">设置</h1>
    <section class="settings-card">
      <header class="card-header">
        <h2 class="card-title">自动刷新</h2>
        <p class="card-subtitle">页面可见且网络正常时自动更新数据</p>
      </header>

      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">启用自动刷新</span>
          <span class="setting-desc">关闭后将暂停所有自动刷新任务</span>
        </div>
        <AppSwitch v-model="refreshEnabled" />
      </div>

      <div class="setting-row">
        <div class="setting-info">
          <span class="setting-label">刷新间隔</span>
          <span class="setting-desc">
            范围 {{ refreshStore.minIntervalSeconds }}-{{ refreshStore.maxIntervalSeconds }} 秒
          </span>
        </div>
        <div class="interval-control">
          <input
            v-model.number="intervalSeconds"
            type="number"
            class="interval-input"
            :min="refreshStore.minIntervalSeconds"
            :max="refreshStore.maxIntervalSeconds"
            step="1"
            :disabled="!refreshEnabled"
          />
          <span class="interval-unit">秒</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppSwitch from '@/components/common/AppSwitch.vue'
import { useRefreshStore } from '@/stores/refresh'

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
</script>

<style scoped>
.settings {
  width: 100%;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 1rem;
}

.settings-card {
  padding: 1.5rem;
  background: var(--color-background-elevated);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.card-title {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-text);
}

.card-subtitle {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem 0;
  border-top: 1px dashed var(--color-border);
}

.setting-row:first-of-type {
  border-top: none;
  padding-top: 0;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  min-width: 0;
}

.setting-label {
  font-weight: 500;
  color: var(--color-text);
}

.setting-desc {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.interval-control {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.interval-input {
  width: 88px;
  padding: 6px 10px;
  border-radius: 6px;
  text-align: center;
}

.interval-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.interval-unit {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

@media (max-width: 640px) {
  .setting-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .interval-control {
    width: 100%;
  }

  .interval-input {
    width: 100%;
    text-align: left;
  }
}
</style>
