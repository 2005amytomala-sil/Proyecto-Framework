import { useEffect, useState } from 'react';
import { useRegisterSW } from 'virtual:pwa-register/react';

export default function PwaStatus() {
  const [online, setOnline] = useState(() => navigator.onLine);
  const {
    needRefresh: [needRefresh, setNeedRefresh],
    updateServiceWorker,
  } = useRegisterSW();

  useEffect(() => {
    const onlineHandler = () => setOnline(true);
    const offlineHandler = () => setOnline(false);
    window.addEventListener('online', onlineHandler);
    window.addEventListener('offline', offlineHandler);
    return () => {
      window.removeEventListener('online', onlineHandler);
      window.removeEventListener('offline', offlineHandler);
    };
  }, []);

  if (needRefresh) {
    return (
      <div className="fixed bottom-4 left-4 right-4 z-50 rounded-lg bg-indigo-900 p-4 text-sm text-white shadow-xl md:left-auto md:w-96">
        <p className="font-semibold">Hay una nueva versión disponible.</p>
        <p className="mt-1 text-indigo-100">Actualiza para cargar los últimos cambios.</p>
        <div className="mt-3 flex gap-2">
          <button className="rounded bg-white px-3 py-1.5 font-medium text-indigo-900" onClick={() => updateServiceWorker(true)}>Actualizar ahora</button>
          <button className="rounded border border-indigo-300 px-3 py-1.5" onClick={() => setNeedRefresh(false)}>Después</button>
        </div>
      </div>
    );
  }

  if (!online) {
    return (
      <div className="fixed bottom-4 left-4 right-4 z-50 rounded-lg bg-amber-600 p-3 text-sm text-white shadow-xl md:left-auto md:w-96">
        Sin conexión: login, pedidos y checkout requieren internet.
      </div>
    );
  }

  return null;
}
