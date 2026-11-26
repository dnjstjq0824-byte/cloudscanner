<template>
  <div class="card p-8 text-center space-y-3" role="alert">
    <h3 class="text-base font-bold text-ink">{{ title }}</h3>
    <p class="text-xs text-ink-secondary">{{ message }}</p>
    <button class="btn-primary !py-2 !text-xs" @click="$emit('retry')">{{ $t('common.retry') }}</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; message: string }>()
defineEmits<{ (e: 'retry'): void }>()
</script>
