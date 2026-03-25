package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"bytes"
	"net/http"
)
func SendToAgent(agentName string, msg A2AMessage)error{
	// fetch url of agent from the secret file
 url := os.Getenv(agentName + "_URL")
    if url == "" {
        return fmt.Errorf("agent URL not found: %s", agentName)
    }
	msg.Timestamp=time.Now().Format(time.RFC3339)
	body,_:=json.Marshal(msg)
	resp, err := http.Post(
		url+"/receive",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err !=nil{
		return fmt.Errorf("failed to send message to agent %s: %v", agentName, err)
	}
fmt.Printf("[A2A] %s → %s | event: %s | status: %d\n",
        msg.From, agentName, msg.EventType, resp.StatusCode)

    return nil
}