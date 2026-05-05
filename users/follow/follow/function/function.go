package follow

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("Follow", follow)
}

func follow(ctx context.Context, e event.Event) error {
	fmt.Printf("received follow event: %+v", e)
	return nil
}
