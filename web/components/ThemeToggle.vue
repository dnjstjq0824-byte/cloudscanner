<template>
  <button
    class="icon-btn w-9 h-9 text-ink-secondary hover:text-ink"
    :aria-label="isDark ? 'Switch to light theme' : 'Switch to dark theme'"
    :title="isDark ? 'Light mode' : 'Dark mode'"
    @click="toggle"
  >
    <span v-if="isDark" aria-hidden="true">☀</span>
    <span v-else aria-hidden="true">☾</span>
  </button>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from '~/composables/useTheme'

const { isDark, toggle, init } = useTheme()
onMounted(() => init())
</script>
