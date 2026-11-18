<template>
  <div class="space-y-8">
    <PageHeader
      title="Cloud Connections"
      subtitle="Connect your own provider accounts (BYOC). API tokens are sealed with KMS envelope encryption and never displayed again."
    />

    <LoadingSkeleton v-if="loading" :rows="3" label="Loading connections" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="space-y-3">
      <div
        v-for="conn in connections"
        :key="conn.id"
        class="card p-5 flex flex-col md:flex-row md:items-center justify-between gap-4"
      >
        <div class="space-y-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <span class="text-base font-bold text-ink">{{ conn.connection_name }}</span>
            <span class="px-2 py-0.5 rounded text-[10px] font-semibold uppercase bg-elevated text-brand border border-subtle">
              {{ conn.provider_slug }}
            </span>
            <StatusBadge :status="conn.connection_status" />
            <StatusBadge :status="conn.access_level" />
          </div>
          <p class="text-xs text-ink-secondary">
            {{ conn.auth_mode }} · {{ conn.external_account_email }}
          </p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <button
            class="btn-danger !py-1.5 !px-3 !text-xs"
            :disabled="disconnecting === conn.id"
            @click="askDisconnect(conn)"
          >
            {{ disconnecting === conn.id ? 'Disconnecting…' : 'Disconnect' }}
          </button>
        </div>
      </div>

      <EmptyState
        v-if="connections.length === 0"
        title="No cloud accounts connected yet"
        message="Add your first connection below to start managing servers."
      />
    </div>

    <div class="card p-6 space-y-5">
      <h3 class="text-sm font-bold text-ink">Connect a new cloud account</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="conn-provider" class="label">Provider</label>
          <select id="conn-provider" v-model="form.providerSlug" class="input">
            <option value="digitalocean">DigitalOcean</option>
            <option value="hetzner">Hetzner Cloud</option>
            <option value="vultr">Vultr</option>
            <option value="upcloud">UpCloud</option>
            <option value="linode">Linode (Akamai)</option>
          </select>
        </div>
        <div class="space-y-1.5">
          <label for="conn-name" class="label">Connection name</label>
          <input id="conn-name" v-model="form.name" type="text" placeholder="Production DO Account" class="input" />
        </div>
        <div class="space-y-1.5">
          <label for="conn-access" class="label">Access level</label>
          <select id="conn-access" v-model="form.accessLevel" class="input">
            <option value="READ_ONLY">Read only (inventory + cost)</option>
            <option value="READ_WRITE">Read + write (provisioning)</option>
          </select>
        </div>
        <div class="space-y-1.5">
          <label for="conn-token" class="label">Scoped API token</label>
          <input id="conn-token" v-model="form.token" type="password" autocomplete="off" placeholder="Paste a scoped token from your provider" class="input" />
        </div>
      </div>

      <div v-if="formError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ formError }}
      </div>

      <button class="btn-primary" :disabled="connecting" @click="connect">
        {{ connecting ? 'Verifying & sealing token…' : 'Connect Account' }}
      </button>
    </div>

    <ConfirmModal
      :open="confirmOpen"
      title="Disconnect this cloud account?"
      message="The stored credential is revoked immediately and synced servers stop updating. This cannot be undone."
      confirm-label="Disconnect"
      danger
      :busy="disconnecting !== ''"
      @confirm="disconnect"
      @cancel="confirmOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { CloudConnection, ItemEnvelope, ListEnvelope } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Cloud Connections' })

const loading = ref(false)
const error = ref<ApiError | null>(null)
const connections = ref<CloudConnection[]>([])
const connecting = ref(false)
const disconnecting = ref('')
const confirmOpen = ref(false)
const pendingDisconnect = ref('')
const formError = ref('')

const form = reactive({
  providerSlug: 'digitalocean',
  name: '',
  accessLevel: 'READ_ONLY',
  token: '',
})

const { request } = useApi()
const toast = useToast()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<CloudConnection>>('/api/v1/cloud/connections')
    connections.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

async function connect() {
  formError.value = ''
  if (!form.token.trim()) {
    formError.value = 'An API token is required.'
    return
  }
  connecting.value = true
  try {
    await request<ItemEnvelope<CloudConnection>>('/api/v1/cloud/connections', {
      method: 'POST',
      body: {
        provider_slug: form.providerSlug,
        provider_id: form.providerSlug,
        connection_name: form.name.trim() || `${form.providerSlug}-default`,
        auth_mode: 'SCOPED_API_TOKEN',
        access_level: form.accessLevel,
        raw_api_token: form.token.trim(),
        granted_scopes: form.accessLevel === 'READ_WRITE' ? ['read', 'write'] : ['read'],
      },
    })
    form.name = ''
    form.token = ''
    toast.success('Cloud account connected and token sealed.')
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    formError.value = apiErr.message
  } finally {
    connecting.value = false
  }
}

function askDisconnect(conn: CloudConnection) {
  pendingDisconnect.value = conn.id
  confirmOpen.value = true
}

async function disconnect() {
  const id = pendingDisconnect.value
  if (!id) {
    confirmOpen.value = false
    return
  }
  disconnecting.value = id
  try {
    await request(`/api/v1/cloud/connections/${id}`, { method: 'DELETE' })
    toast.success('Connection disconnected and credential revoked.')
    confirmOpen.value = false
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr
    confirmOpen.value = false
  } finally {
    disconnecting.value = ''
    pendingDisconnect.value = ''
  }
}

onMounted(() => {
  load()
})
</script>
