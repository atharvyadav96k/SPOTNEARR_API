package follow

import (
	"context"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("Follow", follow)
}

func follow(ctx context.Context, e event.Event) error {
	log.Printf("User Follow function triggered with event")
	return nil
}
