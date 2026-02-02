/**
 * Theme composable - wugou style implementation
 * Uses Pinia store with View Transitions API
 */

import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

export type ThemeMode = 'light' | 'dark' | 'auto'

export function useTheme() {
  const themeStore = useThemeStore()

  const mode = computed(() => themeStore.mode)
  const isDark = computed(() => themeStore.isDark)

  async function setTheme(newMode: ThemeMode, event?: MouseEvent) {
    // Fallback: View Transitions API not supported
    if (!document.startViewTransition) {
      themeStore.setMode(newMode)
      return
    }

    // Get animation start coordinates
    const x = event?.clientX ?? window.innerWidth - 100
    const y = event?.clientY ?? 100

    const endRadius = Math.hypot(
      Math.max(x, window.innerWidth - x),
      Math.max(y, window.innerHeight - y),
    )

    // Use View Transitions API
    const transition = document.startViewTransition(() => {
      themeStore.setMode(newMode)
    })

    // Circular expansion animation
    await transition.ready
    document.documentElement.animate(
      {
        clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`],
      },
      {
        duration: 400,
        easing: 'ease-out',
        pseudoElement: '::view-transition-new(root)',
      },
    )
  }

  async function toggleTheme(event?: MouseEvent) {
    const modes: ThemeMode[] = ['light', 'dark', 'auto']
    const currentIndex = modes.indexOf(mode.value)
    const nextIndex = (currentIndex + 1) % modes.length
    const nextMode = modes[nextIndex]
    if (nextMode) {
      await setTheme(nextMode, event)
    }
  }

  return {
    mode,
    isDark,
    setTheme,
    toggleTheme,
  }
}
