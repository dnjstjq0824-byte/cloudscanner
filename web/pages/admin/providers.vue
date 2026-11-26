<template>
  <div class="space-y-8">
    <PageHeader :title="$t('admin.vendorsTitle')" :subtitle="$t('admin.vendorsSubtitle')" />

    <LoadingSkeleton v-if="loading" :rows="4" :label="$t('admin.loadingVendors')" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table-shell min-w-[820px]">
          <thead>
            <tr>
              <th class="py-3 px-4">{{ $t('admin.colVendor') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colPricing') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colMgmt') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colAffiliate') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colReseller') }}</th>
              <th class="py-3 px-4 text-right">{{ $t('admin.colSave') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in providers" :key="p.slug" class="hover:bg-elevated/60 transition">
              <td class="py-3 px-4">
                <div class="font-bold">{{ p.name }}</div>
                <div class="text-ink-secondary text-[11px]">{{ p.slug }} · {{ p.tier_priority }}</div>
              </td>
              <td v-for="field in statusFields" :key="field.key" class="py-3 px-4">
                <select
                  v-model="drafts[p.slug][field.key]"
                  class="input !w-auto !py-1 !text-[11px]"
                  :disabled="saving === p.slug"
                  :aria-label="`${p.name} ${field.label}`"
                >
                  <option v-for="s in statuses" :key="s" :value="s">{{ s }}</option>
                </select>
              </td>
              <td class="py-3 px-4 text-right">
                <button
                  class="btn-primary !py-1 !px-3 !text-[11px]"
                  :disabled="saving === p.slug || !dirty(p)"
                  @click="save(p)"
                >
                  {{ saving === p.slug ? $t('common.saving') : $t('common.save') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { ListEnvelope, Provider } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

useHead(() => ({ title: t('admin.vendorsTitle') }))

const statuses = ['VERIFIED', 'PENDING', 'FAILED', 'NOT_SUPPORTED', 'CONTRACT_REQUIRED', 'DISABLED']
const { t } = useI18n()
const statusFields = computed(() => [
  { key: 'pricing_api_status', label: t('admin.fieldPricing') },
  { key: 'management_api_status', label: t('admin.fieldMgmt') },
  { key: 'affiliate_status', label: t('admin.fieldAffiliate') },
  { key: 'reseller_program_status', label: t('admin.fieldReseller') },
] as const)

type Draft = Record<(typeof statusFields)['value'][number]['key'], string>

const loading = ref(false)
const error = ref<ApiError | null>(null)
const providers = ref<Provider[]>([])
const drafts = reactive<Record<string, Draft>>({})
const saving = ref('')

const { request } = useApi()
const toast = useToast()

function snapshot(p: Provider): Draft {
  return {
    pricing_api_status: p.pricing_api_status,
    management_api_status: p.management_api_status,
    affiliate_status: p.affiliate_status,
    reseller_program_status: p.reseller_program_status,
  }
}

function dirty(p: Provider): boolean {
  const d = drafts[p.slug]
  const s = snapshot(p)
  return (Object.keys(s) as (keyof Draft)[]).some((k) => d[k] !== s[k])
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<Provider>>('/api/v1/providers')
    providers.value = res.data ?? []
    for (const p of providers.value) {
      drafts[p.slug] = snapshot(p)
    }
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

async function save(p: Provider) {
  saving.value = p.slug
  try {
    const d = drafts[p.slug]
    const res = await request<{ data: Provider }>(`/api/v1/admin/providers/${p.slug}`, {
      method: 'PATCH',
      body: {
        pricing_api_status: d.pricing_api_status,
        management_api_status: d.management_api_status,
        affiliate_status: d.affiliate_status,
        reseller_program_status: d.reseller_program_status,
      },
    })
    const idx = providers.value.findIndex((x) => x.slug === p.slug)
    if (idx >= 0) {
      providers.value[idx] = res.data
      drafts[p.slug] = snapshot(res.data)
    }
    toast.success(t('admin.toastVendorSaved', { name: p.name }))
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    toast.error(apiErr.message)
  } finally {
    saving.value = ''
  }
}

onMounted(() => {
  load()
})
</script>
