// src/context/UIContext.tsx
import { createContext, useContext, useState, type ReactNode } from "react";

interface UIContextType {
  isCollapsed: boolean;
  isMobileOpen: boolean;
  toggleSidebar: () => void;
  closeMobileMenu: () => void;
}

const UIContext = createContext<UIContextType | undefined>(undefined);

export const useUI = () => {
  const context = useContext(UIContext);
  if (!context) {
    throw new Error("useUI debe ser usado dentro de un UIProvider");
  }
  return context;
};

export const UIProvider = ({ children }: { children: ReactNode }) => {
  const [isCollapsed, setIsCollapsed] = useState<boolean>(false);
  const [isMobileOpen, setIsMobileOpen] = useState<boolean>(false);

  const toggleSidebar = () => {
    // Si la pantalla es móvil (< 768px), alterna el menú deslizable
    if (window.innerWidth < 768) {
      setIsMobileOpen((prev) => !prev);
    } else {
      // En pantallas de escritorio, colapsa el sidebar a iconos
      setIsCollapsed((prev) => !prev);
    }
  };

  const closeMobileMenu = () => {
    setIsMobileOpen(false);
  };

  return (
    <UIContext.Provider
      value={{
        isCollapsed,
        isMobileOpen,
        toggleSidebar,
        closeMobileMenu,
      }}
    >
      {children}
    </UIContext.Provider>
  );
};