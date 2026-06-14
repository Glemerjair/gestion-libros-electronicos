package prestamos

import (
	"database/sql"
	"errors"
	"fmt"
)

// Prestamo representa la estructura de un préstamo en el sistema
type Prestamo struct {
	ID              int
	LibroID         int
	UsuarioID       int
	FechaPrestamo   string
	FechaDevolucion string
	Activo          bool
}

// realizarPrestamo asigna un libro a un usuario
func RealizarPrestamo(db *sql.DB, libroID, usuarioID int) error {
	if libroID <= 0 || usuarioID <= 0 {
		return errors.New("el ID del libro y del usuario deben ser válidos")
	}

	// verificar disponibilidad antes de prestar
	disponible, err := VerificarDisponibilidad(db, libroID)
	if err != nil {
		return err
	}
	if !disponible {
		return errors.New("el libro no está disponible para préstamo")
	}

	// registrar el préstamo
	query := "INSERT INTO prestamos (libro_id, usuario_id, activo) VALUES (?, ?, TRUE)"
	_, err = db.Exec(query, libroID, usuarioID)
	if err != nil {
		return fmt.Errorf("error al realizar préstamo: %w", err)
	}

	// marcar el libro como no disponible
	_, err = db.Exec("UPDATE libros SET disponible = FALSE WHERE id = ?", libroID)
	if err != nil {
		return fmt.Errorf("error al actualizar disponibilidad: %w", err)
	}

	fmt.Println("Préstamo realizado exitosamente")
	return nil
}

// devolverLibro marca el libro como disponible nuevamente
func DevolverLibro(db *sql.DB, prestamoID int) error {
	if prestamoID <= 0 {
		return errors.New("el ID del préstamo no es válido")
	}

	// obtener el libro_id del préstamo
	var libroID int
	err := db.QueryRow("SELECT libro_id FROM prestamos WHERE id = ?", prestamoID).Scan(&libroID)
	if err != nil {
		return fmt.Errorf("error al obtener préstamo: %w", err)
	}

	// cerrar el préstamo
	query := "UPDATE prestamos SET activo = FALSE, fecha_devolucion = NOW() WHERE id = ?"
	_, err = db.Exec(query, prestamoID)
	if err != nil {
		return fmt.Errorf("error al devolver libro: %w", err)
	}

	// marcar el libro como disponible
	_, err = db.Exec("UPDATE libros SET disponible = TRUE WHERE id = ?", libroID)
	if err != nil {
		return fmt.Errorf("error al actualizar disponibilidad: %w", err)
	}

	fmt.Println("Libro devuelto exitosamente")
	return nil
}

// listarPrestamos muestra todos los préstamos activos
func ListarPrestamos(db *sql.DB) ([]Prestamo, error) {
	query := `SELECT id, libro_id, usuario_id, fecha_prestamo, activo 
			  FROM prestamos WHERE activo = TRUE`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al listar préstamos: %w", err)
	}
	defer rows.Close()

	var prestamos []Prestamo
	for rows.Next() {
		var p Prestamo
		err := rows.Scan(&p.ID, &p.LibroID, &p.UsuarioID, &p.FechaPrestamo, &p.Activo)
		if err != nil {
			return nil, fmt.Errorf("error al leer préstamo: %w", err)
		}
		prestamos = append(prestamos, p)
	}
	return prestamos, nil
}

// historialPrestamos muestra préstamos pasados por usuario
func HistorialPrestamos(db *sql.DB, usuarioID int) ([]Prestamo, error) {
	if usuarioID <= 0 {
		return nil, errors.New("el ID del usuario no es válido")
	}

	query := `SELECT id, libro_id, usuario_id, fecha_prestamo, fecha_devolucion, activo 
			  FROM prestamos WHERE usuario_id = ? AND activo = FALSE`
	rows, err := db.Query(query, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener historial: %w", err)
	}
	defer rows.Close()

	var prestamos []Prestamo
	for rows.Next() {
		var p Prestamo
		err := rows.Scan(&p.ID, &p.LibroID, &p.UsuarioID, &p.FechaPrestamo, &p.FechaDevolucion, &p.Activo)
		if err != nil {
			return nil, fmt.Errorf("error al leer historial: %w", err)
		}
		prestamos = append(prestamos, p)
	}
	return prestamos, nil
}

// verificarDisponibilidad consulta si un libro está libre
func VerificarDisponibilidad(db *sql.DB, libroID int) (bool, error) {
	if libroID <= 0 {
		return false, errors.New("el ID del libro no es válido")
	}

	var disponible bool
	err := db.QueryRow("SELECT disponible FROM libros WHERE id = ?", libroID).Scan(&disponible)
	if err != nil {
		return false, fmt.Errorf("error al verificar disponibilidad: %w", err)
	}

	return disponible, nil
}

// Encapsulación

// SetLibroID establece el ID del libro validando que sea válido
func (p *Prestamo) SetLibroID(libroID int) error {
	if libroID <= 0 {
		return errors.New("el ID del libro no es válido")
	}
	p.LibroID = libroID
	return nil
}

// SetUsuarioID establece el ID del usuario validando que sea válido
func (p *Prestamo) SetUsuarioID(usuarioID int) error {
	if usuarioID <= 0 {
		return errors.New("el ID del usuario no es válido")
	}
	p.UsuarioID = usuarioID
	return nil
}

// SetActivo establece el estado activo del préstamo
func (p *Prestamo) SetActivo(activo bool) {
	p.Activo = activo
}
