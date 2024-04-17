package wallet

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Wallets() ([]Wallet, error)
	WalletsByType(walletType string) ([]Wallet, error)
	WalletByUser(userID int) (Wallet, error)
	CreateWallet(createWallet CreateWallet) (Wallet, error)
	DeleteWallet(userID int) error
	UpdateWallet(updateWallet UpdateWallet) (Wallet, error)
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

// WalletHandlerByUser
//
//		@Summary		Get wallet by user Id
//		@Description	Get wallet by user Id
//		@Tags			wallet
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/users/{id}/wallets [get]
//		@Failure		500	{object}	Err
//	 	@Param          id path int true "User ID"
func (h *Handler) WalletHandlerByUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Unable to find wallet!"})
	}
	result, err := h.store.WalletByUser(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Unable to find wallet!"})
	}
	if result.UserID != userId {
		return c.JSON(http.StatusNotFound, Err{Message: "Unable to find wallet!"})
	}
	return c.JSON(http.StatusOK, result)
}

// CreateWallet
//
//		@Summary		Create wallet
//		@Description	Create wallet
//		@Tags			wallet
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/wallets [post]
//		@Failure		500	{object}	Err
//	 	@Param 			CreateWallet body CreateWallet true "Body for create wallet"
func (h *Handler) CreateWallet(c echo.Context) error {
	var createWallet CreateWallet
	if err := c.Bind(&createWallet); err != nil {
		return err
	}
	result, err := h.store.CreateWallet(createWallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}

// DeleteWallet
//
//		@Summary		Delete wallet by user Id
//		@Description	Delete wallet by user Id
//		@Tags			wallet
//		@Accept			json
//		@Produce		plain
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/users/{id}/wallets [delete]
//		@Failure		500	{object}	Err
//	 	@Param          id path int true "User ID"
func (h *Handler) DeleteWallet(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Unable to find wallet!"})
	}
	err = h.store.DeleteWallet(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.String(http.StatusOK, "Delete Success")
}

// UpdateWallet
//
//		@Summary		Update wallet
//		@Description	Update wallet
//		@Tags			wallet
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/wallets [patch]
//		@Failure		500	{object}	Err
//	 	@Param 			UpdateWallet body UpdateWallet true "Body for update wallet"
func (h *Handler) UpdateWallet(c echo.Context) error {
	var updateWallet UpdateWallet
	if err := c.Bind(&updateWallet); err != nil {
		return err
	}
	result, err := h.store.UpdateWallet(updateWallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
