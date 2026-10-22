# Usamos la imagen oficial de Go como base.
FROM golang:1.26-alpine AS builder

# Establecemos el directorio de trabajo.
WORKDIR /app

# Copiamos los archivos go.mod y go.sum.
COPY go.mod go.sum ./

# Descargamos las dependencias.
RUN go mod download

# Copiamos el resto del código fuente.
COPY . .

# Construimos la aplicación.
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Usamos una imagen alpine más pequeña para ejecutar la aplicación.
FROM alpine:latest

# Instalamos ca-certificates para las peticiones HTTPS.
RUN apk --no-cache add ca-certificates

RUN adduser -D -u 10001 appuser
WORKDIR /app

# Copiamos el binario desde el builder.
COPY --from=builder /app/main .
USER appuser

# Exponemos el puerto 3000.
EXPOSE 3000

# Ejecutamos la aplicación.
CMD ["./main"]
