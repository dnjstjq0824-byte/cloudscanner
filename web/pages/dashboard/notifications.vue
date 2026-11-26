<template>
  <div class="space-y-8">
    <PageHeader :title="$t('notifcenter.h1')" :subtitle="$t('notifcenter.lede')" />

    <div class="space-y-3">
      <div
        v-for="n in items"
        :key="n.id"
        class="card p-4 flex items-start gap-3"
        :class="n.read ? 'opacity-70' : ''"
      >
        <span class="mt-1.5 w-2 h-2 rounded-full shrink-0" :class="n.read ? 'bg-subtle' : 'bg-brand'" />
        <div class="min-w-0 flex-1">
          <div class="text-sm font-bold text-ink">{{ n.title }}</div>
          <div class="text-xs text-ink-secondary mt-0.5">{{ n.body }}</div>
          <div class="text-[11px] text-ink-secondary mt-1 tabular-nums">{{ n.at }}</div>
        </div>
        <button v-if="!n.read" class="text-xs font-semibold text-brand hover:underline shrink-0" @click="markRead(n.id)">
          {{ $t('notifcenter.markRead') }}
        </button>
      </div>
      <EmptyState
        v-if="items.length === 0"
        :title="$t('notifcenter.emptyTitle')"
        :message="$t('notifcenter.emptyBody')"
      />
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('notifcenter.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockUserNotices } from '~/utils/mock-data'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t } = useI18n()
useHead(() => ({ title: t('notifcenter.title') }))

// TODO(api): replace with user notification endpoints when implemented.
const items = ref(mockUserNotices())

function markRead(id: string) {
  const found = items.value.find((n) => n.id === id)
  if (found) {
    found.read = true
  }
}
</script>
