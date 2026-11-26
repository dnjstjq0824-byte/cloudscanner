<template>
  <div class="max-w-md mx-auto py-12">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('register.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('register.subtitle') }}</p>
    </div>

    <form class="mt-8 space-y-5 p-8 card" @submit.prevent="onSubmit">
      <div class="space-y-1.5">
        <label for="reg-name" class="label">{{ $t('auth.fullName') }}</label>
        <input
          id="reg-name"
          v-model="fullName"
          type="text"
          required
          autocomplete="name"
          :placeholder="$t('auth.namePlaceholder')"
          class="input"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-email" class="label">{{ $t('auth.email') }}</label>
        <input
          id="reg-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          :placeholder="$t('auth.emailPlaceholder')"
          class="input"
        />
      </div>

      <div class="space-y-1.5">
        <label for="reg-password" class="label">{{ $t('auth.password') }}</label>
        <input
          id="reg-password"
          v-model="password"
          type="password"
          required
          minlength="8"
          autocomplete="new-password"
          :placeholder="$t('auth.passwordMinPlaceholder')"
          class="input"
        />
        <p class="text-[11px] text-ink-secondary">{{ $t('auth.hashNote') }}</p>
      </div>

      <div v-if="error" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ error }}
      </div>

      <button type="submit" :disabled="submitting" class="btn-primary w-full !py-3">
        {{ submitting ? $t('auth.creatingAccount') : $t('auth.createAccount') }}
      </button>

      <p class="text-center text-[11px] text-ink-secondary leading-relaxed">
        {{ $t('auth.agreePre') }}
        <NuxtLink :to="$localePath('/terms')" class="text-brand hover:underline">{{ $t('auth.terms') }}</NuxtLink>
        {{ $t('auth.and') }}
        <NuxtLink :to="$localePath('/privacy')" class="text-brand hover:underline">{{ $t('auth.privacy') }}</NuxtLink>.
      </p>

      <p class="text-center text-xs text-ink-secondary">
        {{ $t('auth.haveAccount') }}
        <NuxtLink :to="$localePath('/auth/login')" class="text-brand hover:underline font-semibold">{{ $t('auth.signInLink') }}</NuxtLink>
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
  title: t('register.title'),
  meta: [{ name: 'description', content: t('register.meta') }],
}))

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
    await navigateTo(localePath('/dashboard'))
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    error.value = apiErr.message
  } finally {
    submitting.value = false
  }
}
</script>
