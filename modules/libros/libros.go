package libros

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Libro representa la estructura de un libro en el sistema
type Libro struct {
	ID              int
	Titulo          string
	Autor           string
	AnioPublicacion int
	EnlaceDrive     string
	Disponible      bool
	CategoriaID     int
}

// agregarLibro registra un nuevo libro en el catálogo
func AgregarLibro(db *sql.DB, titulo, autor string, anio int, enlaceDrive string, categoriaID int) error {
	if titulo == "" || autor == "" {
		return errors.New("el título y el autor no pueden estar vacíos")
	}

	query := "INSERT INTO libros (titulo, autor, anio_publicacion, enlace_drive, disponible, categoria_id) VALUES (?, ?, ?, ?, TRUE, ?)"
	_, err := db.Exec(query, titulo, autor, anio, enlaceDrive, categoriaID)
	if err != nil {
		return fmt.Errorf("error al agregar libro: %w", err)
	}

	fmt.Println("Libro agregado exitosamente")
	return nil
}

// listarLibros retorna todos los libros disponibles
func ListarLibros(db *sql.DB) ([]Libro, error) {
	query := "SELECT id, titulo, autor, anio_publicacion, enlace_drive, disponible, categoria_id FROM libros"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al listar libros: %w", err)
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.AnioPublicacion, &l.EnlaceDrive, &l.Disponible, &l.CategoriaID)
		if err != nil {
			return nil, fmt.Errorf("error al leer libro: %w", err)
		}
		libros = append(libros, l)
	}
	return libros, nil
}

// buscarLibro filtra por título o autor
func BuscarLibro(db *sql.DB, criterio string) ([]Libro, error) {
	if criterio == "" {
		return nil, errors.New("el criterio de búsqueda no puede estar vacío")
	}

	criterio = "%" + strings.ToLower(criterio) + "%"
	query := `SELECT id, titulo, autor, anio_publicacion, enlace_drive, disponible, categoria_id 
			  FROM libros 
			  WHERE LOWER(titulo) LIKE ? OR LOWER(autor) LIKE ?`

	rows, err := db.Query(query, criterio, criterio)
	if err != nil {
		return nil, fmt.Errorf("error al buscar libro: %w", err)
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var l Libro
		err := rows.Scan(&l.ID, &l.Titulo, &l.Autor, &l.AnioPublicacion, &l.EnlaceDrive, &l.Disponible, &l.CategoriaID)
		if err != nil {
			return nil, fmt.Errorf("error al leer libro: %w", err)
		}
		libros = append(libros, l)
	}
	return libros, nil
}

// actualizarLibro modifica los datos de un libro existente
func ActualizarLibro(db *sql.DB, id int, titulo, autor string, anio int, enlaceDrive string) error {
	if id <= 0 {
		return errors.New("el ID del libro no es válido")
	}

	query := "UPDATE libros SET titulo = ?, autor = ?, anio_publicacion = ?, enlace_drive = ? WHERE id = ?"
	_, err := db.Exec(query, titulo, autor, anio, enlaceDrive, id)
	if err != nil {
		return fmt.Errorf("error al actualizar libro: %w", err)
	}

	fmt.Println("Libro actualizado exitosamente")
	return nil
}

// eliminarLibro elimina un libro del catálogo
func EliminarLibro(db *sql.DB, id int) error {
	if id <= 0 {
		return errors.New("el ID del libro no es válido")
	}

	query := "DELETE FROM libros WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar libro: %w", err)
	}

	fmt.Println("Libro eliminado exitosamente")
	return nil
}

// verLibro retorna el enlace de Google Drive del libro
func VerLibro(db *sql.DB, id int) (string, error) {
	if id <= 0 {
		return "", errors.New("el ID del libro no es válido")
	}

	var enlace string
	query := "SELECT enlace_drive FROM libros WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&enlace)
	if err != nil {
		return "", fmt.Errorf("error al obtener enlace del libro: %w", err)
	}

	return enlace, nil
}

// Encapsulación

// SetTitulo establece el título del libro validando que no esté vacío
func (l *Libro) SetTitulo(titulo string) error {
	if titulo == "" {
		return errors.New("el título no puede estar vacío")
	}
	l.Titulo = titulo
	return nil
}

// SetAutor establece el autor del libro validando que no esté vacío
func (l *Libro) SetAutor(autor string) error {
	if autor == "" {
		return errors.New("el autor no puede estar vacío")
	}
	l.Autor = autor
	return nil
}

// SetEnlaceDrive establece el enlace de Google Drive del libro
func (l *Libro) SetEnlaceDrive(enlace string) error {
	if enlace == "" {
		return errors.New("el enlace no puede estar vacío")
	}
	l.EnlaceDrive = enlace
	return nil
}

// SetDisponible establece la disponibilidad del libro
func (l *Libro) SetDisponible(disponible bool) {
	l.Disponible = disponible
}
