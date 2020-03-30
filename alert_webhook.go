package alertfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type StackDriveIncident struct {
	IncidentID    string `json:"incident_id"`
	ResourceID    string `json:"resource_id"`
	ResourceName  string `json:"resource_name"`
	State         string `json:"state"`
	StartedAt     Time   `json:"started_at"`
	EndedAt       Time   `json:"ended_at"`
	PolicyName    string `json:"policy_name"`
	ConditionName string `json:"condition_name"`
	URL           string `json:"url"`
	Summary       string `json:"summary"`
}

type StackDriveAlert struct {
	Incident StackDriveIncident `json:"incident"`
	Version  string             `json:"version"`
}

type WebhookRequest struct {
	Text string `json:"text"`
}

func AlertFunc(w http.ResponseWriter, r *http.Request) {
	var alert StackDriveAlert
	token := r.URL.Query().Get("auth_token")
	if token != os.Getenv("AUTH_TOKEN") {
		log.Printf("Unauthorized Access with Token: %s\n", token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webhookReq := WebhookRequest{
		Text: fmt.Sprintf("##### %s\n%s\n%s\n_%s_", alert.Incident.PolicyName,
			alert.Incident.Summary, alert.Incident.URL,
			alert.Incident.StartedAt.Time().Format(time.RFC1123Z)),
	}

	reqBody, err := json.Marshal(webhookReq)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(([]byte("OK")))
}
