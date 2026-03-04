import { watch } from 'vue'
import { useLocalStorage, type RemovableRef } from '@vueuse/core'
import type { ZodType } from 'zod'
import { parseAndCheckChanged } from '@/types/zod'

export function useValidatedLocalStorage<T>(
  key: string,
  schema: ZodType<T>,
  fallback: T,
): RemovableRef<T> {
  const state = useLocalStorage<T>(key, fallback)

  const initial = parseAndCheckChanged(schema, state.value, fallback)
  if (initial.changed) {
    state.value = initial.value
  }

  watch(
    state,
    (nextValue) => {
      const parsed = parseAndCheckChanged(schema, nextValue, fallback)
      if (!parsed.changed) {
        return
      }
      state.value = parsed.value
    },
    { deep: true },
  )

  return state
}
