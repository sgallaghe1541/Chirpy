package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handlerFileServer(filepathRoot)))

	mux.HandleFunc("GET /healthz/", handlerReadiness)
	mux.HandleFunc("GET /metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerMetricReset)

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
