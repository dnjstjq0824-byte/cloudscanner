<template>
  <div class="card p-12 text-center space-y-3">
    <h4 class="text-base font-bold text-ink">{{ title }}</h4>
    <p class="text-xs text-ink-secondary max-w-md mx-auto">{{ message }}</p>
    <div>
      <NuxtLink v-if="ctaTo" :to="ctaTo" class="btn-primary !py-2 !text-xs">{{ ctaLabel }}</NuxtLink>
      <slot v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; message: string; ctaTo?: string; ctaLabel?: string }>()
</script>
