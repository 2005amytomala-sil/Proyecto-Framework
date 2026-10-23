import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";

interface NavbarProps { toggleSidebar: () => void; }

const Navbar = ({ toggleSidebar }: NavbarProps) => {
  const { totalItems } = useCart();
  const { logout, user } = useAuth();
  const navigate = useNavigate();
  const [profileOpen, setProfileOpen] = useState(false);
  const profileRef = useRef<HTMLDivElement>(null);

  const handleLogout = () => {
    if (window.confirm("¿Deseas cerrar sesión?")) {
      logout();
      navigate("/login");
    }
    setProfileOpen(false);
  };

  useEffect(() => {
    const closeOnOutsideClick = (event: MouseEvent) => {
      if (profileRef.current && !profileRef.current.contains(event.target as Node)) setProfileOpen(false);
    };
    document.addEventListener("mousedown", closeOnOutsideClick);
    return () => document.removeEventListener("mousedown", closeOnOutsideClick);
  }, []);

  return (
    <header className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 md:px-8">
      <div className="flex items-center gap-2 md:gap-4">
        <button onClick={toggleSidebar} className="rounded-md p-2 text-slate-600 transition hover:bg-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500" aria-label="Alternar menú">
          <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" /></svg>
        </button>
        <h2 className="hidden text-base font-medium text-slate-600 sm:block md:text-lg">{user?.rol === "admin" ? "Panel de Administración" : "Tienda MultiCatálogo"}</h2>
      </div>

      <div className="flex items-center gap-4 md:gap-6">
        <Link to="/carrito" className="relative rounded-full p-2 transition hover:bg-slate-100" aria-label="Carrito">
          <span className="text-xl">🛒</span>
          {totalItems > 0 && <span className="absolute right-0 top-0 flex h-5 w-5 -translate-y-1 translate-x-1 items-center justify-center rounded-full bg-indigo-600 text-xs font-bold text-white">{totalItems}</span>}
        </Link>
        <div className="flex items-center gap-2 md:gap-4">
          <span className="hidden text-xs text-slate-500 sm:block md:text-sm">{user?.email}</span>
          <span className={`hidden rounded-full px-2 py-1 text-xs font-semibold uppercase md:inline-block ${user?.rol === "admin" ? "bg-amber-100 text-amber-700" : "bg-indigo-100 text-indigo-700"}`}>{user?.rol}</span>
          <div ref={profileRef} className="relative pb-2">
            <button type="button" onClick={() => setProfileOpen((open) => !open)} className="block h-9 w-9 overflow-hidden rounded-full border border-indigo-200 bg-slate-200 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500" aria-label="Abrir menú de usuario" aria-expanded={profileOpen}>
              <span className="flex h-full w-full items-center justify-center bg-gradient-to-br from-indigo-500 via-violet-500 to-fuchsia-500 text-white">
                <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><circle cx="12" cy="8" r="3.25" strokeWidth="1.8" /><path strokeLinecap="round" strokeWidth="1.8" d="M5.5 19a6.5 6.5 0 0 1 13 0" /></svg>
              </span>
            </button>
            <div className={`absolute right-0 top-full z-50 mt-1 w-40 rounded-lg border border-slate-200 bg-white shadow-lg transition-all duration-200 ${profileOpen ? "visible opacity-100" : "invisible opacity-0"}`}>
              <button onClick={handleLogout} className="w-full rounded-md px-4 py-2 text-left text-sm font-semibold text-red-600 transition-colors hover:bg-red-50">Cerrar sesión</button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Navbar;
