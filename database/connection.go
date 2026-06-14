package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Conectar establece la conexión con la base de datos MySQL
func Conectar() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo verificar la conexión:", err)
	}

	fmt.Println("Conexión exitosa a la base de datos")
	return db
}
