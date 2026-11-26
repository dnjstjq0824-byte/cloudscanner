<template>
  <div class="max-w-md mx-auto py-12">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('login.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('login.subtitle') }}</p>
    </div>

    <form class="mt-8 space-y-5 p-8 card" @submit.prevent="onSubmit">
      <div class="space-y-1.5">
        <label for="login-email" class="label">{{ $t('auth.email') }}</label>
        <input
          id="login-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          :placeholder="$t('auth.emailPlaceholder')"
          class="input"
        />
      </div>

      <div class="space-y-1.5">
        <label for="login-password" class="label">{{ $t('auth.password') }}</label>
        <input
          id="login-password"
          v-model="password"
          type="password"
          required
          autocomplete="current-password"
          placeholder="••••••••"
          class="input"
        />
      </div>

      <div v-if="error" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ error }}
      </div>

      <button type="submit" :disabled="submitting" class="btn-primary w-full !py-3">
        {{ submitting ? $t('auth.signingIn') : $t('auth.signIn') }}
      </button>

      <p class="text-center text-xs text-ink-secondary">
        {{ $t('auth.newTo') }}
        <NuxtLink :to="$localePath('/auth/register')" class="text-brand hover:underline font-semibold">{{ $t('auth.createAccountLink') }}</NuxtLink>
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

const { t } = useI18n()
const localePath = useLocalePath()
useHead(() => ({
  title: t('login.title'),
  meta: [{ name: 'description', content: t('login.meta') }],
}))

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
    await navigateTo(localePath('/dashboard'))
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr.status === 401 ? t('auth.invalidCredentials') : apiErr.message
  } finally {
    submitting.value = false
  }
}
</script>
