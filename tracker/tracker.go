package tracker

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

func (t Tracker) FetchTrackDataForFlow(flowID string) ([]EventTrack, error) {
	cursor, err := t.Collection.Find(context.Background(), map[string]string{"entityId": flowID})
	if err != nil {
		return nil, err
	}

	var events []EventTrack
	if err = cursor.All(context.Background(), &events); err != nil {
		return nil, err
	}

	return events, nil
}
