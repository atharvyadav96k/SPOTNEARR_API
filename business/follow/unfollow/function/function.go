package unfollow

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("Unfollow", unfollow)
}

func unfollow(ctx context.Context, e event.Event) error {
	fmt.Printf("received unfollow event: %+v", e)
	return nil
}
