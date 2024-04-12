package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Wallets() ([]Wallet, error)
	WalletsByType(walletType string) ([]Wallet, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// WalletHandler
//
//		@Summary		Get all wallets
//		@Description	Get all wallets
//		@Tags			wallet
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/wallets [get]
//		@Failure		500	{object}	Err
//	 	@Param          wallet_type query string false "wallet type"
func (h *Handler) WalletHandler(c echo.Context) error {
	var wallets []Wallet
	var err error
	walletType := c.QueryParam("wallet_type")
	if walletType != "" {
		wallets, err = h.store.WalletsByType(walletType)
	} else {
		wallets, err = h.store.Wallets()
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}
