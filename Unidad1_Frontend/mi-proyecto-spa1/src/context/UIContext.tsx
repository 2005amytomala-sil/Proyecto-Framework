import { createContext, useContext, useState, useEffect, type ReactNode } from "react";

interface UIContextType {
  isCollapsed: boolean;
  toggleSidebar: () => void;
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
  // Inicia en true (80px, solo iconos) si entra desde un celular (< 768px)
  const [isCollapsed, setIsCollapsed] = useState<boolean>(() => {
    return typeof window !== "undefined" ? window.innerWidth < 768 : false;
  });

  // Detecta cambios de tamaño de pantalla en tiempo real
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth < 768) {
        setIsCollapsed(true); 
      } else {
        setIsCollapsed(false); 
      }
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const toggleSidebar = () => setIsCollapsed((prev) => !prev);

  return (
    <UIContext.Provider value={{ isCollapsed, toggleSidebar }}>
      {children}
    </UIContext.Provider>
  );
};