<template>
  <div class="max-w-md mx-auto py-12">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">Create your account</h1>
      <p class="text-sm text-ink-secondary">One account. Every cloud you own. Zero markup.</p>
    </div>

    <form class="mt-8 space-y-5 p-8 card" @submit.prevent="onSubmit">
      <div class="space-y-1.5">
        <label for="reg-name" class="label">Full name</label>
        <input
          id="reg-name"
          v-model="fullName"
          type="text"
          required
          autocomplete="name"
          placeholder="Ada Lovelace"
          class="input"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-email" class="label">Email</label>
        <input
          id="reg-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@company.com"
          class="input"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-password" class="label">Password</label>
        <input
          id="reg-password"
          v-model="password"
          type="password"
          required
          minlength="8"
          autocomplete="new-password"
          placeholder="Minimum 8 characters"
          class="input"
        />
        <p class="text-[11px] text-ink-secondary">Hashed with Argon2id. Never stored in plain text.</p>
      </div>

      <div v-if="error" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ error }}
      </div>

      <button type="submit" :disabled="submitting" class="btn-primary w-full !py-3">
        {{ submitting ? 'Creating account…' : 'Create Account' }}
      </button>

      <p class="text-center text-[11px] text-ink-secondary leading-relaxed">
        By registering you agree to the
        <NuxtLink to="/terms" class="text-brand hover:underline">Terms of Service</NuxtLink>
        and
        <NuxtLink to="/privacy" class="text-brand hover:underline">Privacy Policy</NuxtLink>.
      </p>

      <p class="text-center text-xs text-ink-secondary">
        Already have an account?
        <NuxtLink to="/auth/login" class="text-brand hover:underline font-semibold">Sign in</NuxtLink>
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
  meta: [{ name: 'description', content: 'Create a TirpScanner account for deal alerts and console access.' }],
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
