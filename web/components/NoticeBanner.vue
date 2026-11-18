<template>
  <div v-if="banners.length > 0" class="space-y-2">
    <div
      v-for="n in banners"
      :key="n.id"
      class="rounded-xl border px-4 py-3 text-xs flex items-start gap-3"
      :class="tone(n.severity)"
      role="status"
    >
      <span class="font-bold shrink-0 mt-px">{{ n.severity }}</span>
      <div class="min-w-0">
        <span class="font-semibold">{{ n.title }}</span>
        <span class="opacity-80"> — {{ n.body }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { ListEnvelope } from '~/types/api'

interface Notice {
  id: string
  title: string
  body: string
  severity: string
}

const banners = ref<Notice[]>([])

function tone(severity: string): string {
  if (severity === 'CRITICAL') return 'bg-rose-500/10 border-rose-500/30 text-rose-700 dark:text-rose-200'
  if (severity === 'WARNING') return 'bg-amber-500/10 border-amber-500/30 text-amber-800 dark:text-amber-200'
  return 'bg-sky-500/10 border-sky-500/30 text-sky-800 dark:text-sky-200'
}

onMounted(async () => {
  try {
    const { request } = useApi()
    const res = await request<ListEnvelope<Notice>>('/api/v1/notices')
    banners.value = (res.data ?? []).slice(0, 3)
  } catch {
    banners.value = []
  }
})
</script>
