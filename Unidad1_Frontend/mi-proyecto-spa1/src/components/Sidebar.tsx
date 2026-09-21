import { Link } from "react-router-dom";

interface SidebarProps {
  isCollapsed: boolean;
  onCloseMobile?: () => void; 
}

const Sidebar = ({ isCollapsed, onCloseMobile }: SidebarProps) => {
  const menu = [
    { to: "/", icono: "📊", texto: "Dashboard" },
    { to: "/catalogo", icono: "📦", texto: "Catálogo" },
    { to: "/mi-red", icono: "👥", texto: "Mi Red" },
  ];

  return (
    <aside
      className={`
        bg-[#440309] text-white flex flex-col transition-all duration-300 h-screen shrink-0
        ${isCollapsed ? "w-20" : "w-64"}
      `}
    >
      <div className="h-16 flex items-center justify-center p-4 text-xl font-bold border-b border-white/20 whitespace-nowrap overflow-hidden">
        {isCollapsed ? "MC" : "MultiCatálogo"}
      </div>

      {/* Opciones del menú */}
      <nav className="flex-1 p-3 space-y-2">
        {menu.map((item) => (
          <Link
            key={item.to}
            to={item.to}
            onClick={onCloseMobile} // Cierra/colapsa al navegar en móvil
            title={item.texto}
            className={`flex items-center gap-3 p-3 rounded-lg hover:bg-[#7A0509] transition-colors ${
              isCollapsed ? "justify-center" : "justify-start"
            }`}
          >
            <span className="text-xl shrink-0">{item.icono}</span>
            {!isCollapsed && <span className="truncate">{item.texto}</span>}
          </Link>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;