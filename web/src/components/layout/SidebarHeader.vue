<template>
  <div class="bg-transparent px-5 pb-3 pt-5">
    <div class="flex items-center gap-3">
      <img
        :src="logoSvg"
        alt="logo"
        class="size-7 shrink-0"
      />
      <h2 class="m-0 text-xl font-semibold text-accent">web</h2>
    </div>
  </div>
</template>

<script setup lang="ts">
import logoSvg from '@/assets/logo.svg'
</script>
