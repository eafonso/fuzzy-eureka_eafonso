package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"fuzzy-eureka_eafonso/config"
	"fuzzy-eureka_eafonso/internal/models"

	"github.com/google/uuid"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}

	select {
	case config.InsertQueue <- req:
	default:
		log.Println("Fila de inserções cheia, descartando requisição")
		http.Error(w, "Servidor sobrecarregado, tente novamente", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
