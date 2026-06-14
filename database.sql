CREATE DATABASE IF NOT EXISTS gestion_libros;
USE gestion_libros;

CREATE TABLE categorias (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT
);

CREATE TABLE libros (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titulo VARCHAR(200) NOT NULL,
    autor VARCHAR(150) NOT NULL,
    anio_publicacion INT,
    enlace_drive VARCHAR(500),
    disponible BOOLEAN DEFAULT TRUE,
    categoria_id INT,
    FOREIGN KEY (categoria_id) REFERENCES categorias(id)
);

CREATE TABLE usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    fecha_registro DATETIME DEFAULT CURRENT_TIMESTAMP,
    activo BOOLEAN DEFAULT TRUE
);

CREATE TABLE prestamos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    libro_id INT NOT NULL,
    usuario_id INT NOT NULL,
    fecha_prestamo DATETIME DEFAULT CURRENT_TIMESTAMP,
    fecha_devolucion DATETIME,
    activo BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (libro_id) REFERENCES libros(id),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);