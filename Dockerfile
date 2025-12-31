# Dockerfile
# Multi-stage build para optimizar el tamaño de la imagen

# Etapa 1: Builder
FROM golang:1.23-alpine AS builder

# Instalar git y ca-certificates (necesarios para go get)
RUN apk add --no-cache git ca-certificates

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de módulos de Go
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
# CGO_ENABLED=0 para crear un binario estático
# -ldflags="-w -s" para reducir el tamaño del binario
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /api ./cmd/api

# Etapa 2: Runtime
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde el builder
COPY --from=builder /api .

# Exponer el puerto de la API
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./api"]
