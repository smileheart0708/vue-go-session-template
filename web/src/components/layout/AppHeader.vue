<template>
  <header
    class="fixed inset-x-0 top-0 z-90 flex h-header items-center gap-4 border-b border-(--cmp-header-border) bg-(--cmp-header-bg) px-6 backdrop-blur-md backdrop-saturate-140 md:left-sidebar max-md:px-4"
  >
    <IconButton class="md:hidden" title="切换侧边栏" @click="emit('toggle-sidebar')">
      <Menu />
    </IconButton>

    <h1
      class="m-0 flex-1 min-w-0 overflow-hidden text-ellipsis whitespace-nowrap text-lg font-semibold text-text-primary"
    >
      {{ currentTitle }}
    </h1>

    <div class="flex items-center gap-2">
      <ThemeToggle
        v-model="mode"
        light-label="浅色"
        dark-label="深色"
        auto-label="自动"
        tooltip-prefix="当前："
        tooltip-suffix="（长按切换）"
        @change="handleThemeChange"
      />
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu } from 'lucide-vue-next'
import { IconButton, ThemeToggle } from '@/components/common'
import { useTheme } from '@/composables'
import type { ThemeMode } from '@/composables'

const emit = defineEmits<{ 'toggle-sidebar': [] }>()

const route = useRoute()
const { mode, setTheme } = useTheme()

const currentTitle = computed<string>(() => {
  const title = route.meta['title']
  return typeof title === 'string' ? title : ''
})

async function handleThemeChange(nextMode: ThemeMode, event?: MouseEvent): Promise<void> {
  await setTheme(nextMode, event)
}
</script>
