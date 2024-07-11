package wallets

import "github.com/jmoiron/sqlx"

type WalletService struct {
	DbConnection *sqlx.DB
}

func NewWalletService(dbConnection *sqlx.DB) *WalletService {
	return &WalletService{
		DbConnection: dbConnection,
	}
}

func (s *WalletService) GetWallet(workspaceId string, userId string) (*UserWallet, error) {
	query := `SELECT id, user_id, current_balance FROM game_engine.user_wallet WHERE user_id = $1 AND workspace_id = $2`
	var wallet UserWallet
	err := s.DbConnection.Get(&wallet, query, userId, workspaceId)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (s *WalletService) CreateWallet(workspaceId string, userId string) (*UserWallet, error) {
	query := `INSERT INTO game_engine.user_wallet (user_id, workspace_id, current_balance) VALUES ($1, $2, 0) RETURNING id, user_id, current_balance`
	var wallet UserWallet
	err := s.DbConnection.Get(&wallet, query, userId, workspaceId)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
