# LoveApp Backend Makefile

.PHONY: help up down restart logs clean build dev

# Default target
help: ## Mostrar ayuda
	@echo '╔════════════════════════════════════════╗'
	@echo '║     LoveApp Backend - Comandos         ║'
	@echo '╚════════════════════════════════════════╝'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================
# Docker - Comandos principales
# ============================================

up: ## Levantar Backend
	docker compose up -d
	@echo "✅ Servicios levantados!"
	@echo "🔗 Backend: http://localhost:8080"
	@echo "📚 Swagger: http://localhost:8080/swagger/index.html"
	@echo "🏥 Health: http://localhost:8080/health"

down: ## Detener todos los servicios
	docker compose down
	@echo "🛑 Servicios detenidos"

restart: ## Reiniciar servicios
	docker compose restart
	@echo "🔄 Servicios reiniciados"

rebuild: ## Reconstruir y levantar
	@echo "🔨 Reconstruyendo servicios..."
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	@echo "✅ Servicios reconstruidos y levantados"

build-fast: ## Construir usando caché (rápido)
	@echo "⚡ Construyendo con caché..."
	docker compose build
	@echo "✅ Construcción completada"

up-build: ## Construir y levantar (usa caché)
	@echo "🚀 Construyendo y levantando servicios..."
	docker compose up -d --build
	@echo "✅ Servicios levantados!"
	@echo "🔗 Backend: http://localhost:8080"
	@echo "📚 Swagger: http://localhost:8080/swagger/index.html"
	@echo "🏥 Health: http://localhost:8080/health"

# ============================================
# Logs
# ============================================

logs: ## Ver logs del backend
	docker compose logs -f backend

logs-all: ## Ver todos los logs
	docker compose logs -f

# ============================================
# Limpieza
# ============================================

clean: ## Limpiar contenedores y volúmenes
	docker compose down -v
	@echo "🧹 Contenedores y volúmenes eliminados"

clean-all: ## Limpieza completa (incluye imágenes)
	docker compose down -v --rmi all
	docker system prune -f
	@echo "🧹 Limpieza completa realizada"

clean-db: ## Limpiar solo la base de datos SQLite
	rm -f ./data/loveapp.db
	@echo "🧹 Base de datos SQLite eliminada"

# ============================================
# Desarrollo local (sin Docker)
# ============================================

dev: ## Ejecutar en modo desarrollo local
	@echo "🚀 Iniciando en modo desarrollo..."
	@mkdir -p ./data
	go run .

build: ## Compilar aplicación
	@echo "🔨 Compilando..."
	swag init -g main.go --output ./docs
	go build -o bin/loveapp-backend .
	@echo "✅ Compilación exitosa: bin/loveapp-backend"

deps: ## Instalar dependencias
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Dependencias instaladas"

# ============================================
# Utilidades
# ============================================

health: ## Verificar salud del servicio
	@curl -f http://localhost:8080/health && echo "\n✅ Servicio saludable" || echo "\n❌ Servicio no disponible"

swagger: ## Abrir Swagger UI
	@echo "📚 Abriendo Swagger UI..."
	@which xdg-open >/dev/null && xdg-open http://localhost:8080/swagger/index.html || echo "Abre: http://localhost:8080/swagger/index.html"

env: ## Crear archivo .env desde example.env
	@if [ ! -f .env ]; then \
		cp example.env .env; \
		echo "✅ Archivo .env creado"; \
	else \
		echo "⚠️  .env ya existe"; \
	fi

# ============================================
# Setup inicial
# ============================================

setup: env deps ## Setup completo del proyecto
	@mkdir -p ./data
	@echo "✅ Setup completo!"
	@echo "Ejecuta: make up"