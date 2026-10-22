import { apiFetch } from "./http";

export interface CheckoutRequest {
  email: string;
  metodo_pago: string;
  items: Array<{
    producto_id: number;
    cantidad: number;
    precio_unitario: number;
  }>;
}

export interface CheckoutResponse {
  message: string;
  pedido: {
    id: number;
    email: string;
    total: number;
    metodo_pago: string;
    estado: string;
    fecha_creacion: string;
  };
}

// POST /api/checkout - Procesar pedido
export const createCheckout = async (request: CheckoutRequest): Promise<CheckoutResponse> => {
  const response = await apiFetch("/checkout", {
    method: "POST",
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || error.error || "Error al procesar el pedido");
  }

  return response.json();
};
