<template>
  <div class="max-w-md mx-auto py-12">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Sign in to CloudScanr</h1>
      <p class="text-sm text-slate-400">Manage your connected clouds from one console.</p>
    </div>

    <form class="mt-8 space-y-5 p-8 rounded-xl bg-slate-900 border border-slate-800" @submit.prevent="onSubmit">
      <div class="space-y-1.5">
        <label for="login-email" class="text-xs font-semibold text-slate-300">Email</label>
        <input
          id="login-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@company.com"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 focus:outline-none focus:border-sky-500 transition text-sm"
        />
      </div>

      <div class="space-y-1.5">
        <label for="login-password" class="text-xs font-semibold text-slate-300">Password</label>
        <input
          id="login-password"
          v-model="password"
          type="password"
          required
          autocomplete="current-password"
          placeholder="••••••••"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 focus:outline-none focus:border-sky-500 transition text-sm"
        />
      </div>

      <div v-if="error" class="px-4 py-3 rounded-lg bg-red-950/40 border border-red-900 text-xs text-red-300" role="alert">
        {{ error }}
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full py-3 rounded-lg bg-sky-600 hover:bg-sky-500 disabled:opacity-60 text-white font-semibold text-sm transition shadow"
      >
        {{ submitting ? 'Signing in…' : 'Sign In' }}
      </button>

      <p class="text-center text-xs text-slate-400">
        New to CloudScanr?
        <NuxtLink to="/auth/register" class="text-sky-400 hover:text-sky-300 font-semibold">Create an account</NuxtLink>
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ApiError } from '~/composables/useApi'

definePageMeta({
  layout: 'public',
  middleware: 'guest',
})

useHead({
  title: 'Sign In',
  meta: [{ name: 'description', content: 'Sign in to the CloudScanr BYOC multi-cloud console.' }],
})

const email = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')

async function onSubmit() {
  submitting.value = true
  error.value = ''
  try {
    const auth = useAuthStore()
    await auth.login(email.value.trim(), password.value)
    await navigateTo('/dashboard')
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr.status === 401 ? 'Invalid email or password.' : apiErr.message
  } finally {
    submitting.value = false
  }
}
</script>
