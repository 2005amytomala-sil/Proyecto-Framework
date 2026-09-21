// src/components/Navbar.tsx
import { Link, useNavigate } from "react-router-dom";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";
import { useUI } from "../context/UIContext";

const Navbar = () => {
  const { totalItems } = useCart();
  const { logout, userEmail } = useAuth();
  const { toggleSidebar } = useUI();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header className="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-4 md:px-8 shrink-0">
      {/* Izquierda: Botón Toggle y Título */}
      <div className="flex items-center gap-3">
        <button
          onClick={toggleSidebar}
          className="p-2 rounded-lg hover:bg-slate-100 transition text-slate-700 text-xl focus:outline-none focus:ring-2 focus:ring-slate-300"
          title="Mostrar u ocultar menú"
          aria-label="Alternar menú lateral"
        >
          ☰
        </button>
        <h2 className="text-slate-700 font-medium text-base md:text-lg truncate max-w-[180px] sm:max-w-none">
          Panel de Administración
        </h2>
      </div>

      {/* Derecha: Carrito y Perfil */}
      <div className="flex items-center gap-4 sm:gap-6">
        {/* Carrito */}
        <Link
          to="/carrito"
          className="relative p-2 hover:bg-slate-100 rounded-full transition text-slate-700"
          title="Ver carrito"
        >
          <span className="text-xl">🛒</span>
          {totalItems > 0 && (
            <span className="absolute top-0 right-0 bg-indigo-600 text-white text-xs font-bold w-5 h-5 flex items-center justify-center rounded-full transform translate-x-1 -translate-y-1">
              {totalItems}
            </span>
          )}
        </Link>

        {/* Usuario y Menú */}
        <div className="flex items-center gap-3">
          <span className="hidden sm:inline-block text-sm text-slate-500 max-w-[150px] truncate">
            {userEmail}
          </span>

          <div className="relative group cursor-pointer py-1">
            <div className="w-9 h-9 rounded-full bg-slate-200 overflow-hidden border border-slate-300 flex items-center justify-center">
              <img
                src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80"
                alt="Avatar del usuario"
                className="w-full h-full object-cover"
              />
            </div>

            {/* Dropdown Logout */}
            <div className="absolute right-0 top-full mt-1 w-40 bg-white border border-slate-200 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
              <div className="px-4 py-2 border-b border-slate-100 sm:hidden">
                <p className="text-xs text-slate-400">Sesión iniciada como</p>
                <p className="text-xs font-medium text-slate-600 truncate">{userEmail}</p>
              </div>
              <button
                onClick={handleLogout}
                className="w-full text-left px-4 py-2 text-sm text-red-600 font-semibold hover:bg-red-50 rounded-b-lg transition-colors"
              >
                Cerrar Sesión
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Navbar;