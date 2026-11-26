<template>
  <div class="max-w-3xl mx-auto space-y-8 py-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">{{ $t('privacy.title') }}</h1>
      <p class="text-xs text-slate-500">{{ $t('privacy.updated') }}</p>
    </div>

    <div class="space-y-6 text-sm text-slate-300 leading-relaxed">
      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s1title') }}</h2>
        <ul class="list-disc list-inside space-y-1">
          <li>{{ $t('privacy.s1a') }}</li>
          <li>{{ $t('privacy.s1b') }}</li>
          <li>{{ $t('privacy.s1c') }}</li>
          <li>{{ $t('privacy.s1d') }}</li>
        </ul>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s2title') }}</h2>
        <p>{{ $t('privacy.s2body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s3title') }}</h2>
        <p>{{ $t('privacy.s3body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s4title') }}</h2>
        <p>{{ $t('privacy.s4body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s5title') }}</h2>
        <p>{{ $t('privacy.s5body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('privacy.s6title') }}</h2>
        <p>{{ $t('privacy.s6body') }}</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public'
})

const { t } = useI18n()
useHead(() => ({
  title: t('privacy.title'),
  meta: [{ name: 'description', content: t('privacy.meta') }],
}))
</script>
