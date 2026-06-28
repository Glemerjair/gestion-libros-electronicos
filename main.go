package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gestion-libros-electronicos/database"
	"gestion-libros-electronicos/modules/categorias"
	"gestion-libros-electronicos/modules/libros"
	"gestion-libros-electronicos/modules/prestamos"
	"gestion-libros-electronicos/modules/reportes"
	"gestion-libros-electronicos/modules/usuarios"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db *sql.DB

// Middleware para verificar sesión
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("admin")
		if user == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	db = database.Conectar()
	defer db.Close()

	godotenv.Load()
	adminUser := os.Getenv("ADMIN_USER")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	fmt.Println("=== Sistema de Gestión de Libros Electrónicos ===")

	r := gin.Default()

	// Configurar sesiones
	store := cookie.NewStore([]byte("secreto-gestion-libros"))
	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// ── Rutas de Login ──
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		usuario := c.PostForm("usuario")
		password := c.PostForm("password")

		if usuario == adminUser && password == adminPassword {
			session := sessions.Default(c)
			session.Set("admin", usuario)
			session.Save()
			c.Redirect(http.StatusSeeOther, "/libros")
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"error": "Usuario o contraseña incorrectos",
			})
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
	})

	// ── Ruta principal ──
	r.GET("/", authRequired(), func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/libros")
	})

	// ── Rutas Libros ──
	r.GET("/libros", authRequired(), func(c *gin.Context) {
		lista, err := libros.ListarLibros(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "libros.html", gin.H{
			"titulo": "Gestión de Libros",
			"libros": lista,
		})
	})

	r.POST("/libros/agregar", authRequired(), func(c *gin.Context) {
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

	r.GET("/libros/eliminar/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := libros.EliminarLibro(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/libros")
	})

	r.GET("/libros/editar/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		lista, err := libros.ListarLibros(db)
		if err != nil {
			log.Println("Error:", err)
		}
		var libroEditar libros.Libro
		for _, l := range lista {
			if l.ID == id {
				libroEditar = l
				break
			}
		}
		c.HTML(http.StatusOK, "libros.html", gin.H{
			"titulo":      "Gestión de Libros",
			"libros":      lista,
			"libroEditar": libroEditar,
		})
	})

	r.POST("/libros/actualizar", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		titulo := c.PostForm("titulo")
		autor := c.PostForm("autor")
		anio, _ := strconv.Atoi(c.PostForm("anio"))
		enlaceDrive := c.PostForm("enlace_drive")

		err := libros.ActualizarLibro(db, id, titulo, autor, anio, enlaceDrive)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/libros")
	})

	r.GET("/libros/buscar", authRequired(), func(c *gin.Context) {
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

	r.GET("/libros/buscar-json", authRequired(), func(c *gin.Context) {
		criterio := c.Query("q")
		lista, err := libros.BuscarLibro(db, criterio)
		if err != nil {
			c.JSON(http.StatusOK, []libros.Libro{})
			return
		}
		c.JSON(http.StatusOK, lista)
	})

	// ── Rutas Usuarios ──
	r.GET("/usuarios", authRequired(), func(c *gin.Context) {
		lista, err := usuarios.ListarUsuarios(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "usuarios.html", gin.H{
			"titulo":   "Gestión de Usuarios",
			"usuarios": lista,
		})
	})

	r.POST("/usuarios/registrar", authRequired(), func(c *gin.Context) {
		nombre := c.PostForm("nombre")
		email := c.PostForm("email")

		err := usuarios.RegistrarUsuario(db, nombre, email)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/eliminar/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := usuarios.EliminarUsuario(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/editar/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		lista, err := usuarios.ListarUsuarios(db)
		if err != nil {
			log.Println("Error:", err)
		}
		var usuarioEditar usuarios.Usuario
		for _, u := range lista {
			if u.ID == id {
				usuarioEditar = u
				break
			}
		}
		c.HTML(http.StatusOK, "usuarios.html", gin.H{
			"titulo":        "Gestión de Usuarios",
			"usuarios":      lista,
			"usuarioEditar": usuarioEditar,
		})
	})

	r.POST("/usuarios/actualizar", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		nombre := c.PostForm("nombre")
		email := c.PostForm("email")

		err := usuarios.ActualizarUsuario(db, id, nombre, email)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/toggleactivo/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		lista, err := usuarios.ListarUsuarios(db)
		if err != nil {
			log.Println("Error:", err)
		}
		var estadoActual bool
		for _, u := range lista {
			if u.ID == id {
				estadoActual = u.Activo
				break
			}
		}
		_, err = db.Exec("UPDATE usuarios SET activo = ? WHERE id = ?", !estadoActual, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/usuarios")
	})

	r.GET("/usuarios/buscar", authRequired(), func(c *gin.Context) {
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

	r.GET("/usuarios/buscar-json", authRequired(), func(c *gin.Context) {
		criterio := c.Query("q")
		lista, err := usuarios.BuscarUsuario(db, criterio)
		if err != nil {
			c.JSON(http.StatusOK, []usuarios.Usuario{})
			return
		}
		c.JSON(http.StatusOK, lista)
	})

	// ── Rutas Categorías ──
	r.GET("/categorias", authRequired(), func(c *gin.Context) {
		lista, err := categorias.ListarCategorias(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "categorias.html", gin.H{
			"titulo":     "Gestión de Categorías",
			"categorias": lista,
		})
	})

	r.POST("/categorias/agregar", authRequired(), func(c *gin.Context) {
		nombre := c.PostForm("nombre")
		descripcion := c.PostForm("descripcion")

		err := categorias.AgregarCategoria(db, nombre, descripcion)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/categorias")
	})

	r.GET("/categorias/eliminar/:id", authRequired(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := categorias.EliminarCategoria(db, id)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/categorias")
	})

	r.GET("/categorias/buscar-json", authRequired(), func(c *gin.Context) {
		q := c.Query("q")
		lista, err := categorias.ListarCategorias(db)
		if err != nil {
			c.JSON(http.StatusOK, []categorias.Categoria{})
			return
		}
		var resultado []categorias.Categoria
		for _, cat := range lista {
			if strings.Contains(strings.ToLower(cat.Nombre), strings.ToLower(q)) {
				resultado = append(resultado, cat)
			}
		}
		c.JSON(http.StatusOK, resultado)
	})

	// ── Rutas Préstamos ──
	r.GET("/prestamos", authRequired(), func(c *gin.Context) {
		lista, err := prestamos.ListarPrestamos(db)
		if err != nil {
			log.Println("Error:", err)
		}
		c.HTML(http.StatusOK, "prestamos.html", gin.H{
			"titulo":    "Gestión de Préstamos",
			"prestamos": lista,
		})
	})

	r.POST("/prestamos/realizar", authRequired(), func(c *gin.Context) {
		libroID, _ := strconv.Atoi(c.PostForm("libro_id"))
		usuarioID, _ := strconv.Atoi(c.PostForm("usuario_id"))

		err := prestamos.RealizarPrestamo(db, libroID, usuarioID)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/prestamos")
	})

	r.POST("/prestamos/devolver", authRequired(), func(c *gin.Context) {
		prestamoID, _ := strconv.Atoi(c.PostForm("prestamo_id"))

		err := prestamos.DevolverLibro(db, prestamoID)
		if err != nil {
			log.Println("Error:", err)
		}
		c.Redirect(http.StatusSeeOther, "/prestamos")
	})

	// ── Rutas Reportes ──
	r.GET("/reportes", authRequired(), func(c *gin.Context) {
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

	// ── API Servicios Web ──

// 1. Listar todos los libros
r.GET("/api/libros", func(c *gin.Context) {
    lista, err := libros.ListarLibros(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 2. Buscar libro por criterio
r.GET("/api/libros/buscar", func(c *gin.Context) {
    criterio := c.Query("q")
    lista, err := libros.BuscarLibro(db, criterio)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 3. Listar todos los usuarios
r.GET("/api/usuarios", func(c *gin.Context) {
    lista, err := usuarios.ListarUsuarios(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 4. Listar todas las categorías
r.GET("/api/categorias", func(c *gin.Context) {
    lista, err := categorias.ListarCategorias(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 5. Listar préstamos activos
r.GET("/api/prestamos", func(c *gin.Context) {
    lista, err := prestamos.ListarPrestamos(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 6. Libros más prestados
r.GET("/api/reportes/libros-mas-prestados", func(c *gin.Context) {
    lista, err := reportes.LibrosMasPrestados(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 7. Usuarios más activos
r.GET("/api/reportes/usuarios-activos", func(c *gin.Context) {
    lista, err := reportes.UsuariosActivos(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

// 8. Libros disponibles
r.GET("/api/reportes/libros-disponibles", func(c *gin.Context) {
    lista, err := reportes.LibrosDisponibles(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, lista)
})

	fmt.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
