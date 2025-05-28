package api

import (
	"encoding/json"
	"net/http"
	"sync"

	server "github.com/SeaBassLab/hyperx-server"
)

var (
	count = 0
	mu    sync.Mutex
)

func counterHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	case http.MethodPost:
		count++
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func init() {
	server.Register("/api/counter", counterHandler)
}
