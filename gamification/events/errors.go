package events

type EventError struct {
	HttpCode int
	Message  string
}

func (e *EventError) Error() string {
	return e.Message
}

func NewEventError(httpCode int, message string) *EventError {
	return &EventError{
		HttpCode: httpCode,
		Message:  message,
	}
}

var Errors = struct {
	KeyExistsError          *EventError
	InvalidEventError       *EventError
	EventNotFoundByKeyError *EventError
}{
	KeyExistsError:          NewEventError(400, "Event with this key already exists"),
	InvalidEventError:       NewEventError(400, "Invalid event"),
	EventNotFoundByKeyError: NewEventError(404, "Event not found by key"),
}
