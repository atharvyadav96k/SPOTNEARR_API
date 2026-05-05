package user_likes

import (
	"context"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("Like", like)
}

func like(ctx context.Context, e event.Event) error {
	log.Printf("User Like function triggered with event")
	return nil
}
