<template>
  <div class="max-w-4xl mx-auto space-y-8 py-4">
    <div class="space-y-3">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-800 border border-slate-700 text-slate-300 text-xs font-semibold">
        <span>{{ $t('disclosure.badge') }}</span>
      </div>
      <h1 class="text-3xl font-extrabold text-white tracking-tight">{{ $t('disclosure.title') }}</h1>
      <p class="text-sm text-slate-400">
        {{ $t('disclosure.intro') }}
      </p>
    </div>

    <!-- Core Disclosures -->
    <div class="space-y-6 text-sm text-slate-300 leading-relaxed bg-slate-900/60 p-8 rounded-2xl border border-slate-800">
      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">{{ $t('disclosure.s1title') }}</h2>
        <p>
          {{ $t('disclosure.s1body') }}
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">{{ $t('disclosure.s2title') }}</h2>
        <p>
          {{ $t('disclosure.s2pre') }} <code class="text-sky-400 bg-slate-800 px-1.5 py-0.5 rounded">/go/{slug}</code>{{ $t('disclosure.s2post') }}
        </p>
        <p class="text-slate-400 text-xs">
          {{ $t('disclosure.s2note') }}
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">{{ $t('disclosure.s3title') }}</h2>
        <p class="text-slate-300">
          {{ $t('disclosure.s3pre') }}<code class="text-sky-400 bg-slate-800 px-1.5 py-0.5 rounded">/go/*</code>{{ $t('disclosure.s3post') }}
        </p>
        <p class="text-slate-400 text-xs">
          {{ $t('disclosure.s3note') }}
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">{{ $t('disclosure.s4title') }}</h2>
        <p class="text-slate-400 text-xs">
          {{ $t('disclosure.s4body') }}
        </p>
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
  title: t('disclosure.title'),
  meta: [
    { name: 'description', content: t('disclosure.meta') },
  ],
}))
</script>
