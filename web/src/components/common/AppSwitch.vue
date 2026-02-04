<script setup lang="ts">
/**
 * WinSwitch.vue
 * Windows 11 风格的开关组件
 */
import { computed } from 'vue'

interface Props {
  disabled?: boolean
  label?: string
}

const props = withDefaults(defineProps<Props>(), { disabled: false, label: '' })
const modelValue = defineModel<boolean>({ required: true })

const emit = defineEmits<{ change: [value: boolean] }>()

const isChecked = computed({
  get: () => modelValue.value,
  set: (val: boolean) => {
    if (props.disabled) return
    modelValue.value = val
    emit('change', val)
  },
})

const toggle = () => {
  isChecked.value = !isChecked.value
}
</script>

<template>
  <div class="win-switch-container" :class="{ 'is-disabled': disabled }" @click="toggle">
    <div class="win-switch" :class="{ 'is-checked': isChecked }">
      <div class="win-switch-thumb" />
    </div>
    <span v-if="label" class="win-switch-label">{{ label }}</span>
  </div>
</template>

<style scoped>
.win-switch-container {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
  transition: opacity 0.2s ease;
}

.win-switch-container.is-disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

/* 轨道样式 */
.win-switch {
  position: relative;
  width: 40px;
  height: 20px;
  border-radius: 10px;
  background-color: transparent;
  border: 1px solid var(--color-border);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

/* 鼠标悬停轨道 */
.win-switch-container:not(.is-disabled):hover .win-switch {
  border-color: var(--color-border-hover);
}

/* 选中状态轨道 */
.win-switch.is-checked {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}

.win-switch-container:not(.is-disabled):hover .win-switch.is-checked {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

/* 滑块样式 */
.win-switch-thumb {
  position: absolute;
  top: 50%;
  left: 3px;
  transform: translateY(-50%);
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: var(--color-text-secondary);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 选中状态滑块 */
.win-switch.is-checked .win-switch-thumb {
  left: calc(100% - 15px);
  background-color: #ffffff;
  width: 12px;
  height: 12px;
}

/* 点击时的拉伸效果 (Windows 11 特色) */
.win-switch-container:not(.is-disabled):active .win-switch-thumb {
  width: 16px;
  border-radius: 8px;
}

.win-switch-container:not(.is-disabled):active .win-switch.is-checked .win-switch-thumb {
  left: calc(100% - 19px);
}

.win-switch-label {
  font-size: 14px;
  color: var(--color-text);
}
</style>
