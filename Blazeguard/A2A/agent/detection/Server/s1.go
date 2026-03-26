package s1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"blazeguard/shared"
)

func StartHTTPServer() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"agent":"detection","status":"ok"}`))
	})

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		var msg shared.A2AMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		fmt.Printf("[Detection] A2A received from %s | event: %s\n", msg.From, msg.EventType)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"received"}`))
	})

	port := "8001"
	fmt.Printf("[Detection Agent] HTTP server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

