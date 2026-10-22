// Declaramos el paquete models para agrupar las estructuras de datos de la aplicación.
package models

// LoginRequest define la estructura esperada para el cuerpo de la petición (JSON) al iniciar sesión.
type LoginRequest struct {
	// Email representa el correo del usuario; la etiqueta `json:"email"` indica cómo se mapea el JSON a esta variable.
	Email string `json:"email"`
	// Password representa la contraseña del usuario; se extrae del campo "password" del JSON entrante.
	Password string `json:"password"`
}

// RegisterRequest define la estructura para el registro de nuevos usuarios.
type RegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ref      string `json:"ref,omitempty"` // Código de referido opcional
}

// Producto define la estructura de los datos de un artículo en nuestro catálogo.
type Producto struct {
	// ID es el identificador único numérico del producto.
	ID int `json:"id"`
	// Nombre es la descripción en texto del producto.
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	// Precio es el costo del producto, almacenado como un número con decimales (float64).
	Precio float64 `json:"precio"`
	// Img es la URL o ruta que apunta a la fotografía o imagen del producto.
	Img       string   `json:"img"`
	Categoria string   `json:"categoria"`
	Galeria   []string `json:"galeria"`
	// Stock es la cantidad disponible en inventario.
	Stock int `json:"stock"`
}

// APIError estandariza la estructura de las respuestas de error HTTP.
type APIError struct {
	// Status es el código de estado HTTP (ej. 400, 401, 500).
	Status int `json:"status"`
	// Message es el mensaje de error descriptivo.
	Message string `json:"message"`
	// Details es información adicional opcional sobre el error.
	Details interface{} `json:"details,omitempty"`
}

// Usuario representa un usuario en el sistema multinivel.
type Usuario struct {
	// Email es el identificador único del usuario.
	Email string `json:"email"`
	// PasswordHash es el hash bcrypt del password.
	PasswordHash string `json:"-"`
	// Nombre es el nombre completo del usuario.
	Nombre string `json:"nombre"`
	// Rol define el tipo de usuario (admin, cliente, distribuidor).
	Rol    string `json:"rol"`
	Activo bool   `json:"activo,omitempty"`
	// ReferidoPor es el email del usuario que lo refirió (nil si es el raíz).
	ReferidoPor *string `json:"referido_por,omitempty"`
}

// NodoRed representa un nodo en el árbol de referidos multinivel.
type NodoRed struct {
	// Usuario contiene la información del usuario en este nodo.
	Usuario Usuario `json:"usuario"`
	Codigo  string  `json:"codigo,omitempty"`
	Ventas  float64 `json:"ventas"`
	// Hijos es la lista de nodos referidos por este usuario.
	Hijos []NodoRed `json:"hijos"`
}

// CheckoutRequest define la estructura esperada para el proceso de pago.
type CheckoutRequest struct {
	// Email del usuario que realiza la compra.
	Email string `json:"email"`
	// MetodoPago define el método de pago (tarjeta, paypal, transferencia).
	MetodoPago string `json:"metodo_pago"`
	// Items es el arreglo de productos en el carrito.
	Items []CarritoItem `json:"items"`
}

// CarritoItem representa un producto en el carrito de compras.
type CarritoItem struct {
	// ID del producto.
	ProductoID int `json:"producto_id"`
	// Cantidad de unidades.
	Cantidad int `json:"cantidad"`
	// Precio unitario al momento de la compra.
	PrecioUnitario float64 `json:"precio_unitario"`
}

// Pedido representa un pedido en el sistema.
type Pedido struct {
	// ID es el identificador único del pedido.
	ID        int `json:"id"`
	UsuarioID int `json:"usuario_id"`
	// Email del usuario que realizó el pedido.
	Email string `json:"email"`
	// Total es el monto total del pedido.
	Total float64 `json:"total"`
	// MetodoPago es el método de pago utilizado.
	MetodoPago string `json:"metodo_pago"`
	// Estado es el estado del pedido (pendiente, completado, cancelado).
	Estado string `json:"estado"`
	// FechaCreacion es la fecha en que se creó el pedido.
	FechaCreacion string `json:"fecha_creacion"`
}

// PedidoItem representa un ítem dentro de un pedido.
type PedidoItem struct {
	// ID es el identificador único del ítem.
	ID int `json:"id"`
	// PedidoID es el ID del pedido al que pertenece.
	PedidoID int `json:"pedido_id"`
	// ProductoID es el ID del producto.
	ProductoID int `json:"producto_id"`
	// Cantidad de unidades compradas.
	Cantidad int `json:"cantidad"`
	// PrecioUnitario al momento de la compra.
	PrecioUnitario float64 `json:"precio_unitario"`
	// Subtotal es cantidad * precio_unitario.
	Subtotal float64 `json:"subtotal"`
}
