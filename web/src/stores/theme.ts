import { computed } from 'vue'
import { useColorMode } from '@vueuse/core'
import { z } from 'zod'
import { defineStore } from 'pinia'
import { useValidatedLocalStorage } from '@/composables/useValidatedLocalStorage'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'vue-go-session-theme-mode'
const themeModeSchema = z.enum(['light', 'dark', 'auto'])

export const useThemeStore = defineStore('theme', () => {
  const colorMode = useColorMode({
    selector: 'html',
    attribute: 'data-theme',
    storageKey: STORAGE_KEY,
    modes: { light: 'light', dark: 'dark' },
  })

  const mode = useValidatedLocalStorage<ThemeMode>(STORAGE_KEY, themeModeSchema, 'auto')

  const isDark = computed(() => {
    return colorMode.value === 'dark'
  })

  function syncColorMode(nextMode: ThemeMode): void {
    if (nextMode === 'auto') {
      colorMode.value = 'auto'
      return
    }
    colorMode.value = nextMode
  }

  function setMode(newMode: ThemeMode) {
    mode.value = newMode
    syncColorMode(newMode)
  }

  function cycleMode() {
    const modes: ThemeMode[] = ['light', 'dark', 'auto']
    const currentIndex = modes.indexOf(mode.value)
    const nextIndex = (currentIndex + 1) % modes.length
    const nextMode = modes[nextIndex]
    if (nextMode) {
      setMode(nextMode)
    }
  }

  function init() {
    syncColorMode(mode.value)
  }

  return { mode, isDark, setMode, cycleMode, init }
})
