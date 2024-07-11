package webhooks

type WebhookEndpoint struct {
	ID          string `json:"id" db:"id"`
	WorkspaceID string `json:"workspace_id" db:"workspace_id"`
	Url         string `json:"url" db:"url"`
}

type WebhookMessage struct {
	ID          string                 `json:"id" db:"id"`
	WorkspaceID string                 `json:"workspace_id" db:"workspace_id"`
	EndpointID  string                 `json:"endpoint_id" db:"endpoint_id"`
	Payload     map[string]interface{} `json:"payload" db:"payload"`
	SentAt      string                 `json:"sent_at" db:"sent_at"`
	Status      WebhookMessageStatus   `json:"status" db:"status"`
}
type WebhookMessageStatus string

const (
	WebhookMessageStatusSent    WebhookMessageStatus = "sent"
	WebhookMessageStatusFailed  WebhookMessageStatus = "failed"
	WebhookMessageStatusPending WebhookMessageStatus = "pending"
)

type WebhookDeliveryLog struct {
	ID              string `json:"id" db:"id"`
	MessageID       string `json:"message_id" db:"message_id"`
	DeliveryAttempt int    `json:"delivery_attempt" db:"delivery_attempt"`
	DeliveredAt     string `json:"delivered_at" db:"delivered_at"`
	ResponseStatus  int    `json:"response_status" db:"response_status"`
	ResponseBody    string `json:"response_body" db:"response_body"`
	Success         bool   `json:"success" db:"success"`
}
