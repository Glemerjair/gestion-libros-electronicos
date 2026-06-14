package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gestion-libros-electronicos/database"
	"gestion-libros-electronicos/modules/categorias"
	"gestion-libros-electronicos/modules/libros"
	"gestion-libros-electronicos/modules/prestamos"
	"gestion-libros-electronicos/modules/reportes"
	"gestion-libros-electronicos/modules/usuarios"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	db = database.Conectar()
	defer db.Close()

	fmt.Println("=== Sistema de Gestión de Libros Electrónicos ===")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// ── Ruta principal ──
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/libros")
	})

	// ── Rutas Libros ──
	r.GET("/libros", func(c *gin.Context) {
		lista, err := libros.ListarLibros(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "libros.html", gin.H{
			"titulo": "Gestión de Libros",
			"libros": lista,
		})
	})

	r.POST("/libros/agregar", func(c *gin.Context) {
		titulo := c.PostForm("titulo")
		autor := c.PostForm("autor")
		anio, _ := strconv.Atoi(c.PostForm("anio"))
		enlaceDrive := c.PostForm("enlace_drive")
		categoriaID, _ := strconv.Atoi(c.PostForm("categoria_id"))

		err := libros.AgregarLibro(db, titulo, autor, anio, enlaceDrive, categoriaID)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/libros")
	})

	r.GET("/libros/eliminar/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := libros.EliminarLibro(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/libros")
	})

	r.GET("/libros/buscar", func(c *gin.Context) {
		criterio := c.Query("criterio")
		lista, err := libros.ListarLibros(db)
		if err != nil {
			log.Println("Error:", err)
		}
		var resultados []libros.Libro
		if criterio != "" {
			resultados, err = libros.BuscarLibro(db, criterio)
			if err != nil {
				log.Println("Error:", err)
			}
		}
		c.HTML(http.StatusOK, "libros.html", gin.H{
			"titulo":     "Gestión de Libros",
			"libros":     lista,
			"resultados": resultados,
		})
	})

	// ── Rutas Usuarios ──
	r.GET("/usuarios", func(c *gin.Context) {
		lista, err := usuarios.ListarUsuarios(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "usuarios.html", gin.H{
			"titulo":   "Gestión de Usuarios",
			"usuarios": lista,
		})
	})

	r.POST("/usuarios/registrar", func(c *gin.Context) {
		nombre := c.PostForm("nombre")
		email := c.PostForm("email")

		err := usuarios.RegistrarUsuario(db, nombre, email)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/eliminar/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := usuarios.EliminarUsuario(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/buscar", func(c *gin.Context) {
		criterio := c.Query("criterio")
		lista, err := usuarios.ListarUsuarios(db)
		if err != nil {
			log.Println("Error:", err)
		}
		var resultados []usuarios.Usuario
		if criterio != "" {
			resultados, err = usuarios.BuscarUsuario(db, criterio)
			if err != nil {
				log.Println("Error:", err)
			}
		}
		c.HTML(http.StatusOK, "usuarios.html", gin.H{
			"titulo":     "Gestión de Usuarios",
			"usuarios":   lista,
			"resultados": resultados,
		})
	})

	// ── Rutas Categorías ──
	r.GET("/categorias", func(c *gin.Context) {
		lista, err := categorias.ListarCategorias(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "categorias.html", gin.H{
			"titulo":     "Gestión de Categorías",
			"categorias": lista,
		})
	})

	r.POST("/categorias/agregar", func(c *gin.Context) {
		nombre := c.PostForm("nombre")
		descripcion := c.PostForm("descripcion")

		err := categorias.AgregarCategoria(db, nombre, descripcion)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/categorias")
	})

	r.GET("/categorias/eliminar/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := categorias.EliminarCategoria(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/categorias")
	})

	// ── Rutas Préstamos ──
	r.GET("/prestamos", func(c *gin.Context) {
		lista, err := prestamos.ListarPrestamos(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "prestamos.html", gin.H{
			"titulo":    "Gestión de Préstamos",
			"prestamos": lista,
		})
	})

	r.POST("/prestamos/realizar", func(c *gin.Context) {
		libroID, _ := strconv.Atoi(c.PostForm("libro_id"))
		usuarioID, _ := strconv.Atoi(c.PostForm("usuario_id"))

		err := prestamos.RealizarPrestamo(db, libroID, usuarioID)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/prestamos")
	})

	r.POST("/prestamos/devolver", func(c *gin.Context) {
		prestamoID, _ := strconv.Atoi(c.PostForm("prestamo_id"))

		err := prestamos.DevolverLibro(db, prestamoID)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/prestamos")
	})

	// ── Rutas Reportes ──
	r.GET("/reportes", func(c *gin.Context) {
		masPrestados, err := reportes.LibrosMasPrestados(db)
		if err != nil {
			log.Println("Error:", err)
		}
		activos, err := reportes.UsuariosActivos(db)
		if err != nil {
			log.Println("Error:", err)
		}
		disponibles, err := reportes.LibrosDisponibles(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "reportes.html", gin.H{
			"titulo":       "Reportes",
			"masPrestados": masPrestados,
			"activos":      activos,
			"disponibles":  disponibles,
		})
	})

	fmt.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
