<template>
  <div class="max-w-3xl mx-auto space-y-8 py-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">{{ $t('terms.title') }}</h1>
      <p class="text-xs text-slate-500">{{ $t('terms.updated') }}</p>
    </div>

    <div class="space-y-6 text-sm text-slate-300 leading-relaxed">
      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s1title') }}</h2>
        <p>{{ $t('terms.s1body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s2title') }}</h2>
        <p>{{ $t('terms.s2body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s3title') }}</h2>
        <p>{{ $t('terms.s3body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s4title') }}</h2>
        <p>{{ $t('terms.s4pre') }} <NuxtLink :to="$localePath('/affiliate-disclosure')" class="text-sky-400 hover:text-sky-300">{{ $t('terms.s4link') }}</NuxtLink> {{ $t('terms.s4post') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s5title') }}</h2>
        <p>{{ $t('terms.s5body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s6title') }}</h2>
        <p>{{ $t('terms.s6body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s7title') }}</h2>
        <p>{{ $t('terms.s7body') }}</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">{{ $t('terms.s8title') }}</h2>
        <p>{{ $t('terms.s8body') }}</p>
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
  title: t('terms.title'),
  meta: [{ name: 'description', content: t('terms.meta') }],
}))
</script>
