<template>
  <div class="space-y-8">
    <PageHeader :title="$t('alerts.h1')" :subtitle="$t('alerts.lede')" />

    <div class="card p-6">
      <form class="grid grid-cols-1 sm:grid-cols-4 gap-3" @submit.prevent="add">
        <input v-model="label" type="text" required :placeholder="$t('alerts.labelPlaceholder')" class="input sm:col-span-2" :aria-label="$t('alerts.label')" />
        <input v-model.number="target" type="number" min="1" step="0.01" required class="input" :aria-label="$t('alerts.target')" />
        <button type="submit" class="btn-primary">{{ $t('alerts.add') }}</button>
      </form>
      <p class="text-[11px] text-ink-secondary mt-3">{{ $t('alerts.localNote') }}</p>
    </div>

    <div class="space-y-3">
      <div v-for="alert in alerts" :key="alert.id" class="card p-4 flex items-center gap-3">
        <button
          class="w-10 h-6 rounded-full transition shrink-0"
          :class="alert.active ? 'bg-brand' : 'bg-subtle'"
          role="switch"
          :aria-checked="alert.active"
          :aria-label="alert.label"
          @click="toggle(alert.id)"
        >
          <span class="block w-4 h-4 m-1 rounded-full bg-white transition" :class="alert.active ? 'translate-x-4' : ''" />
        </button>
        <div class="min-w-0 flex-1">
          <div class="text-sm font-bold text-ink truncate">{{ alert.label }}</div>
          <div class="text-xs text-ink-secondary tabular-nums">
            <PriceDisplay :amount="alert.targetPrice" :currency="alert.currency" />
          </div>
        </div>
        <button class="text-xs font-semibold text-rose-500 hover:underline shrink-0" @click="remove(alert.id)">
          {{ $t('common.delete') }}
        </button>
      </div>
      <EmptyState
        v-if="alerts.length === 0"
        :title="$t('alerts.emptyTitle')"
        :message="$t('alerts.emptyBody')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { mockPriceAlerts, type MockPriceAlert } from '~/utils/mock-data'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t } = useI18n()
useHead(() => ({ title: t('alerts.title') }))

const alerts = ref<MockPriceAlert[]>([])
const label = ref('')
const target = ref<number | null>(null)

// TODO(api): replace localStorage with alert CRUD endpoints when implemented.
function persist() {
  localStorage.setItem('tirp-alerts', JSON.stringify(alerts.value))
}

function add() {
  if (!label.value.trim() || !target.value || target.value <= 0) {
    return
  }
  alerts.value.unshift({
    id: `al-${Date.now()}`,
    label: label.value.trim(),
    targetPrice: target.value,
    currency: 'USD',
    active: true,
  })
  label.value = ''
  target.value = null
  persist()
}

function toggle(id: string) {
  const found = alerts.value.find((a) => a.id === id)
  if (found) {
    found.active = !found.active
    persist()
  }
}

function remove(id: string) {
  alerts.value = alerts.value.filter((a) => a.id !== id)
  persist()
}

onMounted(() => {
  try {
    const saved = localStorage.getItem('tirp-alerts')
    alerts.value = saved ? (JSON.parse(saved) as MockPriceAlert[]) : mockPriceAlerts()
  } catch {
    alerts.value = mockPriceAlerts()
  }
})
</script>
