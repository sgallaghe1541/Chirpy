package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf(`
		<html>
		<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`, cfg.fileserverHits)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hits))
}

func (cfg *apiConfig) handlerMetricReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type respMsg struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpToPost := chirp{}
	err := decoder.Decode(&chirpToPost)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		w.WriteHeader(500)
		return
	}

	if checkMsgLength(chirpToPost.Body) {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respBody := respMsg{
		CleanedBody: cleanString(chirpToPost.Body),
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
