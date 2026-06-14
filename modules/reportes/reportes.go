package reportes

import (
	"database/sql"
	"fmt"
)

// ReporteLibro representa los datos de un libro en el reporte
type ReporteLibro struct {
	ID     int
	Titulo string
	Autor  string
	Total  int
}

// ReporteUsuario representa los datos de un usuario en el reporte
type ReporteUsuario struct {
	ID     int
	Nombre string
	Email  string
	Total  int
}

// librosMasPrestados retorna el ranking de libros con más préstamos
func LibrosMasPrestados(db *sql.DB) ([]ReporteLibro, error) {
	query := `SELECT l.id, l.titulo, l.autor, COUNT(p.id) AS total
			  FROM prestamos p
			  JOIN libros l ON p.libro_id = l.id
			  GROUP BY l.id, l.titulo, l.autor
			  ORDER BY total DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reporte de libros: %w", err)
	}
	defer rows.Close()

	var reporte []ReporteLibro
	for rows.Next() {
		var r ReporteLibro
		err := rows.Scan(&r.ID, &r.Titulo, &r.Autor, &r.Total)
		if err != nil {
			return nil, fmt.Errorf("error al leer reporte: %w", err)
		}
		reporte = append(reporte, r)
	}
	return reporte, nil
}

// usuariosActivos retorna los usuarios con más actividad
func UsuariosActivos(db *sql.DB) ([]ReporteUsuario, error) {
	query := `SELECT u.id, u.nombre, u.email, COUNT(p.id) AS total
			  FROM prestamos p
			  JOIN usuarios u ON p.usuario_id = u.id
			  GROUP BY u.id, u.nombre, u.email
			  ORDER BY total DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reporte de usuarios: %w", err)
	}
	defer rows.Close()

	var reporte []ReporteUsuario
	for rows.Next() {
		var r ReporteUsuario
		err := rows.Scan(&r.ID, &r.Nombre, &r.Email, &r.Total)
		if err != nil {
			return nil, fmt.Errorf("error al leer reporte: %w", err)
		}
		reporte = append(reporte, r)
	}
	return reporte, nil
}

// librosDisponibles retorna el listado actual de libros sin prestar
func LibrosDisponibles(db *sql.DB) ([]ReporteLibro, error) {
	query := `SELECT id, titulo, autor, 0 AS total
			  FROM libros
			  WHERE disponible = TRUE`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener libros disponibles: %w", err)
	}
	defer rows.Close()

	var reporte []ReporteLibro
	for rows.Next() {
		var r ReporteLibro
		err := rows.Scan(&r.ID, &r.Titulo, &r.Autor, &r.Total)
		if err != nil {
			return nil, fmt.Errorf("error al leer reporte: %w", err)
		}
		reporte = append(reporte, r)
	}
	return reporte, nil
}
