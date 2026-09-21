// src/components/Layout.tsx
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";
import { useUI } from "../context/UIContext";

const Layout = () => {
  const { isCollapsed, toggleSidebar } = useUI();

  // Si está abierto en celular y tocamos un enlace, lo colapsa a solo iconos
  const handleCloseMobile = () => {
    if (!isCollapsed && window.innerWidth < 768) {
      toggleSidebar();
    }
  };

  return (
    <div className="flex h-screen bg-slate-50 relative overflow-hidden">
      <Sidebar 
        isCollapsed={isCollapsed} 
        onCloseMobile={handleCloseMobile} 
      />

      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <Navbar />
        <main className="flex-1 overflow-x-hidden overflow-y-auto p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;