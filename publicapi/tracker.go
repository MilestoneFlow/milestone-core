package publicapi

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tracker struct {
	Collection *mongo.Collection
}

func (t Tracker) TrackEvents(workspaceId string, externalUserId string, events []EventTrack) error {
	rowsToInsert := make([]interface{}, len(events))
	for i, event := range events {
		event.WorkspaceID = workspaceId
		event.ExternalUserID = externalUserId
		rowsToInsert[i] = event
	}

	_, err := t.Collection.InsertMany(context.Background(), rowsToInsert)
	if err != nil {
		return err
	}

	return nil
}
