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
</script>
