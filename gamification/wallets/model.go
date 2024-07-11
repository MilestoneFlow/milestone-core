package wallets

type UserWallet struct {
	ID             string `json:"id" db:"id"`
	UserID         string `json:"user_id" db:"user_id"`
	CurrentBalance int    `json:"current_balance" db:"current_balance"`
}

type WalletTransaction struct {
	ID              string                 `json:"id" db:"id"`
	WalletID        string                 `json:"wallet_id" db:"wallet_id"`
	Amount          int                    `json:"amount" db:"amount"`
	TransactionType WalletTransactionType  `json:"transaction_type" db:"transaction_type"`
	TransactionDate string                 `json:"transaction_date" db:"transaction_date"`
	TrackData       map[string]interface{} `json:"track_data" db:"track_data"`
}

type WalletTransactionType string

const (
	DepositTransactionType  WalletTransactionType = "deposit"
	WithdrawTransactionType WalletTransactionType = "withdraw"
)
