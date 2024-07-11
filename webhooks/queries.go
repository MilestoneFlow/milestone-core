package webhooks

import "github.com/jmoiron/sqlx"

func FetchPendingWebhookMessages(dbConnection *sqlx.DB, bufferSize int) ([]WebhookMessage, error) {
	var messages []WebhookMessage
	query := `
		SELECT * FROM webhooks.webhook_messages
		WHERE status = 'pending'
		ORDER BY created_at
		LIMIT $1
		`
	err := dbConnection.Select(&messages, query, bufferSize)
	return messages, err
}

func UpdateWebhookMessageStatus(dbConnection *sqlx.DB, id, status string) error {
	query := `
		UPDATE webhooks.webhook_messages
		SET status = $2
		WHERE id = $1
	`
	_, err := dbConnection.Exec(query, id, status)
	return err
}

func InsertWebhookDeliveryLog(dbConnection *sqlx.DB, log WebhookDeliveryLog) error {
	query := `
		INSERT INTO webhooks.webhook_delivery_logs (id, message_id, delivery_attempt, delivered_at, response_status, response_body, success)
		VALUES (:id, :message_id, :delivery_attempt, :delivered_at, :response_status, :response_body, :success)
	`
	_, err := dbConnection.NamedExec(query, log)
	return err
}
