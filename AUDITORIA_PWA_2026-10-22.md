# Auditoría PWA — 2026-10-22

## Implementado en esta etapa

- `vite-plugin-pwa` instalado.
- Manifest generado con nombre, alcance, colores e íconos.
- Service Worker generado por Workbox.
- Precaché del shell de la SPA.
- Caché `NetworkFirst` del catálogo público.
- Caché `CacheFirst` de imágenes.
- `VITE_API_URL` configurable por entorno.
- Mutaciones y rutas autenticadas no se incluyen en runtime caching.

## Evidencia

`npm run build` genera:

- `dist/manifest.webmanifest`
- `dist/sw.js`
- `dist/registerSW.js`
- `dist/workbox-*.js`

`npm run lint` pasa correctamente.

## Pendiente para cerrar PWA-1 y PWA-2

- Instalar la aplicación desde Chrome/Edge o Android.
- Confirmar apertura offline mediante `npm run preview` o un servidor HTTPS.
- Verificar actualización del Service Worker y comportamiento en DevTools.

Las operaciones de login, checkout, pedidos y administración requieren red y
no se almacenan en caché.
