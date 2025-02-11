package config

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var once sync.Once

func InitDB() {
	once.Do(func() {
		if err_ := godotenv.Load(); err_ != nil {
			log.Println("No .env file found")
		}

		connStr := os.Getenv("DATABASE_URL")
		var err error

		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Erro ao conectar ao banco:", err)
		}

		DB.SetMaxOpenConns(10)
		DB.SetMaxIdleConns(5)

		if err = DB.Ping(); err != nil {
			log.Fatal("Erro ao conectar ao banco (Ping):", err)
		}

		log.Println("Banco de dados conectado com sucesso!")

		RunMigrations()
	})
}
