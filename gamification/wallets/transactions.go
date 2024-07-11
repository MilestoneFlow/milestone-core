package wallets

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type TransactionsService struct {
	DbConnection *sqlx.DB
}

func NewTransactionsService(dbConnection *sqlx.DB) *TransactionsService {
	return &TransactionsService{
		DbConnection: dbConnection,
	}
}

func (s *TransactionsService) GetTransactions(workspaceId string, userId string) ([]WalletTransaction, error) {
	query := `SELECT wt.id as id, wt.user_wallet_id as wallet_id, wt.amount as amount, wt.transaction_type as transaction_type, wt.created_at as transaction_date, wt.track_data as	track_data
	FROM game_engine.user_wallet_transaction wt
	JOIN game_engine.user_wallet uw ON wt.user_wallet_id = uw.id
	WHERE uw.user_id = $1 AND uw.workspace_id = $2`

	var transactions []WalletTransaction
	err := s.DbConnection.Select(&transactions, query, userId, workspaceId)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionsService) CreateDeposit(workspaceId string, userId string, amount int, trackData map[string]interface{}) (*WalletTransaction, error) {
	tx := s.DbConnection.MustBegin()

	query := `SELECT id, current_balance FROM game_engine.user_wallet WHERE user_id = $1 AND workspace_id = $2 FOR UPDATE`
	var wallet UserWallet
	err := tx.Get(&wallet, query, userId, workspaceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	wallet.CurrentBalance += amount
	_, err = tx.Exec(`UPDATE game_engine.user_wallet SET current_balance = $1 WHERE id = $2`, wallet.CurrentBalance, wallet.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := WalletTransaction{
		WalletID:        wallet.ID,
		Amount:          amount,
		TransactionType: DepositTransactionType,
		TrackData:       trackData,
	}

	_, err = tx.NamedExec(`INSERT INTO game_engine.user_wallet_transaction (user_wallet_id, amount, transaction_type, track_data) VALUES (:wallet_id, :amount, :transaction_type, :track_data)`, transaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionsService) CreateWithdraw(workspaceId string, userId string, amount int, trackData map[string]interface{}) (*WalletTransaction, error) {
	tx := s.DbConnection.MustBegin()

	query := `SELECT id, current_balance FROM game_engine.user_wallet WHERE user_id = $1 AND workspace_id = $2 FOR UPDATE`
	var wallet UserWallet
	err := tx.Get(&wallet, query, userId, workspaceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if wallet.CurrentBalance < amount {
		tx.Rollback()
		return nil, errors.New("insufficient funds")
	}

	wallet.CurrentBalance -= amount
	_, err = tx.Exec(`UPDATE game_engine.user_wallet SET current_balance = $1 WHERE id = $2`, wallet.CurrentBalance, wallet.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := WalletTransaction{
		WalletID:        wallet.ID,
		Amount:          amount,
		TransactionType: WithdrawTransactionType,
		TrackData:       trackData,
	}

	_, err = tx.NamedExec(`INSERT INTO game_engine.user_wallet_transaction (user_wallet_id, amount, transaction_type, track_data) VALUES (:wallet_id, :amount, :transaction_type, :track_data)`, transaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &transaction, nil
}
