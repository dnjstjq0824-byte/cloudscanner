// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-10-28',
  devtools: { enabled: false },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],

  app: {
    head: {
      titleTemplate: '%s — CloudScanr',
      title: 'CloudScanr: Independent Global Cloud Price Comparison & BYOC Multi-Cloud Management',
      htmlAttrs: {
        lang: 'en'
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Compare prices across 18+ cloud providers with 36-dimension normalization and manage your servers with zero markup.' },
        { name: 'theme-color', content: '#0f172a' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'canonical', href: 'https://cloudscanr.com' }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    // NUXT_PUBLIC_API_BASE env var overrides this default automatically.
    public: {
      apiBase: 'http://localhost:8080',
    },
  },

  // Hybrid rendering strategy (Article 16 & Section F)
  routeRules: {
    '/': { ssr: true, swr: 300 },
    '/search/**': { ssr: true, swr: 120 },
    '/compare/**': { ssr: true, swr: 300 },
    '/providers/**': { ssr: true, swr: 600 },
    '/articles/**': { ssr: true, swr: 3600 },
    '/dashboard/**': { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } },
    '/admin/**': { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } },
    '/api/**': { proxy: 'http://localhost:8080/api/**' }
  },

  typescript: {
    strict: true,
    typeCheck: false
  }
})

