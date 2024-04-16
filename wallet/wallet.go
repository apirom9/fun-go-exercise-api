package wallet

import "time"

type Wallet struct {
	ID         int       `json:"id" example:"1"`
	UserID     int       `json:"user_id" example:"1"`
	UserName   string    `json:"user_name" example:"John Doe"`
	WalletName string    `json:"wallet_name" example:"John's Wallet"`
	WalletType string    `json:"wallet_type" example:"Create Card"`
	Balance    float64   `json:"balance" example:"100.00"`
	CreatedAt  time.Time `json:"created_at" example:"2024-03-25T14:19:00.729237Z"`
}

type CreateWallet struct {
	UserID     int     `json:"user_id"`
	UserName   string  `json:"user_name"`
	WalletName string  `json:"wallet_name"`
	WalletType string  `json:"wallet_type"`
	Balance    float64 `json:"balance"`
}

type UpdateWallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	UserName   string  `json:"user_name"`
	WalletName string  `json:"wallet_name"`
	WalletType string  `json:"wallet_type"`
	Balance    float64 `json:"balance"`
}
