package likes

import (
	"context"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("DisLike", disLike)
}

func disLike(ctx context.Context, e event.Event) error {
	log.Printf("Business DisLike function triggered with event")
	return nil
}
