.PHONY: help run deps build clean resetdb test deploy

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

SERVICE_NAME ?= loveapp-backend
DISPLAY_NAME ?= $(if $(APP_NAME),$(APP_NAME),$(SERVICE_NAME))
BIN_DIR ?= ./bin
LOG_DIR ?= ./logs
RUN_DIR ?= ./run
DB_PATH ?= ./data/loveapp.db
SERVER_PORT ?= 4418
BINARY="$(BIN_DIR)/$(SERVICE_NAME)"
PID_FILE="$(RUN_DIR)/$(SERVICE_NAME).pid"
LOG_FILE="$(LOG_DIR)/$(SERVICE_NAME).log"

help: ## Muestra esta ayuda
	@echo 'Uso: make <comando>'
	@echo ''
	@echo 'Comandos:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Ejecutar la aplicación en modo desarrollo
	@echo "🚀 Iniciando la aplicación en http://localhost:$(SERVER_PORT)..."
	@mkdir -p ./data
	go run . --host=0.0.0.0

deps: ## Instalar/actualizar dependencias y herramientas
	@echo "📦 Instalando dependencias..."
	go mod tidy
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Dependencias actualizadas."

build: ## Compilar el binario para producción (incluye Swagger)
	@echo "🔨 Compilando binario..."
	@echo "📚 Generando documentación Swagger..."
	swag init -g main.go --output ./docs
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BINARY) .
	@echo "✅ Compilación finalizada: $(BINARY)"

clean: ## Limpiar binarios compilados y la base de datos
	@echo "🧹 Limpiando..."
	rm -f $(BINARY)
	rm -f $(DB_PATH)
	rm -f $(PID_FILE)
	@echo "✅ Limpieza finalizada."

resetdb: ## Vaciar la base de datos, correr migraciones y sembrar usuarios iniciales
	@echo "🗃️ Reiniciando base de datos en $(DB_PATH)..."
	@mkdir -p ./data
	rm -f $(DB_PATH)
	LOVEAPP_RESET_ONLY=1 GOCACHE=/tmp/go-build go run .
	@echo "✅ Base de datos reiniciada."

test: ## Ejecutar todos los tests unitarios con salida detallada
	@echo "🧪 Ejecutando tests unitarios..."
	GOCACHE=$(CURDIR)/.gocache go test -v -count=1 ./...
	@echo "✅ Tests finalizados."

deploy: build ## Recompilar y desplegar producción (reemplaza instancia previa si existe)
	@echo "🚀 Desplegando $(DISPLAY_NAME) en producción..."
	@mkdir -p $(LOG_DIR) $(RUN_DIR) ./data
	@if [ -f $(PID_FILE) ] && kill -0 $$(cat $(PID_FILE)) 2>/dev/null; then \
		echo "🛑 Deteniendo instancia actual: $$(cat $(PID_FILE))"; \
		kill $$(cat $(PID_FILE)); \
		sleep 1; \
	fi
	@rm -f $(PID_FILE)
	@nohup $(BINARY) >> $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)
	@sleep 1
	@if kill -0 $$(cat $(PID_FILE)) 2>/dev/null; then \
		echo "✅ Producción desplegada. PID: $$(cat $(PID_FILE))"; \
		echo "📄 Logs: $(LOG_FILE)"; \
	else \
		echo "❌ El proceso no inició correctamente. Revisa logs en $(LOG_FILE)"; \
		exit 1; \
	fi
