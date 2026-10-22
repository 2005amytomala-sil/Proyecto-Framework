import { apiFetch } from "./http";
import { type Producto } from "../data/productos";

// GET /api/productos (lista completa desde backend)
export const getProductos = async (): Promise<Producto[]> => {
  const response = await apiFetch("/productos");
  if (!response.ok) {
    throw new Error("Error al obtener productos");
  }
  return response.json();
};

// GET /api/productos/:id (detalle de un producto desde backend)
export const getProductoById = async (id: number): Promise<Producto | undefined> => {
  const response = await apiFetch(`/productos/${id}`);
  if (!response.ok) {
    throw new Error("Error al obtener producto");
  }
  const data = await response.json();
  return data;
};
