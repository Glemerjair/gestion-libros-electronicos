package main

import "database/sql"

// InterfazLibros define las operaciones del módulo de libros
type InterfazLibros interface {
	AgregarLibro(db *sql.DB, titulo, autor, isbn string, anio int, enlaceDrive string, categoriaID int) error
	ListarLibros(db *sql.DB) error
	BuscarLibro(db *sql.DB, criterio string) error
	ActualizarLibro(db *sql.DB, id int, titulo, autor, anio int, enlaceDrive string) error
	EliminarLibro(db *sql.DB, id int) error
	VerLibro(db *sql.DB, id int) (string, error)
}

// InterfazUsuarios define las operaciones del módulo de usuarios
type InterfazUsuarios interface {
	RegistrarUsuario(db *sql.DB, nombre, email string) error
	ListarUsuarios(db *sql.DB) error
	BuscarUsuario(db *sql.DB, criterio string) error
	ActualizarUsuario(db *sql.DB, id int, nombre, email string) error
	EliminarUsuario(db *sql.DB, id int) error
}

// InterfazPrestamos define las operaciones del módulo de préstamos
type InterfazPrestamos interface {
	RealizarPrestamo(db *sql.DB, libroID, usuarioID int) error
	DevolverLibro(db *sql.DB, prestamoID int) error
	ListarPrestamos(db *sql.DB) error
	HistorialPrestamos(db *sql.DB, usuarioID int) error
	VerificarDisponibilidad(db *sql.DB, libroID int) (bool, error)
}

// InterfazCategorias define las operaciones del módulo de categorías
type InterfazCategorias interface {
	AgregarCategoria(db *sql.DB, nombre, descripcion string) error
	ListarCategorias(db *sql.DB) error
	AsignarCategoria(db *sql.DB, libroID, categoriaID int) error
	EliminarCategoria(db *sql.DB, id int) error
}

// InterfazReportes define las operaciones del módulo de reportes
type InterfazReportes interface {
	LibrosMasPrestados(db *sql.DB) error
	UsuariosActivos(db *sql.DB) error
	LibrosDisponibles(db *sql.DB) error
}
