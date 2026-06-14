package usuarios

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Usuario representa la estructura de un usuario en el sistema
type Usuario struct {
	ID            int
	Nombre        string
	Email         string
	FechaRegistro string
	Activo        bool
}

// registrarUsuario crea un nuevo usuario en el sistema
func RegistrarUsuario(db *sql.DB, nombre, email string) error {
	if nombre == "" || email == "" {
		return errors.New("el nombre y el email no pueden estar vacíos")
	}

	query := "INSERT INTO usuarios (nombre, email, activo) VALUES (?, ?, TRUE)"
	_, err := db.Exec(query, nombre, email)
	if err != nil {
		return fmt.Errorf("error al registrar usuario: %w", err)
	}

	fmt.Println("Usuario registrado exitosamente")
	return nil
}

// listarUsuarios retorna todos los usuarios registrados
func ListarUsuarios(db *sql.DB) ([]Usuario, error) {
	query := "SELECT id, nombre, email, fecha_registro, activo FROM usuarios"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &u.FechaRegistro, &u.Activo)
		if err != nil {
			return nil, fmt.Errorf("error al leer usuario: %w", err)
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// buscarUsuario busca por nombre o correo
func BuscarUsuario(db *sql.DB, criterio string) ([]Usuario, error) {
	if criterio == "" {
		return nil, errors.New("el criterio de búsqueda no puede estar vacío")
	}

	criterio = "%" + strings.ToLower(criterio) + "%"
	query := `SELECT id, nombre, email, fecha_registro, activo 
			  FROM usuarios 
			  WHERE LOWER(nombre) LIKE ? OR LOWER(email) LIKE ?`

	rows, err := db.Query(query, criterio, criterio)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		err := rows.Scan(&u.ID, &u.Nombre, &u.Email, &u.FechaRegistro, &u.Activo)
		if err != nil {
			return nil, fmt.Errorf("error al leer usuario: %w", err)
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// actualizarUsuario modifica los datos de un usuario
func ActualizarUsuario(db *sql.DB, id int, nombre, email string) error {
	if id <= 0 {
		return errors.New("el ID del usuario no es válido")
	}

	query := "UPDATE usuarios SET nombre = ?, email = ? WHERE id = ?"
	_, err := db.Exec(query, nombre, email, id)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	fmt.Println("Usuario actualizado exitosamente")
	return nil
}

// eliminarUsuario desactiva o elimina un usuario del sistema
func EliminarUsuario(db *sql.DB, id int) error {
	if id <= 0 {
		return errors.New("el ID del usuario no es válido")
	}

	query := "UPDATE usuarios SET activo = FALSE WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	fmt.Println("Usuario desactivado exitosamente")
	return nil
}

// Encapsulación

// SetNombre establece el nombre del usuario validando que no esté vacío
func (u *Usuario) SetNombre(nombre string) error {
	if nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	u.Nombre = nombre
	return nil
}

// SetEmail establece el email del usuario validando que no esté vacío
func (u *Usuario) SetEmail(email string) error {
	if email == "" {
		return errors.New("el email no puede estar vacío")
	}
	u.Email = email
	return nil
}

// SetActivo establece el estado activo del usuario
func (u *Usuario) SetActivo(activo bool) {
	u.Activo = activo
}
