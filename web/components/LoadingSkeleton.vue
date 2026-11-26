<template>
  <div class="space-y-3" aria-busy="true" :aria-label="ariaText">
    <div v-for="n in rows" :key="n" class="card p-5">
      <div class="shimmer h-4 w-1/2 rounded" />
      <div v-if="lines > 1" class="shimmer h-3 w-2/3 rounded mt-3" />
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const props = withDefaults(defineProps<{ rows?: number; lines?: number; label?: string }>(), {
  rows: 3,
  lines: 1,
  label: undefined,
})
const ariaText = computed(() => props.label || t('common.loading'))
</script>
