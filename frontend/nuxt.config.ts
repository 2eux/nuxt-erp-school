import type { NuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  future: {
    compatibilityVersion: 4,
  },

  ssr: true,

  devtools: { enabled: true },

  modules: [
    '@nuxt/ui',
    '@nuxtjs/i18n',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxt/image',
    '@nuxt/fonts',
    '@nuxtjs/color-mode',
    '@vite-pwa/nuxt',
  ],

  css: ['~/assets/css/tailwind.css'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
      appName: process.env.NUXT_PUBLIC_APP_NAME || 'Islamic School ERP',
      appUrl: process.env.NUXT_PUBLIC_APP_URL || 'http://localhost:3000',
    },
  },

  colorMode: {
    preference: 'system',
    fallback: 'light',
    classSuffix: '',
  },

  i18n: {
    strategy: 'prefix_except_default',
    defaultLocale: 'id',
    langDir: 'i18n/locales',
    locales: [
      { code: 'id', iso: 'id-ID', name: 'Bahasa Indonesia', file: 'id.ts' },
      { code: 'en', iso: 'en-US', name: 'English', file: 'en.ts' },
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      alwaysRedirect: false,
    },
    vueI18n: './i18n.config.ts',
  },

  ui: {
    icons: ['heroicons', 'lucide'],
  },

  image: {
    domains: ['localhost'],
    format: ['webp'],
  },

  fonts: {
    families: [
      { name: 'Inter', provider: 'google' },
      { name: 'Amiri', provider: 'google' },
    ],
  },

  pwa: {
    registerType: 'autoUpdate',
    manifest: {
      name: 'Islamic School ERP',
      short_name: 'SchoolERP',
      theme_color: '#047857',
      background_color: '#ffffff',
      display: 'standalone',
      orientation: 'portrait',
      icons: [
        {
          src: '/icon-192.png',
          sizes: '192x192',
          type: 'image/png',
        },
        {
          src: '/icon-512.png',
          sizes: '512x512',
          type: 'image/png',
        },
      ],
    },
    workbox: {
      navigateFallback: '/',
    },
  },

  routeRules: {
    '/auth/**': { ssr: false },
    '/dashboard/**': { ssr: false },
    '/api/**': { cors: true },
  },

  nitro: {
    preset: 'node-server',
    devProxy: {
      '/api': {
        target: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
        changeOrigin: true,
      },
    },
  },

  vite: {
    optimizeDeps: {
      include: ['apexcharts', 'vue3-apexcharts', 'date-fns'],
    },
  },
}) satisfies NuxtConfig
