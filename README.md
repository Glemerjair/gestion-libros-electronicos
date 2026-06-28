# Sistema de Gestión de Libros Electrónicos

## Datos del autor
- **Nombre:** Jair Glemer Armijos
- **Materia:** Programación Orientada a Objetos
- **Universidad:** Universidad Internacional del Ecuador (UIDE)
- **Fecha:** Junio 2026

## Objetivo
Desarrollar un sistema de gestión de libros electrónicos mediante programación funcional y orientada a objetos en Go, que permita administrar un catálogo digital de libros, gestionar el registro de usuarios y controlar el acceso a los recursos digitales disponibles.

## Tecnologías utilizadas
- **Lenguaje:** Go 1.21+
- **Framework web:** Gin
- **Base de datos:** MySQL 8.0
- **Frontend:** HTML, CSS
- **Otros:** godotenv, uuid, gin-contrib/sessions

## Módulos del sistema
- **Gestión de Libros:** agregar, listar, buscar, editar y eliminar libros
- **Gestión de Usuarios:** registrar, listar, buscar, editar y activar/desactivar usuarios
- **Gestión de Préstamos:** realizar y devolver préstamos con verificación de disponibilidad
- **Gestión de Categorías:** agregar, listar, asignar y eliminar categorías
- **Reportes:** libros más prestados, usuarios más activos y libros disponibles

## Servicios Web API
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/libros | Lista todos los libros |
| GET | /api/libros/buscar?q= | Busca libros por título o autor |
| GET | /api/usuarios | Lista todos los usuarios |
| GET | /api/categorias | Lista todas las categorías |
| GET | /api/prestamos | Lista préstamos activos |
| GET | /api/reportes/libros-mas-prestados | Ranking de libros más prestados |
| GET | /api/reportes/usuarios-activos | Usuarios con más actividad |
| GET | /api/reportes/libros-disponibles | Libros disponibles actualmente |

## Cómo ejecutar el proyecto

### Requisitos
- Go 1.21 o superior
- MySQL 8.0 o superior

### Pasos
1. Clonar el repositorio
2. Crear la base de datos ejecutando `database.sql`
3. Crear el archivo `.env` con las credenciales:
DB_HOST=127.0.0.1

DB_PORT=3306

DB_USER=root

DB_PASSWORD=1234

DB_NAME=gestion_libros

ADMIN_USER=admin

ADMIN_PASSWORD=1234
4. Instalar dependencias:
go mod tidy
5. Compilar y ejecutar:
go build -o app.exe main.go interfaces.go

.\app.exe
6. Abrir en el navegador: http://localhost:8080

## Estructura del proyecto
gestion-libros-electronicos/

├── main.go

├── interfaces.go

├── go.mod

├── go.sum

├── .env

├── database.sql

├── modules/

│   ├── libros/

│   ├── usuarios/

│   ├── prestamos/

│   ├── categorias/

│   └── reportes/

├── database/

├── templates/

└── static/