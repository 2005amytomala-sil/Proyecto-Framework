import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';

// 1. Tipos de rol que maneja la aplicación
export type Rol = 'admin' | 'cliente';

// 2. Usuario autenticado: correo + rol + token (lo entrega la API en /api/login)
export interface Usuario {
  email: string;
  rol: Rol;
  token: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: Usuario | null;
  login: (usuario: Usuario) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<Usuario | null>(null);

  // Cargar sesión desde localStorage al montar
  useEffect(() => {
    const savedAuth = localStorage.getItem('multicatalogo_sesion');
    if (savedAuth) {
      try {
        const parsed = JSON.parse(savedAuth);
        setUser(parsed);
        setIsAuthenticated(true);
      } catch {
        localStorage.removeItem('multicatalogo_sesion');
      }
    }
  }, []);

  const login = (usuario: Usuario) => {
    setIsAuthenticated(true);
    setUser(usuario);
    // Guardar en localStorage
    localStorage.setItem('multicatalogo_sesion', JSON.stringify(usuario));
  };

  const logout = () => {
    setIsAuthenticated(false);
    setUser(null);
    // Limpiar localStorage
    localStorage.removeItem('multicatalogo_sesion');
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};