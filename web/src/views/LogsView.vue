<template>
  <div class="logs">
    <LogBoard />
  </div>
</template>

<script setup lang="ts">
import LogBoard from '@/components/logs/LogBoard.vue'
</script>

<style scoped>
.logs {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
