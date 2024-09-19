package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/PaytonWebber/go-webhook-listener/config"
)

type GitHubPushEvent struct {
	Ref string `json:"ref"`
}

type WebhookHandler struct {
	cfg *config.Config
}

func NewWebhookHandler(cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{cfg: cfg}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.cfg.Repository.Path
	branch := h.cfg.Repository.Branch
	restart := h.cfg.RestartCommand

	// Read the signature from the header
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusUnauthorized)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Reset r.Body so it can be read again
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	// Verify the signature
	if !h.verifySignature(signature, body) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Now decode the JSON payload
	var event GitHubPushEvent
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if event.Ref == "refs/heads/"+branch {
		log.Printf("Received push to %s", event.Ref)

		// Pull from the path and branch
		log.Printf("Pulling from %s", path)
		cmd := exec.Command("git", "pull", "origin", branch)
		cmd.Dir = path
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error pulling code: %v", err)
			http.Error(w, "Error pulling code", http.StatusInternalServerError)
			return
		}
		log.Printf("Pull output: %s", output)

		// Execute the restart command
		log.Printf("Executing %s", restart)
		restartCmd := exec.Command("sh", "-c", restart)
		output, err = restartCmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing restart command: %v", err)
			http.Error(w, "Error executing restart command", http.StatusInternalServerError)
			return
		}
		log.Printf("Restart output: %s", output)
	} else {
		log.Printf("Ignoring push to %s", event.Ref)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) verifySignature(signature string, body []byte) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(signature, prefix) {
		return false
	}
	signature = strings.TrimPrefix(signature, prefix)

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	// Compute HMAC SHA256
	mac := hmac.New(sha256.New, []byte(h.cfg.Secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(signatureBytes, expectedMAC)
}
