package likes

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("DisLike", disLike)
}

func disLike(ctx context.Context, e event.Event) error {
	fmt.Println("Event: ", e)
	return nil
}
