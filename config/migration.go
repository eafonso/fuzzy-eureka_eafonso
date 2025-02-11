package config

import (
	"log"
)

func RunMigrations() {
	if DB == nil {
		log.Fatal("Banco de dados não está conectado")
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Fatal("Erro ao iniciar transação de migration:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS requests (
		id UUID PRIMARY KEY,
		github_username TEXT NOT NULL,
		commit_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err = tx.Exec(createTableQuery)
	if err != nil {
		tx.Rollback()
		log.Fatal("Erro ao criar a tabela:", err)
	}

	createIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_github_username ON requests (github_username);",
		"CREATE INDEX IF NOT EXISTS idx_commit_hash ON requests (commit_hash);",
	}

	for _, query := range createIndexes {
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			log.Fatal("Erro ao criar índice:", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Erro ao confirmar transação:", err)
	}

	log.Println("Migration executada com sucesso! (Tabela + Índices)")
}
