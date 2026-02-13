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
  inline-size: 1px;
  block-size: 1px;
  overflow: hidden;
  clip-path: inset(50%);
  white-space: nowrap;
}

.app-switch__control {
  position: relative;
  display: inline-flex;
  align-items: center;
  inline-size: 40px;
  block-size: 20px;
  border-radius: 10px;
  background-color: transparent;
  border: 1px solid var(--color-border);
  box-sizing: border-box;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: inherit;
}

.app-switch:not(.app-switch--disabled):hover .app-switch__control {
  border-color: var(--color-border-hover);
}

.app-switch__input:checked + .app-switch__control {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}

.app-switch:not(.app-switch--disabled):hover .app-switch__input:checked + .app-switch__control {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.app-switch__thumb {
  position: absolute;
  inset-block-start: 50%;
  inset-inline-start: 3px;
  transform: translateY(-50%);
  inline-size: 12px;
  block-size: 12px;
  border-radius: 50%;
  background-color: var(--color-text-secondary);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.app-switch__input:checked + .app-switch__control .app-switch__thumb {
  inset-inline-start: calc(100% - 15px);
  background-color: var(--color-on-primary);
}

.app-switch:not(.app-switch--disabled):active .app-switch__thumb {
  inline-size: 16px;
  border-radius: 8px;
}

.app-switch:not(.app-switch--disabled):active .app-switch__input:checked + .app-switch__control .app-switch__thumb {
  inset-inline-start: calc(100% - 19px);
}

.app-switch__input:focus-visible + .app-switch__control {
  outline: 2px solid var(--color-focus-ring);
  outline-offset: 2px;
}

.app-switch__label {
  font-size: 14px;
  color: var(--color-text);
  cursor: inherit;
}
</style>
