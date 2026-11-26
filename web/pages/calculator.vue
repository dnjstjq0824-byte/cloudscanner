<template>
  <div class="max-w-3xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('calc.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('calc.lede') }}</p>
    </div>

    <div class="card p-6 grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="space-y-1.5">
        <label for="calc-plan" class="label">{{ $t('calc.plan') }}</label>
        <select id="calc-plan" v-model="planId" class="input">
          <option v-for="p in plans" :key="p.id" :value="p.id">{{ p.name }} — {{ money(p.monthly, 'USD') }}{{ $t('money.perMonth') }}</option>
        </select>
      </div>
      <div class="space-y-1.5">
        <label for="calc-seats" class="label">{{ $t('calc.seats') }}</label>
        <input id="calc-seats" v-model.number="seats" type="number" min="1" max="10000" class="input" />
      </div>
      <div class="space-y-1.5">
        <label for="calc-cycle" class="label">{{ $t('calc.cycle') }}</label>
        <select id="calc-cycle" v-model="cycle" class="input">
          <option value="month">{{ $t('calc.monthly') }}</option>
          <option value="year">{{ $t('calc.annual') }}</option>
        </select>
      </div>
    </div>

    <div class="card p-6 space-y-2">
      <div class="flex items-baseline justify-between">
        <span class="text-sm text-ink-secondary">{{ $t('calc.total') }}</span>
        <span class="text-2xl font-extrabold text-brand tabular-nums">
          <PriceDisplay :amount="total" currency="USD" :cycle="cycle" />
        </span>
      </div>
      <p class="text-[11px] text-ink-secondary">{{ $t('calc.note') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
const { money } = useMoney()
useHead(() => ({
  title: t('calc.title'),
  meta: [{ name: 'description', content: t('calc.meta') }],
}))

const plans = [
  { id: 'starter', name: 'Starter', monthly: 12 },
  { id: 'pro', name: 'Pro', monthly: 24 },
  { id: 'team', name: 'Team', monthly: 49 },
]
const planId = ref('pro')
const seats = ref(5)
const cycle = ref<'month' | 'year'>('month')

const total = computed(() => {
  const plan = plans.find((p) => p.id === planId.value) || plans[1]
  const safeSeats = Math.min(10000, Math.max(1, seats.value || 1))
  const monthly = plan.monthly * safeSeats
  return cycle.value === 'year' ? monthly * 12 * 0.8 : monthly
})
</script>
