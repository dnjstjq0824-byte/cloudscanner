<template>
  <div class="space-y-8">
    <PageHeader :title="$t('admin.settingsTitle')" :subtitle="$t('admin.settingsSubtitle')" />

    <LoadingSkeleton v-if="loading" :rows="1" :lines="2" :label="$t('admin.loadingSettings')" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <form v-else class="card p-6 space-y-5" @submit.prevent="save">
      <div class="p-4 rounded-xl border" :class="form.maintenance ? 'bg-amber-500/10 border-amber-500/30' : 'bg-elevated border-subtle'">
        <label class="flex items-start gap-3 cursor-pointer">
          <input v-model="form.maintenance" type="checkbox" class="mt-1 rounded accent-sky-600" />
          <span>
            <span class="block text-sm font-bold text-ink">{{ $t('admin.maintenanceTitle') }}</span>
            <span class="block text-xs text-ink-secondary mt-0.5">
              {{ $t('admin.maintenanceBody') }}
            </span>
          </span>
        </label>
      </div>

      <div class="space-y-1.5">
        <label for="settings-announcement" class="label">{{ $t('admin.announcementLabel') }}</label>
        <input id="settings-announcement" v-model="form.announcement" type="text" :placeholder="$t('admin.announcementPlaceholder')" class="input" />
      </div>

      <p v-if="updatedAt" class="text-[11px] text-ink-secondary">{{ $t('admin.updatedAt', { date: updatedAt }) }}</p>

      <div v-if="formError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ formError }}
      </div>

      <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? $t('common.saving') : $t('admin.saveSettings') }}</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { ItemEnvelope, SiteSettings } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t, locale } = useI18n()
useHead(() => ({ title: t('admin.settingsTitle') }))

const loading = ref(false)
const error = ref<ApiError | null>(null)
const saving = ref(false)
const formError = ref('')
const updatedAt = ref('')

const form = reactive({ maintenance: false, announcement: '' })

const { request } = useApi()
const toast = useToast()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ItemEnvelope<SiteSettings>>('/api/v1/admin/settings')
    form.maintenance = res.data.maintenance_mode
    form.announcement = res.data.announcement
    try {
      updatedAt.value = new Date(res.data.updated_at).toLocaleString(locale.value)
    } catch {
      updatedAt.value = res.data.updated_at
    }
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  formError.value = ''
  try {
    const res = await request<ItemEnvelope<SiteSettings>>('/api/v1/admin/settings', {
      method: 'PUT',
      body: { maintenance_mode: form.maintenance, announcement: form.announcement },
    })
    form.maintenance = res.data.maintenance_mode
    form.announcement = res.data.announcement
    toast.success(t('admin.toastSettingsSaved'))
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
