# Compila con CGO (go-sqlite3 lo necesita) y corre en una imagen mínima con glibc.
FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /out/loveapp-backend . && mkdir -p /out/data

# distroless/base trae glibc y certificados CA (Firebase usa HTTPS).
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/loveapp-backend /app/loveapp-backend
# /data debe ser del usuario nonroot para que SQLite pueda escribir en el volumen.
COPY --from=build --chown=65532:65532 /out/data /data
WORKDIR /app
ENV DB_PATH=/data/loveapp.db SERVER_PORT=4418
EXPOSE 4418
ENTRYPOINT ["/app/loveapp-backend"]
