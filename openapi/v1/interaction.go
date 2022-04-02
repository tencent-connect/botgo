package v1

import (
	"context"
)

// PutInteraction 更新 interaction
func (o *openAPI) PutInteraction(ctx context.Context,
	interactionID string, body string) error {
	_, err := o.request(ctx).
		SetPathParam("interaction_id", interactionID).
		SetBody(body).
		Put(o.getURL(interactionsURI))
	return err
}
