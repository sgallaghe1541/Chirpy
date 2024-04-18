package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnErrs struct {
		Err string `json:"error"`
	}

	respBody := returnErrs{
		Err: msg,
	}

	respondWithJSON(w, code, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func cleanString(msg string) string {
	badWords := [...]string{"kerfuffle", "sharbert", "fornax"}

	msgWords := strings.Split(msg, " ")

	for i := 0; i < len(msgWords); i++ {
		for j := 0; j < len(badWords); j++ {
			if strings.ToLower(msgWords[i]) == badWords[j] {
				msgWords[i] = "****"
			}
		}
	}

	cleanMsg := strings.Join(msgWords, " ")
	return cleanMsg
}

func checkMsgLength(msg string) bool {
	return len(msg) > 140
}
