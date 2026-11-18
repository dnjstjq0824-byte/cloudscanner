<template>
  <div class="space-y-8">
    <PageHeader
      title="Budgets & Alerts"
      subtitle="Hard caps are enforced before provisioning executes. A blocked approval can never bill you."
    />

    <LoadingSkeleton v-if="loading" :rows="1" :lines="2" label="Loading budget" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <form v-else class="card p-6 space-y-5" @submit.prevent="save">
      <div v-if="!exists" class="px-4 py-3 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs text-amber-700 dark:text-amber-300">
        No budget policy configured yet. Provisioning is allowed until you set a cap.
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="budget-name" class="label">Policy name</label>
          <input id="budget-name" v-model="form.name" type="text" required placeholder="org-monthly-cap" class="input" />
        </div>
        <div class="space-y-1.5">
          <label for="budget-limit" class="label">Monthly limit (USD)</label>
          <input id="budget-limit" v-model.number="form.monthlyLimit" type="number" required min="1" step="1" class="input" />
        </div>
        <div class="space-y-1.5">
          <label for="budget-single" class="label">Max single server (USD/mo)</label>
          <input id="budget-single" v-model.number="form.singleMax" type="number" required min="1" step="1" class="input" />
        </div>
        <div class="space-y-1.5">
          <label for="budget-alert" class="label">Alert threshold (%)</label>
          <input id="budget-alert" v-model.number="form.alertPct" type="number" required min="1" max="100" class="input" />
        </div>
      </div>

      <div class="flex flex-col sm:flex-row gap-4 text-xs text-ink-secondary">
        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="form.hardBlock" type="checkbox" class="rounded accent-sky-600" />
          Hard-block provisioning when exceeded
        </label>
        <label class="flex items-center gap-2 cursor-pointer">
          <input v-model="form.isActive" type="checkbox" class="rounded accent-sky-600" />
          Policy active
        </label>
      </div>

      <div v-if="formError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ formError }}
      </div>

      <button type="submit" :disabled="saving" class="btn-primary">
        {{ saving ? 'Saving…' : 'Save Budget Policy' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
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
const toast = useToast()

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
    toast.success('Budget policy saved and enforced.')
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
