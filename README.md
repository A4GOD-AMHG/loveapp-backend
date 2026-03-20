# LoveApp Backend 💕

Este es el backend de **LoveApp**, una aplicación móvil privada diseñada para mi novia y para mí. El objetivo es centralizar nuestra comunicación, organización y los momentos especiales que compartimos.

## Tech Stack

<div style="display: flex; align-items: center; gap: 10px;">
  <img src="https://raw.githubusercontent.com/A4GOD-AMHG/Utils-for-repos/main/icons/go/go-original-wordmark.svg" alt="Go" width="65" height="65" />
  <img src="https://raw.githubusercontent.com/A4GOD-AMHG/Utils-for-repos/main/icons/sqlite/sqlite-original.svg" alt="SQLite" width="65" height="65" />
  <img src="https://raw.githubusercontent.com/A4GOD-AMHG/Utils-for-repos/main/icons/swagger/swagger-original.svg" alt="Swagger" width="65" height="65" />
  <img src="https://raw.githubusercontent.com/A4GOD-AMHG/Utils-for-repos/main/icons/socketio/socketio-original.svg" alt="socketio" width="65" height="65" />
</div>

## Características

-   💬 **Chat Privado**: Comunicación en tiempo real y segura entre nosotros.
-   ✅ **Lista de Tareas Compartida**: Para organizar nuestras metas y pendientes juntos.
-   ❤️ **Contador de Días Juntos**: Un espacio para llevar la cuenta de nuestro tiempo juntos.
-   🔔 **Notificaciones y Alarmas**: Recordatorios en tiempo real para fechas y eventos importantes.
-   🔐 **Autenticación Segura con JWT**: Protegiendo nuestra información.
-   📚 **Documentación Interactiva con Swagger**: Para probar y entender la API de forma sencilla.

## Guía de Inicio Rápido

Sigue estos pasos para levantar el backend en tu entorno local.

1.  **Clona el repositorio**:

    ```bash
    git clone https://github.com/A4GOD-AMHG/LoveApp-Backend.git
    cd LoveApp-Backend
    ```

2.  **Configura las variables de entorno**:

    Crea un archivo `.env` a partir del ejemplo.

    ```bash
    make env
    ```

    Puedes editar el archivo `.env` si necesitas cambiar el puerto o el secreto del JWT.

3.  **Instala las dependencias**:

    Este comando descargará los módulos de Go y las herramientas necesarias como `swag`.

    ```bash
    make setup
    ```

4.  **¡Lanza la aplicación!**:

    ```bash
    make run
    ```

    El servidor se iniciará en `http://localhost:8080`.

## Comandos Útiles del Makefile

-   `make run`: Inicia la aplicación en modo desarrollo.
-   `make build`: Genera la documentación de Swagger y compila el binario para producción.
-   `make swagger`: Regenera la documentación de Swagger manualmente.
-   `make setup`: Instala todas las dependencias del proyecto.
-   `make clean`: Elimina los binarios compilados y la base de datos local.

## Frontend

Este backend está diseñado para ser consumido por la aplicación móvil privada de LoveApp.

## Autor

-   Alexis Manuel Hurtado García ([@A4GOD-AMHG](https://github.com/A4GOD-AMHG))
