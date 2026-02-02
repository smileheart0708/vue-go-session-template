import { ref, computed, watch } from 'vue'
import { defineStore } from 'pinia'
import { useColorMode } from '@vueuse/core'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'vue-go-session-theme-mode'

export const useThemeStore = defineStore('theme', () => {
  const colorMode = useColorMode({
    selector: 'html',
    attribute: 'data-theme',
    storageKey: STORAGE_KEY,
    modes: {
      light: 'light',
      dark: 'dark',
    },
  })

  const mode = ref<ThemeMode>('auto')

  const isDark = computed(() => {
    if (mode.value === 'auto') {
      return colorMode.value === 'dark'
    }
    return colorMode.value === 'dark'
  })

  function setMode(newMode: ThemeMode) {
    mode.value = newMode
    if (newMode === 'auto') {
      colorMode.value = 'auto'
    } else {
      colorMode.value = newMode
    }
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
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'light' || stored === 'dark' || stored === 'auto') {
      mode.value = stored
    }

    watch(mode, (newMode) => {
      localStorage.setItem(STORAGE_KEY, newMode)
      if (newMode === 'auto') {
        colorMode.value = 'auto'
      } else {
        colorMode.value = newMode
      }
    })
  }

  return {
    mode,
    isDark,
    setMode,
    cycleMode,
    init,
  }
})