// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-10-28',
  devtools: { enabled: false },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@nuxtjs/i18n'
  ],

  // TirpScanner locales (G_GLOBALIZATION.md): subpath routing (/en/...),
  // lazy JSON bundles, browser detection with cookie memory.
  i18n: {
    baseUrl: 'https://tirpscanr.com',
    locales: [
      { code: 'en', language: 'en-US', file: 'en.json', name: 'English' },
      { code: 'ko', language: 'ko-KR', file: 'ko.json', name: '한국어' },
      { code: 'ja', language: 'ja-JP', file: 'ja.json', name: '日本語' },
      { code: 'de', language: 'de-DE', file: 'de.json', name: 'Deutsch' },
      { code: 'es', language: 'es-ES', file: 'es.json', name: 'Español' }
    ],
    defaultLocale: 'en',
    langDir: 'locales',
    lazy: true,
    strategy: 'prefix',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'tirp-locale',
      redirectOn: 'root'
    },
    bundle: {
      optimizeTranslationDirective: false
    }
  },

  app: {
    head: {
      titleTemplate: '%s — TirpScanner',
      title: 'TirpScanner: Independent SaaS & AI-Tool Deal Scanner',
      htmlAttrs: {
        class: 'dark'
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'theme-color', content: '#0b0d12' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
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
    // Locales share the same rendering rules ('/' itself is locale-redirected
    // by @nuxtjs/i18n; retired-page redirects live in server/plugins/).
    ...Object.fromEntries(
      ['en', 'ko', 'ja', 'de', 'es'].flatMap((loc) => [
        [`/${loc}`, { ssr: true, swr: 300 }],
        [`/${loc}/dashboard/**`, { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } }],
        [`/${loc}/admin/**`, { ssr: false, headers: { 'X-Robots-Tag': 'noindex, nofollow' } }]
      ])
    ),
    '/api/**': { proxy: `${process.env.NUXT_API_PROXY_TARGET || 'http://localhost:8080'}/api/**` }
  },

  typescript: {
    strict: true,
    typeCheck: false
  }
})
