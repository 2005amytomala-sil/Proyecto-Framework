// src/components/Layout.tsx
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Navbar from "./Navbar";
import { useUI } from "../context/UIContext";

const Layout = () => {
  const { isCollapsed, isMobileOpen, toggleSidebar, closeMobileMenu } = useUI();

  return (
    <div className="flex h-screen bg-slate-50 relative overflow-hidden">
      {/* Fondo oscuro para cerrar al tocar fuera en pantallas móviles */}
      {isMobileOpen && (
        <div
          onClick={closeMobileMenu}
          className="fixed inset-0 bg-slate-900/50 z-40 md:hidden"
        />
      )}

      {/* Sidebar con las props exactas que espera Sidebar.tsx */}
      <Sidebar 
        isCollapsed={isCollapsed} 
        isMobileOpen={isMobileOpen}
        closeMobileMenu={closeMobileMenu} 
      />

      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Navbar ahora sí recibe toggleSidebar para accionar el botón */}
        <Navbar toggleSidebar={toggleSidebar} />
        
        <main className="flex-1 overflow-x-hidden overflow-y-auto p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;