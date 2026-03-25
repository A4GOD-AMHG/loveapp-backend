# LoveApp Backend

Backend en Go para LoveApp. Expone la API, genera documentación Swagger y usa SQLite para el entorno local actual.

## Requisitos

- Go instalado
- `make`
- Un archivo `.env` válido

## Estructura Operativa

Estos paths se usan al correr o desplegar la app:

- `bin/`: binario compilado
- `data/`: base de datos SQLite local
- `logs/`: logs de ejecución
- `run/`: archivos de control del proceso
- `docs/`: Swagger generado

### Qué es el archivo `.pid`

El archivo `.pid` vive en `run/loveapp-backend.pid`. Solo guarda el PID del proceso que quedó corriendo en background.

Sirve para que `make deploy` pueda detectar si hay una instancia previa y reemplazarla.

Ese archivo no se versiona. Es runtime puro.

## Qué debe ir al repo y qué no

Sí debe ir al repo:

- código fuente
- `go.mod` y `go.sum`
- `Makefile`
- `README.md`
- `example.env`
- `docs/` si quieres dejar Swagger ya generado dentro del repo

No debe ir al repo:

- `.env`
- `bin/`
- `data/`
- `logs/`
- `run/`
- `*.pid`
- caches locales como `.gocache/`

El `.gitignore` ya quedó alineado con eso.

## Desarrollo Local

1. Instalar dependencias y herramientas:

```bash
make deps
```

2. Levantar en desarrollo:

```bash
make run
```

La app arranca usando el puerto definido en `.env`.

## Despliegue

### Deploy / Redeploy a producción

```bash
make deploy
```

Ese comando:

- regenera Swagger
- recompila el binario
- mata la instancia anterior si existe
- levanta la nueva
- deja logs en `logs/loveapp-backend.log`
- deja PID en `run/loveapp-backend.pid`

## Make Targets

- `make run`: ejecuta la app en desarrollo
- `make build`: genera Swagger y compila el binario
- `make deps`: descarga dependencias e instala `swag`
- `make test`: tests unitarios
- `make resetdb`: reinicia la base SQLite local
- `make clean`: borra binario, base local y PID
- `make deploy`: redeploy completo de producción

## Notas de Producción

- Si cambias variables en `.env`, vuelve a correr `make deploy`.
- Si quieres que el proceso sobreviva reinicios del servidor, lo correcto después es moverlo a `systemd`, `supervisor` o Docker. El flujo actual con `nohup` sirve para dejarlo desplegado ya, pero no es un process manager formal.
