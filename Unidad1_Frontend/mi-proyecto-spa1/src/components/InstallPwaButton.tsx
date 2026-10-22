import { useEffect, useState } from 'react';

type InstallPromptEvent = Event & {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>;
};

export default function InstallPwaButton() {
  const [installEvent, setInstallEvent] = useState<InstallPromptEvent | null>(null);
  const [installed, setInstalled] = useState(false);

  useEffect(() => {
    const handleBeforeInstall = (event: Event) => {
      event.preventDefault();
      setInstallEvent(event as InstallPromptEvent);
    };
    const handleInstalled = () => {
      setInstalled(true);
      setInstallEvent(null);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstall);
    window.addEventListener('appinstalled', handleInstalled);
    if (window.matchMedia('(display-mode: standalone)').matches) setInstalled(true);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstall);
      window.removeEventListener('appinstalled', handleInstalled);
    };
  }, []);

  if (!installEvent || installed) return null;

  const install = async () => {
    await installEvent.prompt();
    const choice = await installEvent.userChoice;
    if (choice.outcome === 'accepted') setInstalled(true);
    setInstallEvent(null);
  };

  return (
    <button
      type="button"
      onClick={install}
      className="fixed bottom-5 right-5 z-[60] inline-flex items-center gap-2 rounded-full bg-indigo-600 px-5 py-3 text-sm font-semibold text-white shadow-[0_8px_24px_rgba(49,46,129,0.4)] ring-2 ring-white/20 transition hover:-translate-y-0.5 hover:bg-indigo-700 hover:shadow-[0_10px_28px_rgba(49,46,129,0.5)] focus:outline-none focus:ring-4 focus:ring-indigo-300 sm:right-7 sm:bottom-7"
      aria-label="Instalar MultiCatálogo"
    >
      <span aria-hidden="true">⇩</span>
      Instalar app
    </button>
  );
}
