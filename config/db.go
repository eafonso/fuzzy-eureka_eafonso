package config

import (
	"fmt"
	"fuzzy-eureka_eafonso/internal/models"
	"log"
	"strings"
	"time"
)

var InsertQueue chan models.Request

const batchSize = 100

func InitWorkerPool(workerCount int) {
	InsertQueue = make(chan models.Request, 1000)

	for i := 0; i < workerCount; i++ {
		go worker()
	}
}

func worker() {
	var batch []models.Request

	for {
		select {
		case req := <-InsertQueue:
			batch = append(batch, req)

			if len(batch) >= batchSize {
				insertBatch(batch)
				batch = nil
			}
		case <-time.After(2 * time.Second):
			if len(batch) > 0 {
				insertBatch(batch)
				batch = nil
			}
		}
	}
}

func insertBatch(batch []models.Request) {
	if DB == nil {
		log.Println("Banco de dados ainda não está pronto. Aguardando...")
		InitDB()
		time.Sleep(2 * time.Second)
		return
	}

	if err := DB.Ping(); err != nil {
		log.Println("Conexão com banco perdida. Tentando reconectar...")
		InitDB()
		time.Sleep(1 * time.Second)
		return
	}

	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*3)

	for i, req := range batch {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, req.ID, req.GitHubUsername, req.CommitHash)
	}

	query := fmt.Sprintf("INSERT INTO requests (id, github_username, commit_hash) VALUES %s",
		strings.Join(valueStrings, ","))

	_, err := DB.Exec(query, valueArgs...)
	if err != nil {
		log.Println("Erro ao inserir no banco:", err)
	}
}
