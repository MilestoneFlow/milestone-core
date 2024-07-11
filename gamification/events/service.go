package events

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"milestone_core/shared/sql"
)

type Service struct {
	DbConnection *sqlx.DB
}

func (s Service) CreateEvent(workspaceId string, event Event) error {
	keyExists, err := s.KeyExists(workspaceId, event.Key)
	if err != nil {
		return err
	}
	if keyExists {
		return Errors.KeyExistsError
	}

	_, err = s.DbConnection.NamedExec("INSERT INTO game_engine.event (workspace_id, key, name) VALUES (:workspaceId, :key, :name)", map[string]interface{}{
		"workspaceId": workspaceId,
		"key":         event.Key,
		"name":        event.Name,
	})
	return err
}

func (s Service) GetEventById(workspaceId string, id string) (*Event, error) {
	event, err := sql.FetchOne[Event](s.DbConnection, "SELECT id, key, name FROM game_engine.event WHERE workspace_id = $1 AND id = $2 AND deleted_at IS NULL", workspaceId, id)
	return event, err
}

func (s Service) GetEventByKey(workspaceId string, key string) (*Event, error) {
	event, err := sql.FetchOne[Event](s.DbConnection, "SELECT id, key, name FROM game_engine.event WHERE workspace_id = $1 AND key = $2 AND deleted_at IS NULL", workspaceId, key)
	return event, err
}

func (s Service) KeyExists(workspaceId string, key string) (bool, error) {
	count, err := sql.FetchOne[int](s.DbConnection, "SELECT COUNT(*) FROM game_engine.event WHERE workspace_id = $1 AND key = $2 AND deleted_at IS NULL LIMIT 1", workspaceId, key)
	return count != nil && *count != 0, err
}

func (s Service) GetEvents(workspaceId string) ([]Event, error) {
	events, err := sql.FetchMultiple[Event](s.DbConnection, "SELECT id, key, name FROM game_engine.event WHERE workspace_id = $1 AND deleted_at IS NULL", workspaceId)
	return events, err
}

func (s Service) UpdateEvent(workspaceId string, id string, event *Event) error {
	if event.Key == "" || event.Name == "" {
		return Errors.InvalidEventError
	}
	currentEvent, err := s.GetEventById(workspaceId, id)
	if currentEvent == nil {
		return Errors.InvalidEventError
	}
	if err != nil {
		return err
	}

	if event.Key != currentEvent.Key {
		keyExists, err := s.KeyExists(workspaceId, event.Key)
		if err != nil {
			return err
		}
		if keyExists {
			return Errors.KeyExistsError
		}
	}

	_, err = s.DbConnection.NamedExec("UPDATE game_engine.event SET key = :key, name = :name WHERE workspace_id = :workspaceId AND id = :id AND deleted_at IS NULL", map[string]interface{}{
		"workspaceId": workspaceId,
		"id":          id,
		"key":         event.Key,
		"name":        event.Name,
	})
	return err
}

func (s Service) DeleteEvent(workspaceId string, id string) error {
	item, err := s.GetEventById(workspaceId, id)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}

	tx := s.DbConnection.MustBegin()
	_, err = tx.Exec("SELECT game_engine.delete_event_record($1)", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s Service) CreateEventUser(workspaceId string, eventKey string, userID string, metadata *json.RawMessage) error {
	event, err := s.GetEventByKey(workspaceId, eventKey)
	if event == nil {
		return Errors.EventNotFoundByKeyError
	}
	if err != nil {
		return err
	}

	_, err = s.DbConnection.Exec("INSERT INTO game_engine.user_events(workspace_id, user_id, event_id, metadata) VALUES ($1, $2, $3, $4)", workspaceId, userID, event.ID, metadata)
	return err
}
