package main

import (
	"log"
	"net/http"

	"github.com/sgallaghe1541/chirpy/internal/database"
)

type apiConfig struct {
	db             *database.DB
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	dbChirps, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		db:             dbChirps,
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz/", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handlerMetricReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
