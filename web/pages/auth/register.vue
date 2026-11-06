<template>
  <div class="max-w-md mx-auto py-12">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Create your account</h1>
      <p class="text-sm text-slate-400">One account. Every cloud you own. Zero markup.</p>
    </div>

    <form class="mt-8 space-y-5 p-8 rounded-xl bg-slate-900 border border-slate-800" @submit.prevent="onSubmit">
      <div class="space-y-1.5">
        <label for="reg-name" class="text-xs font-semibold text-slate-300">Full name</label>
        <input
          id="reg-name"
          v-model="fullName"
          type="text"
          required
          autocomplete="name"
          placeholder="Ada Lovelace"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 focus:outline-none focus:border-sky-500 transition text-sm"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-email" class="text-xs font-semibold text-slate-300">Email</label>
        <input
          id="reg-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@company.com"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 focus:outline-none focus:border-sky-500 transition text-sm"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-password" class="text-xs font-semibold text-slate-300">Password</label>
        <input
          id="reg-password"
          v-model="password"
          type="password"
          required
          minlength="8"
          autocomplete="new-password"
          placeholder="Minimum 8 characters"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white placeholder-slate-600 focus:outline-none focus:border-sky-500 transition text-sm"
        />
        <p class="text-[11px] text-slate-500">Hashed with Argon2id. Never stored in plain text.</p>
      </div>

      <div v-if="error" class="px-4 py-3 rounded-lg bg-red-950/40 border border-red-900 text-xs text-red-300" role="alert">
        {{ error }}
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full py-3 rounded-lg bg-sky-600 hover:bg-sky-500 disabled:opacity-60 text-white font-semibold text-sm transition shadow"
      >
        {{ submitting ? 'Creating account…' : 'Create Account' }}
      </button>

      <p class="text-center text-[11px] text-slate-500 leading-relaxed">
        By registering you agree to the
        <NuxtLink to="/terms" class="text-sky-400 hover:text-sky-300">Terms of Service</NuxtLink>
        and
        <NuxtLink to="/privacy" class="text-sky-400 hover:text-sky-300">Privacy Policy</NuxtLink>.
      </p>

      <p class="text-center text-xs text-slate-400">
        Already have an account?
        <NuxtLink to="/auth/login" class="text-sky-400 hover:text-sky-300 font-semibold">Sign in</NuxtLink>
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
  title: 'Create Account',
  meta: [{ name: 'description', content: 'Create a CloudScanr account and connect your own cloud providers.' }],
})

const fullName = ref('')
const email = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')

async function onSubmit() {
  submitting.value = true
  error.value = ''
  try {
    const auth = useAuthStore()
    await auth.register(email.value.trim(), password.value, fullName.value.trim())
    await navigateTo('/dashboard')
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr.message
  } finally {
    submitting.value = false
  }
}
</script>
