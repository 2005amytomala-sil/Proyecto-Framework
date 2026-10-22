# Plan de implementación PWA — MultiCatálogo

Este plan transforma la SPA React/Vite existente en una Progressive Web App
instalable y resistente a cortes de conexión. La API Go/PostgreSQL continúa
siendo necesaria para autenticación, checkout, pedidos, red, comisiones y
operaciones administrativas.

## FASE PWA-0 — Línea base y configuración

- Confirmar que `npm run build` y `npm run lint` pasan.
- Definir `VITE_API_URL` para desarrollo, Docker y producción.
- Instalar y configurar `vite-plugin-pwa`.
- Mantener separados los recursos estáticos del tráfico de la API.

**Criterio:** el build genera una versión de producción reproducible y la API
no se cachea como si fuera contenido estático.

## FASE PWA-1 — Manifest e identidad instalable

- Crear `manifest.webmanifest` con `name`, `short_name`, `description`,
  `start_url`, `scope`, `display: standalone`, colores y orientación.
- Generar íconos PNG de 192×192 y 512×512, incluyendo versión maskable.
- Configurar favicon, theme color y metadatos móviles en `index.html`.
- Definir el nombre visible como `MultiCatálogo`.

**Criterio:** Lighthouse/Application panel reconoce la aplicación como
instalable y muestra nombre e ícono correctos.

## FASE PWA-2 — Service Worker y caché del shell

- Registrar el Service Worker generado por Vite PWA.
- Cachear `index.html`, JavaScript, CSS, fuentes e imágenes estáticas.
- Usar estrategia `NetworkFirst` para navegación SPA con fallback al shell.
- Usar `CacheFirst` para assets versionados.
- No guardar respuestas autenticadas ni mutaciones (`POST`, `PUT`, `PATCH`,
  `DELETE`) en caché.

**Criterio:** después de una primera visita, la aplicación abre sin red y la
navegación básica de la SPA continúa funcionando.

## FASE PWA-3 — Catálogo y experiencia degradada offline

- Cachear el catálogo público con una estrategia `StaleWhileRevalidate` o
  `NetworkFirst` con fallback.
- Mostrar fecha o indicador de datos almacenados cuando no haya conexión.
- Mostrar mensajes claros de “sin conexión” en login, checkout, pedidos,
  red, comisiones y administración.
- Evitar enviar checkout offline; conservar el carrito local y permitir
  reintentar cuando vuelva la conexión.

**Criterio:** el catálogo previamente visitado sigue visible sin red y las
operaciones transaccionales nunca se confirman sin respuesta del backend.

## FASE PWA-4 — Sesión, actualización y ciclo de vida

- Mantener la sesión actual en `localStorage`, sin almacenar contraseñas.
- Detectar expiración del JWT y solicitar login nuevamente.
- Mostrar aviso cuando haya una versión nueva disponible.
- Añadir acciones “Actualizar ahora” y “Más tarde”.
- Limpiar cachés antiguas durante la activación del Service Worker.

**Criterio:** una actualización no deja una mezcla de assets viejos y nuevos,
y el usuario puede actualizar sin borrar manualmente los datos de sesión o
carrito.

## FASE PWA-5 — Seguridad y despliegue

- Servir la PWA mediante HTTPS en producción.
- Configurar la URL HTTPS real de la API mediante `VITE_API_URL`.
- Verificar CORS para el dominio de la PWA.
- No incluir secretos en el bundle frontend.
- Configurar headers adecuados: `Content-Security-Policy`,
  `Cache-Control` y `Service-Worker-Allowed` cuando corresponda.

**Criterio:** la instalación funciona en navegador de escritorio y móvil sin
contenido mixto ni errores CORS.

## FASE PWA-6 — Pruebas y aceptación

- Probar instalación en Chrome/Edge de escritorio y Android.
- Probar carga inicial, recarga offline y navegación interna.
- Probar catálogo offline después de una visita previa.
- Probar recuperación de red.
- Confirmar que login, checkout, pedidos y administración requieren conexión.
- Ejecutar `npm run build`, `npm run lint` y auditoría Lighthouse.
- Verificar que el Service Worker no sirve respuestas privadas a otro usuario.

**Criterio de cierre:** Lighthouse PWA pasa los criterios aplicables, la app
se instala, abre offline con shell y catálogo previamente almacenado, y las
operaciones transaccionales permanecen protegidas por la API.

## Orden de ejecución

PWA-0 → PWA-1 → PWA-2 → PWA-3 → PWA-4 → PWA-5 → PWA-6.

No se implementará sincronización automática de checkout offline: una compra
requiere validación online de identidad, precio y stock en PostgreSQL.
