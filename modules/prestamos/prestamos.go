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
	TituloLibro     string
	NombreUsuario   string
}

// Métodos Setter (Encapsulación)

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

// realizarPrestamo asigna un libro a un usuario.
// Antes de registrar el préstamo verifica que el libro esté disponible.
// Si el libro no está disponible retorna un error sin modificar la base de datos.
// Si el libro está disponible registra el préstamo y actualiza la disponibilidad del libro.
func RealizarPrestamo(db *sql.DB, libroID, usuarioID int) error {
	if libroID <= 0 || usuarioID <= 0 {
		return errors.New("el ID del libro y del usuario deben ser válidos")
	}

	// Verificar disponibilidad antes de registrar el préstamo
	disponible, err := VerificarDisponibilidad(db, libroID)
	if err != nil {
		return err
	}
	if !disponible {
		return errors.New("el libro no está disponible para préstamo")
	}

	// Registrar el préstamo en la base de datos
	query := "INSERT INTO prestamos (libro_id, usuario_id, activo) VALUES (?, ?, TRUE)"
	_, err = db.Exec(query, libroID, usuarioID)
	if err != nil {
		return fmt.Errorf("error al realizar préstamo: %w", err)
	}

	// Marcar el libro como no disponible para evitar préstamos duplicados
	_, err = db.Exec("UPDATE libros SET disponible = FALSE WHERE id = ?", libroID)
	if err != nil {
		return fmt.Errorf("error al actualizar disponibilidad: %w", err)
	}

	fmt.Println("Préstamo realizado exitosamente")
	return nil
}

// devolverLibro marca el libro como disponible nuevamente.
// Primero obtiene el libro_id del préstamo, luego cierra el préstamo
// registrando la fecha de devolución y finalmente libera el libro.
func DevolverLibro(db *sql.DB, prestamoID int) error {
	if prestamoID <= 0 {
		return errors.New("el ID del préstamo no es válido")
	}

	// Obtener el libro_id asociado al préstamo
	var libroID int
	err := db.QueryRow("SELECT libro_id FROM prestamos WHERE id = ?", prestamoID).Scan(&libroID)
	if err != nil {
		return fmt.Errorf("error al obtener préstamo: %w", err)
	}

	// Cerrar el préstamo registrando la fecha de devolución
	query := "UPDATE prestamos SET activo = FALSE, fecha_devolucion = NOW() WHERE id = ?"
	_, err = db.Exec(query, prestamoID)
	if err != nil {
		return fmt.Errorf("error al devolver libro: %w", err)
	}

	// Marcar el libro como disponible nuevamente
	_, err = db.Exec("UPDATE libros SET disponible = TRUE WHERE id = ?", libroID)
	if err != nil {
		return fmt.Errorf("error al actualizar disponibilidad: %w", err)
	}

	fmt.Println("Libro devuelto exitosamente")
	return nil
}

// listarPrestamos muestra todos los préstamos activos del sistema.
// Usa JOIN para obtener el título del libro y el nombre del usuario
// en lugar de mostrar solo los IDs.
func ListarPrestamos(db *sql.DB) ([]Prestamo, error) {
	query := `SELECT p.id, p.libro_id, p.usuario_id, p.fecha_prestamo, p.activo,
			  l.titulo, u.nombre
			  FROM prestamos p
			  JOIN libros l ON p.libro_id = l.id
			  JOIN usuarios u ON p.usuario_id = u.id
			  WHERE p.activo = TRUE`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al listar préstamos: %w", err)
	}
	defer rows.Close()

	var prestamos []Prestamo
	for rows.Next() {
		var p Prestamo
		err := rows.Scan(&p.ID, &p.LibroID, &p.UsuarioID, &p.FechaPrestamo, &p.Activo, &p.TituloLibro, &p.NombreUsuario)
		if err != nil {
			return nil, fmt.Errorf("error al leer préstamo: %w", err)
		}
		prestamos = append(prestamos, p)
	}
	return prestamos, nil
}

// historialPrestamos muestra los préstamos pasados de un usuario específico.
// Solo devuelve préstamos con activo = FALSE, es decir, ya devueltos.
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

// verificarDisponibilidad consulta si un libro está libre para préstamo.
// Retorna true si el libro está disponible, false si ya está prestado.
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