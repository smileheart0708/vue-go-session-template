<script setup lang="ts">
import { computed, useId } from 'vue'

defineOptions({ name: 'AppSwitch' })

interface Props {
  disabled?: boolean
  id?: string
  label?: string
}

const { disabled = false, id, label = '' } = defineProps<Props>()

const modelValue = defineModel<boolean>({ required: true })

const emit = defineEmits<{ change: [value: boolean] }>()

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
  if (!(event.target instanceof HTMLInputElement)) return
  isChecked.value = event.target.checked
}
</script>

<template>
  <div
    class="relative inline-flex select-none items-center gap-3 transition-opacity"
    :class="disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'"
  >
    <input
      :id="inputId"
      class="peer sr-only"
      type="checkbox"
      role="switch"
      :aria-labelledby="label ? labelId : undefined"
      :checked="isChecked"
      :disabled="disabled"
      @change="handleChange"
    />

    <label
      class="relative inline-flex h-5 w-10 items-center rounded-full border border-border transition-colors duration-200 ease-out peer-checked:border-accent peer-checked:bg-accent peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-(--sys-color-focus-ring)"
      :for="inputId"
      aria-hidden="true"
    >
      <span
        class="absolute left-0.75 top-1/2 h-3 w-3 -translate-y-1/2 rounded-full bg-text-secondary transition-all duration-200 ease-out peer-checked:left-[calc(100%-15px)] peer-checked:bg-on-accent"
      />
    </label>

    <label
      v-if="label"
      :id="labelId"
      class="cursor-inherit text-sm text-text-primary"
      :for="inputId"
    >
      {{ label }}
    </label>
  </div>
</template>
