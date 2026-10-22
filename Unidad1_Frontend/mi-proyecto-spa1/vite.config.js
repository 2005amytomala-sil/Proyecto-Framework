import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'


// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(), VitePWA({
    registerType: 'prompt',
    includeAssets: ['favicon.svg', 'pwa-192.png', 'pwa-512.png'],
    manifest: {
      name: 'MultiCatálogo',
      short_name: 'MultiCatálogo',
      description: 'Tienda y red de ventas MultiCatálogo',
      start_url: '/',
      scope: '/',
      display: 'standalone',
      orientation: 'portrait-primary',
      theme_color: '#312e81',
      background_color: '#f8fafc',
      icons: [
        { src: '/pwa-192.png', sizes: '192x192', type: 'image/png', purpose: 'any' },
        { src: '/pwa-192.png', sizes: '192x192', type: 'image/png', purpose: 'maskable' },
        { src: '/pwa-512.png', sizes: '512x512', type: 'image/png', purpose: 'any' },
        { src: '/pwa-512.png', sizes: '512x512', type: 'image/png', purpose: 'maskable' },
      ],
    },
    workbox: {
      navigateFallback: '/index.html',
      runtimeCaching: [
        {
          urlPattern: ({ url }) => url.pathname.startsWith('/api/productos'),
          handler: 'NetworkFirst',
          options: { cacheName: 'catalogo-publico', expiration: { maxEntries: 30, maxAgeSeconds: 300 } },
        },
        {
          urlPattern: ({ request }) => request.destination === 'image',
          handler: 'CacheFirst',
          options: { cacheName: 'imagenes-productos', expiration: { maxEntries: 80, maxAgeSeconds: 604800 } },
        },
      ],
    },
    devOptions: { enabled: false },
  }),
  ],
   server: {
    host: true,
    port: 5173,
  },
})
