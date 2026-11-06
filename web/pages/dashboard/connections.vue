<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Cloud Connections</h2>
      <p class="text-sm text-slate-400 mt-1">
        Connect your own provider accounts (BYOC). API tokens are sealed with KMS envelope encryption and never displayed again.
      </p>
    </div>

    <div v-if="loading" class="space-y-3" aria-busy="true" aria-label="Loading connections">
      <div v-for="n in 3" :key="n" class="p-5 rounded-xl bg-slate-900 border border-slate-800 animate-pulse">
        <div class="h-4 w-1/3 rounded bg-slate-800" />
      </div>
    </div>

    <div v-else-if="error" class="p-8 rounded-xl bg-red-950/20 border border-red-900 text-center space-y-3">
      <h3 class="text-base font-bold text-white">{{ error.title }}</h3>
      <p class="text-xs text-slate-400">{{ error.message }}</p>
      <button class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold" @click="load">
        Retry
      </button>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="conn in connections"
        :key="conn.id"
        class="p-5 rounded-xl bg-slate-900 border border-slate-800 flex flex-col md:flex-row justify-between gap-4"
      >
        <div class="space-y-1">
          <div class="flex items-center gap-2 flex-wrap">
            <span class="text-base font-bold text-white">{{ conn.connection_name }}</span>
            <span class="px-2 py-0.5 rounded text-[10px] font-semibold uppercase bg-slate-800 text-sky-400 border border-slate-700">
              {{ conn.provider_slug }}
            </span>
            <span
              class="px-2 py-0.5 rounded-full text-[10px] font-semibold border"
              :class="conn.connection_status === 'ACTIVE' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-400 border-slate-700'"
            >
              {{ conn.connection_status }}
            </span>
          </div>
          <p class="text-xs text-slate-400">
            {{ conn.auth_mode }} · {{ conn.access_level }} · {{ conn.external_account_email }}
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="px-3 py-1.5 rounded bg-red-950 hover:bg-red-900 text-red-300 border border-red-900 text-xs font-semibold transition disabled:opacity-50"
            :disabled="disconnecting === conn.id"
            @click="disconnect(conn.id)"
          >
            {{ disconnecting === conn.id ? 'Disconnecting…' : 'Disconnect' }}
          </button>
        </div>
      </div>

      <div v-if="connections.length === 0" class="text-center py-12 bg-slate-900/40 rounded-xl border border-slate-800">
        <p class="text-sm text-slate-300 font-semibold">No cloud accounts connected yet</p>
        <p class="text-xs text-slate-500 mt-1">Add your first connection below to start managing servers.</p>
      </div>
    </div>

    <!-- Connect form -->
    <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-5">
      <h3 class="text-sm font-bold text-white">Connect a new cloud account</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="conn-provider" class="text-xs font-semibold text-slate-300">Provider</label>
          <select
            id="conn-provider"
            v-model="form.providerSlug"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          >
            <option value="digitalocean">DigitalOcean</option>
            <option value="hetzner">Hetzner Cloud</option>
            <option value="vultr">Vultr</option>
            <option value="upcloud">UpCloud</option>
            <option value="linode">Linode (Akamai)</option>
          </select>
        </div>
        <div class="space-y-1.5">
          <label for="conn-name" class="text-xs font-semibold text-slate-300">Connection name</label>
          <input
            id="conn-name"
            v-model="form.name"
            type="text"
            placeholder="Production DO Account"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
        <div class="space-y-1.5">
          <label for="conn-access" class="text-xs font-semibold text-slate-300">Access level</label>
          <select
            id="conn-access"
            v-model="form.accessLevel"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          >
            <option value="READ_ONLY">Read only (inventory + cost)</option>
            <option value="READ_WRITE">Read + write (provisioning)</option>
          </select>
        </div>
        <div class="space-y-1.5">
          <label for="conn-token" class="text-xs font-semibold text-slate-300">Scoped API token</label>
          <input
            id="conn-token"
            v-model="form.token"
            type="password"
            autocomplete="off"
            placeholder="Paste a scoped token from your provider"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 text-sm focus:outline-none focus:border-sky-500"
          />
        </div>
      </div>

      <div v-if="formError" class="px-4 py-3 rounded-lg bg-red-950/40 border border-red-900 text-xs text-red-300" role="alert">
        {{ formError }}
      </div>

      <button
        class="px-5 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 disabled:opacity-60 text-white text-sm font-semibold transition"
        :disabled="connecting"
        @click="connect"
      >
        {{ connecting ? 'Verifying & sealing token…' : 'Connect Account' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
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
const formError = ref('')

const form = reactive({
  providerSlug: 'digitalocean',
  name: '',
  accessLevel: 'READ_ONLY',
  token: '',
})

const { request } = useApi()

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
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    formError.value = apiErr.message
  } finally {
    connecting.value = false
  }
}

async function disconnect(id: string) {
  disconnecting.value = id
  try {
    await request(`/api/v1/cloud/connections/${id}`, { method: 'DELETE' })
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr
  } finally {
    disconnecting.value = ''
  }
}

onMounted(() => {
  load()
})
</script>
