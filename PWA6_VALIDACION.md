# Validacion PWA-6

Fecha: 2026-10-22

## Comprobaciones automatizadas

- `npm run build`: correcto.
- `npm run lint`: correcto.
- `dist/manifest.webmanifest`: generado con `start_url`, `scope`, `standalone` y cuatro entradas de icono.
- `dist/sw.js`: generado y contiene caché del catálogo e imágenes.
- Las rutas de login y checkout no aparecen en las reglas de caché runtime.

## Comprobaciones manuales realizadas

- Service Worker activo en `http://localhost:5173`.
- Catálogo visible después de activar el modo Offline tras una visita previa.
- Login bloqueado offline con `ERR_INTERNET_DISCONNECTED`.
- La navegación transaccional requiere conexión con la API.

## Pendiente del entorno de publicación

- Instalar desde Chrome/Edge en escritorio.
- Instalar en Android y comprobar apertura standalone.
- Ejecutar Lighthouse sobre una URL HTTPS real.
- Desactivar Offline y confirmar recuperación de login, catálogo y checkout.
- Confirmar CORS contra el dominio HTTPS definitivo.

La validación local no sustituye la prueba de HTTPS, porque los navegadores
requieren un contexto seguro para instalar una PWA fuera de `localhost`.
