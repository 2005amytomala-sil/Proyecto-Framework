import { productosMock, type Producto } from "../data/productos";

// Simula GET /api/productos (lista completa)
export const getProductos = (): Promise<Producto[]> => {
  return new Promise((resolve) => {
    setTimeout(() => resolve(productosMock), 400);
  });
};

// Simula GET /api/productos/:id (detalle de un producto)
export const getProductoById = (id: number): Promise<Producto | undefined> => {
  return new Promise((resolve) => {
    setTimeout(() => resolve(productosMock.find((p) => p.id === id)), 300);
  });
};