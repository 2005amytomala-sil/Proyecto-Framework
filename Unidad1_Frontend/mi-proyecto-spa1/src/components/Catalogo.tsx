import {useCart} from '../context/CartContext';
import { useState, useEffect } from "react";

// Definimos la estructura del producto coincidente con el modelo de Go
interface Producto {
  id: number;
  nombre: string;
  precio: number;
  img: string;
}
const Catalogo = () => {
  const { addToCart } = useCart();
  const [productos, setProductos] = useState<Producto[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    const fetchProductos = async () => {
      try {
        const response = await fetch("http://localhost:3000/api/productos");

        if (!response.ok) {
          throw new Error("No se pudo cargar la lista de productos");
        }

        const data: Producto[] = await response.json();
        setProductos(data);
      } catch (err: any) {
        setError(err.message || "Error al conectar con la API de Go");
      } finally {
        setLoading(false);
      }
    };

    fetchProductos();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-slate-500 font-medium animate-pulse text-lg">
          Cargando productos desde el servidor...
        </p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg">
        <p className="font-semibold">Error al cargar el catálogo:</p>
        <p className="text-sm">{error}</p>
        <p className="text-xs text-slate-500 mt-2">
          Verifica que el backend de Go esté corriendo en el puerto 3000 y tenga CORS habilitado.
        </p>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">
        Catálogo de Productos
      </h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {productos.map((prod) => (
          <div
            key={prod.id}
            className="bg-white rounded-lg overflow-hidden border border-slate-200 shadow-sm flex flex-col hover:shadow-md transition"
          >
            <img
              src={prod.img}
              alt={prod.nombre}
              className="w-full h-40 object-cover"
            />
            <div className="p-4 flex flex-col flex-1">
              <h3 className="font-semibold text-slate-700">{prod.nombre}</h3>
              <p className="text-indigo-600 font-bold mt-2 mb-4">
                ${prod.precio.toFixed(2)}
              </p>

              <button
                onClick={() => addToCart(prod)}
                className="mt-auto w-full bg-slate-900 text-white py-2 rounded text-sm hover:bg-indigo-600 transition"
              >
                Añadir al Carrito
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
export default Catalogo;