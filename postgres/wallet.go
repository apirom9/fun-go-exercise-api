package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/KKGo-Software-engineering/fun-exercise-api/wallet"
)

type Wallet struct {
	ID         int       `postgres:"id"`
	UserID     int       `postgres:"user_id"`
	UserName   string    `postgres:"user_name"`
	WalletName string    `postgres:"wallet_name"`
	WalletType string    `postgres:"wallet_type"`
	Balance    float64   `postgres:"balance"`
	CreatedAt  time.Time `postgres:"created_at"`
}

func (p *Postgres) Wallets() ([]wallet.Wallet, error) {
	rows, err := p.Db.Query("SELECT * FROM user_wallet")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []wallet.Wallet
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		})
	}
	return wallets, nil
}

func (p *Postgres) WalletsByType(walletType string) ([]wallet.Wallet, error) {

	wallets := []wallet.Wallet{}

	rows, err := p.Db.Query("SELECT * FROM user_wallet WHERE wallet_type = $1", walletType)
	if err != nil {
		return wallets, nil
	}
	defer rows.Close()

	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		})
	}
	return wallets, nil
}

func (p *Postgres) WalletByUser(userId int) (wallet.Wallet, error) {

	var result wallet.Wallet

	rows, err := p.Db.Query("SELECT * FROM user_wallet WHERE user_id = $1", userId)
	if err != nil {
		return result, nil
	}
	defer rows.Close()

	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return result, err
		}
		result = wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		}
	}
	return result, nil
}

func (p *Postgres) CreateWallet(createWallet wallet.CreateWallet) (wallet.Wallet, error) {
	var result wallet.Wallet
	sqlStr := "INSERT INTO user_wallet(user_id,user_name,wallet_name,wallet_type,balance) VALUES($1,$2,$3,$4,$5) " +
		"RETURNING id,user_id,user_name,wallet_name,wallet_type,balance,created_at"
	rows, err := p.Db.Query(sqlStr, createWallet.UserID, createWallet.UserName,
		createWallet.WalletName, createWallet.WalletType, createWallet.Balance)
	fmt.Printf("%v\n", rows)
	fmt.Printf("%v\n", err)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&result.ID,
			&result.UserID,
			&result.UserName,
			&result.WalletName,
			&result.WalletType,
			&result.Balance,
			&result.CreatedAt)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (p *Postgres) DeleteWallet(userID int) error {
	_, err := p.Db.Exec("DELETE FROM user_wallet WHERE user_id = $1", userID)
	if err != nil {
		return err
	}
	return nil

func (p *Postgres) UpdateWallet(updateWallet wallet.UpdateWallet) (wallet.Wallet, error) {
	var result wallet.Wallet
	sqlStr := "UPDATE user_wallet SET user_id=$1, user_name=$2, wallet_name=$3," +
		"wallet_type=$4, balance=$5, created_at=$6 WHERE id=$7 " +
		"RETURNING id, user_id, user_name, wallet_name, wallet_type, balance, created_at"
	rows, err := p.Db.Query(sqlStr, updateWallet.UserID, updateWallet.UserName, updateWallet.WalletName,
		updateWallet.WalletType, updateWallet.Balance, time.Now(),
		updateWallet.ID)
	if err != nil {
		return result, errors.New("unable to update row")
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&result.ID,
			&result.UserID,
			&result.UserName,
			&result.WalletName,
			&result.WalletType,
			&result.Balance,
			&result.CreatedAt)
		return result, err
	}
	return result, errors.New("unable to update row")
}
