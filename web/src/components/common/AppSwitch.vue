<script setup lang="ts">
import { computed, useId } from 'vue'

defineOptions({
  name: 'AppSwitch',
})

interface Props {
  disabled?: boolean
  id?: string
  label?: string
}

const { disabled = false, id, label = '' } = defineProps<Props>()

const modelValue = defineModel<boolean>({ required: true })

const emit = defineEmits<{
  change: [value: boolean]
}>()

const generatedId = useId()

const inputId = computed<string>(() => id || `app-switch-${generatedId}`)
const labelId = computed<string>(() => `${inputId.value}-label`)

const isChecked = computed<boolean>({
  get: () => modelValue.value,
  set: (value: boolean) => {
    if (disabled) return
    modelValue.value = value
    emit('change', value)
  },
})

function handleChange(event: Event): void {
  const target = event.target as HTMLInputElement
  isChecked.value = target.checked
}
</script>

<template>
  <div class="app-switch" :class="{ 'app-switch--disabled': disabled }">
    <input
      :id="inputId"
      class="app-switch__input"
      type="checkbox"
      role="switch"
      :aria-labelledby="label ? labelId : undefined"
      :checked="isChecked"
      :disabled="disabled"
      @change="handleChange"
    />

    <label class="app-switch__control" :for="inputId" aria-hidden="true">
      <span class="app-switch__thumb" />
    </label>

    <label v-if="label" :id="labelId" class="app-switch__label" :for="inputId">
      {{ label }}
    </label>
  </div>
</template>

<style scoped>
.app-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
  transition: opacity 0.2s ease;
}

.app-switch--disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.app-switch__input {
  position: absolute;
  inset: 0;
  inline-size: 100%;
  block-size: 100%;
  margin: 0;
  opacity: 0;
  pointer-events: none;
}

.app-switch__control {
  position: relative;
  display: inline-flex;
  align-items: center;
  inline-size: 40px;
  block-size: 20px;
  border-radius: 10px;
  background-color: transparent;
  border: 1px solid var(--sys-color-border);
  box-sizing: border-box;
  transition:
    background-color 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    border-color 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: inherit;
}

.app-switch:not(.app-switch--disabled):hover .app-switch__control {
  border-color: var(--sys-color-border-hover);
}

.app-switch__input:checked + .app-switch__control {
  background-color: var(--sys-color-accent);
  border-color: var(--sys-color-accent);
}

.app-switch:not(.app-switch--disabled):hover .app-switch__input:checked + .app-switch__control {
  background-color: var(--sys-color-accent-hover);
  border-color: var(--sys-color-accent-hover);
}

.app-switch__thumb {
  position: absolute;
  inset-block-start: 50%;
  inset-inline-start: 3px;
  transform: translateY(-50%);
  inline-size: 12px;
  block-size: 12px;
  border-radius: 50%;
  background-color: var(--sys-color-text-secondary);
  transition:
    inset-inline-start 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    inline-size 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    border-radius 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    background-color 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.app-switch__input:checked + .app-switch__control .app-switch__thumb {
  inset-inline-start: calc(100% - 15px);
  background-color: var(--sys-color-on-accent);
}

.app-switch:not(.app-switch--disabled):active .app-switch__thumb {
  inline-size: 16px;
  border-radius: 8px;
}

.app-switch:not(.app-switch--disabled):active .app-switch__input:checked + .app-switch__control .app-switch__thumb {
  inset-inline-start: calc(100% - 19px);
}

.app-switch__input:focus-visible + .app-switch__control {
  outline: 2px solid var(--sys-color-focus-ring);
  outline-offset: 2px;
}

.app-switch__label {
  font-size: 14px;
  color: var(--sys-color-text-primary);
  cursor: inherit;
}
</style>
