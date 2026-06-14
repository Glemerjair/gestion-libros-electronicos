package categorias

import (
	"database/sql"
	"errors"
	"fmt"
)

// Categoria representa la estructura de una categoría en el sistema
type Categoria struct {
	ID          int
	Nombre      string
	Descripcion string
}

// agregarCategoria registra una nueva categoría en la base de datos
func AgregarCategoria(db *sql.DB, nombre, descripcion string) error {
	if nombre == "" {
		return errors.New("el nombre de la categoría no puede estar vacío")
	}

	query := "INSERT INTO categorias (nombre, descripcion) VALUES (?, ?)"
	_, err := db.Exec(query, nombre, descripcion)
	if err != nil {
		return fmt.Errorf("error al agregar categoría: %w", err)
	}

	fmt.Println("Categoría agregada exitosamente")
	return nil
}

// listarCategorias retorna todas las categorías registradas
func ListarCategorias(db *sql.DB) ([]Categoria, error) {
	query := "SELECT id, nombre, descripcion FROM categorias"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al listar categorías: %w", err)
	}
	defer rows.Close()

	var categorias []Categoria
	for rows.Next() {
		var c Categoria
		err := rows.Scan(&c.ID, &c.Nombre, &c.Descripcion)
		if err != nil {
			return nil, fmt.Errorf("error al leer categoría: %w", err)
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}

// asignarCategoria vincula un libro a una categoría
func AsignarCategoria(db *sql.DB, libroID, categoriaID int) error {
	query := "UPDATE libros SET categoria_id = ? WHERE id = ?"
	_, err := db.Exec(query, categoriaID, libroID)
	if err != nil {
		return fmt.Errorf("error al asignar categoría: %w", err)
	}

	fmt.Println("Categoría asignada exitosamente")
	return nil
}

// eliminarCategoria elimina una categoría del sistema
func EliminarCategoria(db *sql.DB, id int) error {
	if id <= 0 {
		return errors.New("el ID de la categoría no es válido")
	}

	query := "DELETE FROM categorias WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar categoría: %w", err)
	}

	fmt.Println("Categoría eliminada exitosamente")
	return nil
}

// Encapsulación

// SetNombre establece el nombre de la categoría validando que no esté vacío
func (c *Categoria) SetNombre(nombre string) error {
	if nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	c.Nombre = nombre
	return nil
}

// SetDescripcion establece la descripción de la categoría
func (c *Categoria) SetDescripcion(descripcion string) error {
	if descripcion == "" {
		return errors.New("la descripción no puede estar vacía")
	}
	c.Descripcion = descripcion
	return nil
}
