<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Budgets & Alerts</h2>
      <p class="text-sm text-slate-400 mt-1">
        Hard caps are enforced before provisioning executes. A blocked approval can never bill you.
      </p>
    </div>

    <div v-if="loading" class="p-6 rounded-xl bg-slate-900 border border-slate-800 animate-pulse" aria-busy="true" aria-label="Loading budget">
      <div class="h-5 w-1/3 rounded bg-slate-800" />
      <div class="mt-4 h-10 w-full rounded bg-slate-800" />
    </div>

    <div v-else-if="error" class="p-8 rounded-xl bg-red-950/20 border border-red-900 text-center space-y-3">
      <h3 class="text-base font-bold text-white">{{ error.title }}</h3>
      <p class="text-xs text-slate-400">{{ error.message }}</p>
      <button class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold" @click="load">
        Retry
      </button>
    </div>

    <form v-else class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-5" @submit.prevent="save">
      <div v-if="!exists" class="px-4 py-3 rounded-lg bg-amber-950/40 border border-amber-900 text-xs text-amber-300">
        No budget policy configured yet. Provisioning is allowed until you set a cap.
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="budget-name" class="text-xs font-semibold text-slate-300">Policy name</label>
          <input
            id="budget-name"
            v-model="form.name"
            type="text"
            required
            placeholder="org-monthly-cap"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
        <div class="space-y-1.5">
          <label for="budget-limit" class="text-xs font-semibold text-slate-300">Monthly limit (USD)</label>
          <input
            id="budget-limit"
            v-model.number="form.monthlyLimit"
            type="number"
            required
            min="1"
            step="1"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
        <div class="space-y-1.5">
          <label for="budget-single" class="text-xs font-semibold text-slate-300">Max single server (USD/mo)</label>
          <input
            id="budget-single"
            v-model.number="form.singleMax"
            type="number"
            required
            min="1"
            step="1"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
        <div class="space-y-1.5">
          <label for="budget-alert" class="text-xs font-semibold text-slate-300">Alert threshold (%)</label>
          <input
            id="budget-alert"
            v-model.number="form.alertPct"
            type="number"
            required
            min="1"
            max="100"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
      </div>

      <div class="flex flex-col sm:flex-row gap-4 text-xs text-slate-300">
        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="form.hardBlock" type="checkbox" class="rounded text-sky-600 bg-slate-800 border-slate-700" />
          Hard-block provisioning when exceeded
        </label>
        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="form.isActive" type="checkbox" class="rounded text-sky-600 bg-slate-800 border-slate-700" />
          Policy active
        </label>
      </div>

      <div v-if="formError" class="px-4 py-3 rounded-lg bg-red-950/40 border border-red-900 text-xs text-red-300" role="alert">
        {{ formError }}
      </div>
      <div v-if="saved" class="px-4 py-3 rounded-lg bg-emerald-950/40 border border-emerald-900 text-xs text-emerald-300" role="status">
        Budget policy saved and enforced.
      </div>

      <button
        type="submit"
        :disabled="saving"
        class="px-5 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 disabled:opacity-60 text-white text-sm font-semibold transition"
      >
        {{ saving ? 'Saving…' : 'Save Budget Policy' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ItemEnvelope } from '~/types/api'
import type { BudgetPolicy } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Budgets & Alerts' })

const loading = ref(false)
const error = ref<ApiError | null>(null)
const exists = ref(false)
const saving = ref(false)
const saved = ref(false)
const formError = ref('')

const form = reactive({
  name: 'org-monthly-cap',
  monthlyLimit: 500,
  singleMax: 100,
  alertPct: 80,
  hardBlock: true,
  isActive: true,
})

const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ItemEnvelope<BudgetPolicy>>('/api/v1/cloud/budgets')
    const p = res.data
    exists.value = true
    form.name = p.name
    form.monthlyLimit = Number(p.monthly_limit_usd)
    form.singleMax = Number(p.max_single_server_monthly_usd)
    form.alertPct = p.alert_threshold_pct
    form.hardBlock = p.hard_block_provisioning_on_exceed
    form.isActive = p.is_active
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    if (apiErr.status === 404) {
      exists.value = false
    } else {
      error.value = apiErr
    }
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  saved.value = false
  formError.value = ''
  try {
    await request<ItemEnvelope<BudgetPolicy>>('/api/v1/cloud/budgets', {
      method: 'PUT',
      body: {
        name: form.name.trim(),
        monthly_limit_usd: form.monthlyLimit,
        max_single_server_monthly_usd: form.singleMax,
        alert_threshold_pct: form.alertPct,
        hard_block_provisioning_on_exceed: form.hardBlock,
        is_active: form.isActive,
      },
    })
    exists.value = true
    saved.value = true
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    formError.value = apiErr.message
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  load()
})
</script>
