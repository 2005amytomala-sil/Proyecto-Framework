-- Migración idempotente para volúmenes creados por versiones anteriores.
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS referido_id INTEGER REFERENCES referidos(id);
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS creado_en TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'usuarios' AND column_name = 'password'
    ) THEN
        EXECUTE $q$UPDATE usuarios SET password_hash = COALESCE(NULLIF(password_hash, ''), password) WHERE password_hash IS NULL OR password_hash = ''$q$;
    END IF;
END $$;
ALTER TABLE usuarios ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE productos ADD COLUMN IF NOT EXISTS activo BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS usuario_id INTEGER REFERENCES usuarios(id);
ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS estado VARCHAR(50) NOT NULL DEFAULT 'pendiente';
ALTER TABLE referidos ADD COLUMN IF NOT EXISTS codigo VARCHAR(50);
UPDATE referidos SET codigo = UPPER(SUBSTRING(REPLACE(email,'@','') FROM 1 FOR 12)) WHERE codigo IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_referidos_codigo_unique ON referidos(codigo);
CREATE TABLE IF NOT EXISTS liquidaciones (id SERIAL PRIMARY KEY, referido_id INTEGER NOT NULL REFERENCES referidos(id), total NUMERIC(12,2) NOT NULL, estado VARCHAR(50) NOT NULL DEFAULT 'pendiente', creado_en TIMESTAMP NOT NULL DEFAULT NOW(), pagada_en TIMESTAMP);
CREATE TABLE IF NOT EXISTS comisiones (id SERIAL PRIMARY KEY, pedido_id INTEGER NOT NULL REFERENCES pedidos(id), referido_id INTEGER NOT NULL REFERENCES referidos(id), nivel INTEGER NOT NULL, porcentaje NUMERIC(5,2) NOT NULL, monto NUMERIC(12,2) NOT NULL, estado VARCHAR(50) NOT NULL DEFAULT 'pendiente', liquidacion_id INTEGER REFERENCES liquidaciones(id), creado_en TIMESTAMP NOT NULL DEFAULT NOW());
CREATE INDEX IF NOT EXISTS idx_comisiones_referido_estado ON comisiones(referido_id,estado);
