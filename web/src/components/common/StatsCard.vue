<template>
  <div class="stats-card">
    <div class="card-bg"></div>
    <div class="stats-icon">
      <slot name="icon">
        <Circle :size="28" />
      </slot>
    </div>
    <div class="stats-content">
      <p class="stats-label">{{ label }}</p>
      <p class="stats-value">{{ value }}</p>
      <p v-if="change" class="stats-change">{{ change }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Circle } from 'lucide-vue-next'

defineOptions({
  name: 'StatsCard',
})

interface Props {
  label: string
  value: string | number
  change?: string
  color: string
}

const { label, value, change = '', color } = defineProps<Props>()
</script>

<style scoped>
.stats-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--color-background-elevated);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
  max-width: 100%;
  overflow: hidden;
  z-index: 1;
}

.card-bg {
  position: absolute;
  top: 1.5rem;
  left: 1.5rem;
  width: 48px;
  height: 48px;
  border-radius: 10px;
  z-index: -1;
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  background: v-bind('color');
}

/* 仅在支持真正 hover 的设备上启用 hover 效果 */
@media (hover: hover) {
  .stats-card:hover {
    border-color: transparent;
  }

  .stats-card:hover .card-bg {
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 12px;
    transform: scale(1.1);
  }

  .stats-card:hover .stats-label,
  .stats-card:hover .stats-value,
  .stats-card:hover .stats-change {
    color: white;
  }

  .stats-card:hover .stats-icon {
    transform: scale(1.5);
  }
}

.stats-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  padding: 10px;
  border-radius: 10px;
  color: white;
  flex-shrink: 0;
  box-sizing: border-box;
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.stats-icon svg {
  width: 100%;
  height: 100%;
}

.stats-content {
  flex: 1;
  min-width: 0;
}

.stats-label {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0 0 0.25rem;
  transition: color 0.3s ease;
}

.stats-value {
  font-size: clamp(1.25rem, 2vw, 1.875rem);
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 0.25rem;
  line-height: 1.2;
  transition: color 0.3s ease;
}

.stats-change {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  transition: color 0.3s ease;
}
</style>
