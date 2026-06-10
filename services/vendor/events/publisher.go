package events

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

// Notifier triggers the Search Service to process pending outbox rows via HTTP
// POST instead of Redis pub/sub. Call as go notifier.Notify(context.Background()).
type Notifier struct {
	url    string
	client *http.Client
}

func NewNotifier(searchServiceURL string) *Notifier {
	return &Notifier{
		url:    searchServiceURL + "/internal/sync",
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (n *Notifier) Notify(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, bytes.NewReader(nil))
	if err != nil {
		log.Printf("notifier: build request: %v", err)
		return
	}
	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("notifier: POST %s: %v", n.url, err)
		return
	}
	resp.Body.Close()
}
