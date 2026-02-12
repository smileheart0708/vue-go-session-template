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
  flex: 1;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}
</style>
