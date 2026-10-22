-- Esquema inicial para instalaciones nuevas. No borra datos existentes.

CREATE TABLE IF NOT EXISTS referidos (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    codigo VARCHAR(50) UNIQUE NOT NULL,
    nivel INTEGER NOT NULL DEFAULT 0,
    ventas NUMERIC(12, 2) NOT NULL DEFAULT 0,
    parent_id INTEGER REFERENCES referidos(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL,
    referido_por VARCHAR(255),
    referido_id INTEGER REFERENCES referidos(id),
    activo BOOLEAN NOT NULL DEFAULT true,
    saldo_comisiones NUMERIC(12, 2) NOT NULL DEFAULT 0,
    creado_en TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS productos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT NOT NULL DEFAULT '',
    precio NUMERIC(10, 2) NOT NULL CHECK (precio >= 0),
    categoria VARCHAR(100) NOT NULL DEFAULT '',
    img VARCHAR(500) NOT NULL DEFAULT '',
    galeria JSONB NOT NULL DEFAULT '[]'::jsonb,
    stock INTEGER NOT NULL DEFAULT 100 CHECK (stock >= 0),
    activo BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pedidos (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    usuario_id INTEGER REFERENCES usuarios(id),
    total NUMERIC(12, 2) NOT NULL CHECK (total >= 0),
    metodo_pago VARCHAR(50) NOT NULL,
    estado VARCHAR(50) NOT NULL DEFAULT 'pendiente',
    fecha_creacion TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pedido_items (
    id SERIAL PRIMARY KEY,
    pedido_id INTEGER NOT NULL REFERENCES pedidos(id),
    producto_id INTEGER NOT NULL REFERENCES productos(id),
    cantidad INTEGER NOT NULL CHECK (cantidad > 0),
    precio_unitario NUMERIC(10, 2) NOT NULL CHECK (precio_unitario >= 0),
    subtotal NUMERIC(12, 2) NOT NULL CHECK (subtotal >= 0)
);

CREATE TABLE IF NOT EXISTS liquidaciones (
    id SERIAL PRIMARY KEY,
    referido_id INTEGER NOT NULL REFERENCES referidos(id),
    total NUMERIC(12, 2) NOT NULL CHECK (total >= 0),
    estado VARCHAR(50) NOT NULL DEFAULT 'pendiente',
    creado_en TIMESTAMP NOT NULL DEFAULT NOW(),
    pagada_en TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comisiones (
    id SERIAL PRIMARY KEY,
    pedido_id INTEGER NOT NULL REFERENCES pedidos(id),
    referido_id INTEGER NOT NULL REFERENCES referidos(id),
    nivel INTEGER NOT NULL CHECK (nivel BETWEEN 1 AND 3),
    porcentaje NUMERIC(5, 2) NOT NULL,
    monto NUMERIC(12, 2) NOT NULL CHECK (monto >= 0),
    estado VARCHAR(50) NOT NULL DEFAULT 'pendiente',
    liquidacion_id INTEGER REFERENCES liquidaciones(id),
    creado_en TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_referidos_parent_id ON referidos(parent_id);
CREATE INDEX IF NOT EXISTS idx_referidos_email ON referidos(email);
CREATE INDEX IF NOT EXISTS idx_referidos_codigo ON referidos(codigo);
CREATE INDEX IF NOT EXISTS idx_usuarios_email ON usuarios(email);
CREATE INDEX IF NOT EXISTS idx_usuarios_referido_id ON usuarios(referido_id);
CREATE INDEX IF NOT EXISTS idx_productos_categoria ON productos(categoria);
CREATE INDEX IF NOT EXISTS idx_productos_activo ON productos(activo);
CREATE INDEX IF NOT EXISTS idx_pedidos_email ON pedidos(email);
CREATE INDEX IF NOT EXISTS idx_pedidos_usuario_id ON pedidos(usuario_id);
CREATE INDEX IF NOT EXISTS idx_pedidos_estado ON pedidos(estado);
CREATE INDEX IF NOT EXISTS idx_pedido_items_pedido_id ON pedido_items(pedido_id);
CREATE INDEX IF NOT EXISTS idx_comisiones_pedido_id ON comisiones(pedido_id);
CREATE INDEX IF NOT EXISTS idx_comisiones_referido_id ON comisiones(referido_id);
CREATE INDEX IF NOT EXISTS idx_comisiones_estado ON comisiones(estado);
CREATE INDEX IF NOT EXISTS idx_liquidaciones_referido_id ON liquidaciones(referido_id);

-- Seed multinivel: cada fila queda relacionada con su padre por ID.
INSERT INTO referidos (email, nombre, codigo, nivel, parent_id) VALUES
    ('admin@upse.edu.ec', 'Admin Principal', 'ADMIN', 0, NULL),
    ('usuario1@upse.edu.ec', 'Usuario 1', 'USUARIO1', 1, 1),
    ('usuario2@upse.edu.ec', 'Usuario 2', 'USUARIO2', 1, 1),
    ('usuario3@upse.edu.ec', 'Usuario 3', 'USUARIO3', 2, 2),
    ('usuario4@upse.edu.ec', 'Usuario 4', 'USUARIO4', 2, 2),
    ('usuario5@upse.edu.ec', 'Usuario 5', 'USUARIO5', 2, 3)
ON CONFLICT (email) DO NOTHING;

-- Contraseña demo: 123456. Solo se persiste su hash bcrypt.
INSERT INTO usuarios (email, password_hash, nombre, rol, referido_id) VALUES
    ('admin@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Admin Principal', 'admin', 1),
    ('cliente@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Cliente Demo', 'cliente', NULL),
    ('usuario1@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Usuario 1', 'cliente', 2),
    ('usuario2@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Usuario 2', 'cliente', 3),
    ('usuario3@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Usuario 3', 'distribuidor', 4),
    ('usuario4@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Usuario 4', 'cliente', 5),
    ('usuario5@upse.edu.ec', '$2a$10$sT4kxI/54nUJ1awYcbRhAunor716D80UQyR/dLnbRsNAVYcgb8Ahq', 'Usuario 5', 'cliente', 6)
ON CONFLICT (email) DO NOTHING;

INSERT INTO productos (nombre, descripcion, precio, categoria, img, galeria, stock) VALUES
    ('Serum Revitalizante', 'Serum antioxidante para piel radiante', 45.00, 'Serum', 'https://picsum.photos/seed/serum/150', '["https://picsum.photos/seed/serum1/300", "https://picsum.photos/seed/serum2/300"]', 100),
    ('Crema Hidratante Pro', 'Crema hidratante de 24 horas', 32.50, 'Crema', 'https://picsum.photos/seed/crema/150', '["https://picsum.photos/seed/crema1/300", "https://picsum.photos/seed/crema2/300"]', 100),
    ('Tónico Purificante', 'Tónico equilibrante para todo tipo de piel', 28.00, 'Tónico', 'https://picsum.photos/seed/tonico/150', '["https://picsum.photos/seed/tonico1/300", "https://picsum.photos/seed/tonico2/300"]', 100),
    ('Mascarilla Nocturna', 'Mascarilla reparadora para uso nocturno', 50.00, 'Mascarilla', 'https://picsum.photos/seed/mascarilla/150', '["https://picsum.photos/seed/mascarilla1/300", "https://picsum.photos/seed/mascarilla2/300"]', 100)
ON CONFLICT DO NOTHING;
