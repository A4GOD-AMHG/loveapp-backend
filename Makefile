.PHONY: help run build swagger deps clean setup env reset-db

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BINARY="$(BIN_DIR)/$(APP_NAME)"

help: ## Muestra esta ayuda
	@echo 'Uso: make <comando>'
	@echo ''
	@echo 'Comandos:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Ejecutar la aplicación en modo desarrollo
	@echo "🚀 Iniciando la aplicación en http://localhost:$(SERVER_PORT)..."
	@mkdir -p ./data
	go run .

build: swagger ## Compilar el binario para producción
	@echo "🔨 Compilando binario..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BINARY) .
	@echo "✅ Compilación finalizada: $(BINARY)"

swagger: ## Generar la documentación de Swagger
	@echo "📚 Generando documentación de Swagger..."
	swag init -g main.go --output ./docs
	@echo "✅ Documentación de Swagger generada."

setup: env deps ## Realizar la configuración inicial del proyecto
	@mkdir -p ./data
	@echo "✅ ¡Proyecto configurado! Ahora puedes usar 'make run'."

deps: ## Instalar/actualizar dependencias y herramientas
	@echo "📦 Instalando dependencias..."
	go mod tidy
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Dependencias actualizadas."

clean: ## Limpiar binarios compilados y la base de datos
	@echo "🧹 Limpiando..."
	rm -f $(BINARY)
	rm -f $(DB_PATH)
	@echo "✅ Limpieza finalizada."

reset-db: ## Vaciar la base de datos, correr migraciones y sembrar usuarios iniciales
	@echo "🗃️ Reiniciando base de datos en $(DB_PATH)..."
	GOCACHE=/tmp/go-build go run ./cmd/resetdb
	@echo "✅ Base de datos reiniciada."

test: ## Ejecutar todos los tests unitarios con salida detallada
	@echo "🧪 Ejecutando tests unitarios..."
	go test -v -count=1 ./...
	@echo "✅ Tests finalizados."

env: ## Crear .env desde example.env si no existe
	@if [ ! -f .env ]; then \
		cp example.env .env; \
		echo "✅ Archivo .env creado."; \
	else \
		echo "⚠️  El archivo .env ya existe."; \
	fi
