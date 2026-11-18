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
        lang: 'en',
        class: 'dark'
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Compare prices across 18+ cloud providers with 36-dimension normalization and manage your servers with zero markup.' },
        { name: 'theme-color', content: '#0b0d12' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'canonical', href: 'https://cloudscanr.com' },
        { rel: 'preconnect', href: 'https://cdn.jsdelivr.net' },
        {
          rel: 'stylesheet',
          href: 'https://cdn.jsdelivr.net/gh/orioncactus/pretendard@v1.3.9/dist/web/variable/pretendardvariable-dynamic-subset.min.css'
        }
      ],
      script: [
        {
          innerHTML: `(function(){try{var t=localStorage.getItem('cs-theme')||'dark';document.documentElement.classList.toggle('dark',t==='dark');}catch(e){}})();`
        }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  // Nitro honors the NITRO_PRESET env var natively: Cloudflare Pages builds
  // set NITRO_PRESET=cloudflare-pages without any config change here.

  runtimeConfig: {
    // NUXT_PUBLIC_API_BASE env var overrides this default automatically.
    public: {
      apiBase: 'http://localhost:8080',
    },
  },

  // Hybrid rendering strategy (Article 16 & Section F).
  // Development proxies same-origin /api/** to the Go backend; override the
  // target per environment with NUXT_API_PROXY_TARGET. Production browsers
  // call NUXT_PUBLIC_API_BASE directly (see useApi), so Pages needs no proxy.
  routeRules: {
    '/': { ssr: true, swr: 300 },
    '/search/**': { ssr: true, swr: 120 },
    '/compare/**': { ssr: true, swr: 300 },
    '/providers/**': { ssr: true, swr: 600 },
    '/articles/**': { ssr: true, swr: 3600 },
    '/dashboard/**': { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } },
    '/admin/**': { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } },
    '/api/**': { proxy: `${process.env.NUXT_API_PROXY_TARGET || 'http://localhost:8080'}/api/**` }
  },

  typescript: {
    strict: true,
    typeCheck: false
  }
})
