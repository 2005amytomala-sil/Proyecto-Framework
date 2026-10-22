const API_BASE_URL = (import.meta.env.VITE_API_URL ?? "http://localhost:3000/api").replace(/\/$/, "");

// Obtener el token JWT del localStorage
const getToken = (): string | null => {
  const authData = localStorage.getItem("multicatalogo_sesion");
  if (!authData) return null;
  try {
    const parsed = JSON.parse(authData);
    return parsed.token || null;
  } catch {
    return null;
  }
};

// apiFetch centralizado con Authorization automático
export const apiFetch = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const token = getToken();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  // Agregar Authorization header si hay token
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  return response;
};

// Helper para manejar errores de API
export const handleApiError = async (response: Response): Promise<never> => {
  const error = await response.json();
  throw new Error(error.message || error.error || "Error desconocido");
};
