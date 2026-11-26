<template>
  <div class="min-h-screen bg-canvas text-ink flex flex-col">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
// TirpScanner Core Root Component: rehydrate the session once per render so
// SSR and client agree on the authenticated header state.
const auth = useAuthStore()
if (!auth.hydrated) {
  await useAsyncData('session', () => auth.fetchMe(), { server: true, lazy: false })
}

// Localized document head: translated title/description plus the i18n
// module's hreflang alternates, canonical link, and html lang attribute.
const { t, locale } = useI18n()
const i18nHead = useLocaleHead()
useHead(() => ({
  title: t('meta.title'),
  htmlAttrs: { lang: locale.value, ...(i18nHead.value.htmlAttrs || {}) },
  link: [...(i18nHead.value.link || [])],
  meta: [
    { name: 'description', content: t('meta.description') },
    ...(i18nHead.value.meta || [])
  ]
}))
</script>
