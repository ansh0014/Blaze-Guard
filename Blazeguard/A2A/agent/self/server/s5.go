package s5
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"blazeguard/shared"
)
type MessageHandler func(eventType string, payload map[string]interface{})
var messageHandler MessageHandler
func SetMessageHandler(handler MessageHandler) {
	messageHandler = handler
}
func StartHTTPServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent":"self","status":"ok"}`))
	})
	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		var msg shared.A2AMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		fmt.Printf("[Self Agent] A2A received from %s | event: %s\n", msg.From, msg.EventType)

		if messageHandler != nil {
			messageHandler(msg.EventType, msg.Payload)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"received"}`))
	})
	port := "8005"
	fmt.Printf("[Self Agent] HTTP server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
