<template>
  <div class="space-y-8">
    <PageHeader :title="$t('profile.h1')" :subtitle="$t('profile.lede')" />

    <LoadingSkeleton v-if="loading" :rows="2" :label="$t('profile.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <template v-else>
      <div class="card p-6">
        <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('auth.fullName') }}</dt>
            <dd class="font-semibold text-ink">{{ me?.full_name || '—' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('auth.email') }}</dt>
            <dd class="font-semibold text-ink">{{ me?.email || '—' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('admin.colStatus') }}</dt>
            <dd><StatusBadge :status="me?.status || '—'" /></dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('admin.colJoined') }}</dt>
            <dd class="font-semibold text-ink">{{ joined }}</dd>
          </div>
        </dl>
      </div>

      <div class="card p-6 space-y-4">
        <h2 class="text-sm font-bold text-ink">{{ $t('profile.prefsTitle') }}</h2>
        <p class="text-xs text-ink-secondary">
          {{ $t('profile.prefsNote') }}
          <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
        </p>
        <label class="flex items-center gap-2 text-sm text-ink-secondary cursor-pointer">
          <input v-model="prefs.dealAlerts" type="checkbox" class="rounded accent-sky-600" />
          {{ $t('profile.prefAlerts') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-ink-secondary cursor-pointer">
          <input v-model="prefs.weeklyDigest" type="checkbox" class="rounded accent-sky-600" />
          {{ $t('profile.prefDigest') }}
        </label>
        <div>
          <button class="btn-primary !py-2 !text-xs" @click="save">{{ $t('common.save') }}</button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { AuthUser } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t, locale } = useI18n()
useHead(() => ({ title: t('profile.title') }))

const me = ref<AuthUser | null>(null)
const joined = ref('—')
const loading = ref(false)
const error = ref<ApiError | null>(null)
const prefs = reactive({ dealAlerts: true, weeklyDigest: false })
const { request } = useApi()
const toast = useToast()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<{ data: AuthUser }>('/api/v1/auth/me')
    me.value = res.data
    try {
      joined.value = new Date(res.data.created_at).toLocaleDateString(locale.value)
    } catch {
      joined.value = res.data.created_at
    }
    try {
      const saved = localStorage.getItem('tirp-prefs')
      if (saved) {
        Object.assign(prefs, JSON.parse(saved))
      }
    } catch {
      // Ignore corrupt local preferences.
    }
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function save() {
  // TODO(api): replace with PATCH /api/v1/users/me/preferences when implemented.
  localStorage.setItem('tirp-prefs', JSON.stringify(prefs))
  toast.success(t('profile.toastSaved'))
}

onMounted(load)
</script>
