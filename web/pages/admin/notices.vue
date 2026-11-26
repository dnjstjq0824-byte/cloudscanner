<template>
  <div class="space-y-8">
    <PageHeader :title="$t('admin.noticesTitle')" :subtitle="$t('admin.noticesSubtitle')" />

    <LoadingSkeleton v-if="loading" :rows="3" :label="$t('admin.loadingNotices')" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="space-y-3">
      <div v-for="n in notices" :key="n.id" class="card p-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div class="space-y-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <span class="text-sm font-bold text-ink">{{ n.title }}</span>
            <StatusBadge :status="n.severity" />
            <StatusBadge :status="n.is_active ? 'ACTIVE' : 'DISABLED'" />
          </div>
          <p class="text-xs text-ink-secondary">{{ n.body }}</p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <button class="btn-ghost !py-1.5 !px-3 !text-xs" @click="startEdit(n)">{{ $t('common.edit') }}</button>
          <button class="btn-danger !py-1.5 !px-3 !text-xs" :disabled="deleting === n.id" @click="askDelete(n)">
            {{ deleting === n.id ? $t('common.deleting') : $t('common.delete') }}
          </button>
        </div>
      </div>
      <EmptyState v-if="notices.length === 0" :title="$t('admin.noNoticesTitle')" :message="$t('admin.noNoticesBody')" />
    </div>

    <div class="card p-6 space-y-5">
      <h3 class="text-sm font-bold text-ink">{{ editing ? $t('admin.formEdit') : $t('admin.formPublish') }}</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5 md:col-span-2">
          <label for="notice-title" class="label">{{ $t('admin.fieldTitle') }}</label>
          <input id="notice-title" v-model="form.title" type="text" required :placeholder="$t('admin.titlePlaceholder')" class="input" />
        </div>
        <div class="space-y-1.5 md:col-span-2">
          <label for="notice-body" class="label">{{ $t('admin.fieldBody') }}</label>
          <textarea id="notice-body" v-model="form.body" rows="3" required :placeholder="$t('admin.bodyPlaceholder')" class="input" />
        </div>
        <div class="space-y-1.5">
          <label for="notice-severity" class="label">{{ $t('admin.fieldSeverity') }}</label>
          <select id="notice-severity" v-model="form.severity" class="input">
            <option>INFO</option>
            <option>WARNING</option>
            <option>CRITICAL</option>
          </select>
        </div>
        <div class="flex items-end pb-1">
          <label class="flex items-center gap-2 text-xs text-ink-secondary cursor-pointer">
            <input v-model="form.isActive" type="checkbox" class="rounded accent-sky-600" />
            {{ $t('admin.activeLabel') }}
          </label>
        </div>
      </div>
      <div v-if="formError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ formError }}
      </div>
      <div class="flex gap-2">
        <button class="btn-primary" :disabled="saving" @click="save">{{ saving ? $t('common.saving') : editing ? $t('admin.updateNotice') : $t('admin.publishNotice') }}</button>
        <button v-if="editing" class="btn-ghost" @click="cancelEdit">{{ $t('common.cancel') }}</button>
      </div>
    </div>

    <ConfirmModal
      :open="confirmOpen"
      :title="$t('admin.modalDeleteTitle')"
      :message="$t('admin.modalDeleteBody')"
      :confirm-label="$t('admin.modalDeleteConfirm')"
      danger
      :busy="deleting !== ''"
      @confirm="remove"
      @cancel="confirmOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { AdminNotice, ItemEnvelope, ListEnvelope } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('admin.noticesTitle') }))

const loading = ref(false)
const error = ref<ApiError | null>(null)
const notices = ref<AdminNotice[]>([])
const saving = ref(false)
const deleting = ref('')
const confirmOpen = ref(false)
const pendingDelete = ref('')
const editing = ref('')
const formError = ref('')

const form = reactive({ title: '', body: '', severity: 'INFO', isActive: true })

const { request } = useApi()
const toast = useToast()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<AdminNotice>>('/api/v1/admin/notices')
    notices.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.title = ''
  form.body = ''
  form.severity = 'INFO'
  form.isActive = true
  formError.value = ''
}

function startEdit(n: AdminNotice) {
  editing.value = n.id
  form.title = n.title
  form.body = n.body
  form.severity = n.severity
  form.isActive = n.is_active
  formError.value = ''
}

function cancelEdit() {
  editing.value = ''
  resetForm()
}

async function save() {
  formError.value = ''
  if (!form.title.trim() || !form.body.trim()) {
    formError.value = t('admin.requiredError')
    return
  }
  saving.value = true
  try {
    const body = { title: form.title.trim(), body: form.body.trim(), severity: form.severity, is_active: form.isActive }
    if (editing.value) {
      await request<ItemEnvelope<AdminNotice>>(`/api/v1/admin/notices/${editing.value}`, { method: 'PUT', body })
      toast.success(t('admin.toastNoticeUpdated'))
    } else {
      await request<ItemEnvelope<AdminNotice>>('/api/v1/admin/notices', { method: 'POST', body })
      toast.success(t('admin.toastNoticePublished'))
    }
    editing.value = ''
    resetForm()
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    formError.value = apiErr.message
  } finally {
    saving.value = false
  }
}

function askDelete(n: AdminNotice) {
  pendingDelete.value = n.id
  confirmOpen.value = true
}

async function remove() {
  const id = pendingDelete.value
  if (!id) {
    confirmOpen.value = false
    return
  }
  deleting.value = id
  try {
    await request(`/api/v1/admin/notices/${id}`, { method: 'DELETE' })
    toast.success(t('admin.toastNoticeDeleted'))
    confirmOpen.value = false
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    toast.error(apiErr.message)
    confirmOpen.value = false
  } finally {
    deleting.value = ''
    pendingDelete.value = ''
  }
}

onMounted(() => {
  load()
})
</script>
