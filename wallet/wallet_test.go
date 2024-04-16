package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type StubStorer struct {
	wallets []Wallet
	err     error
}

func (s StubStorer) CreateWallet(createWallet CreateWallet) (Wallet, error) {
	result := Wallet{
		ID:         1,
		UserID:     createWallet.UserID,
		UserName:   createWallet.UserName,
		WalletName: createWallet.WalletName,
		WalletType: createWallet.WalletType,
		Balance:    createWallet.Balance,
		CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
	}
	_ = append(s.wallets, result)
	return result, nil
}

func (s StubStorer) UpdateWallet(updateWallet UpdateWallet) (Wallet, error) {

	for _, scanWallet := range s.wallets {
		if scanWallet.ID == updateWallet.ID {
			scanWallet.UserID = updateWallet.UserID
			scanWallet.UserName = updateWallet.UserName
			scanWallet.WalletName = updateWallet.WalletName
			scanWallet.WalletType = updateWallet.WalletType
			scanWallet.Balance = updateWallet.Balance
			scanWallet.CreatedAt = time.Now()
			return scanWallet, nil
		}
	}

	return Wallet{}, errors.New("Unable to find update row!")
}

func (s StubStorer) Wallets() ([]Wallet, error) {
	return s.wallets, s.err
}

func (s StubStorer) WalletsByType(walletType string) ([]Wallet, error) {
	var result []Wallet
	for _, wallet := range s.wallets {
		if wallet.WalletType == walletType {
			result = append(result, wallet)
		}
	}
	return result, s.err
}

func (s StubStorer) WalletByUser(userId int) (Wallet, error) {
	var result Wallet
	for _, wallet := range s.wallets {
		if wallet.UserID == userId {
			result = wallet
		}
	}
	return result, s.err
}

type ErrorMessage struct {
	Message string
}

func TestWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		w := New(StubStorer{err: echo.ErrInternalServerError})

		w.WalletHandler(c)

		if res.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, res.Code)
		}
		var errorMessage ErrorMessage
		if err := json.Unmarshal(res.Body.Bytes(), &errorMessage); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		wantErrMsg := "code=500, message=Internal Server Error"
		if errorMessage.Message != wantErrMsg {
			t.Errorf("expected error message %q but got %q", wantErrMsg, errorMessage.Message)
		}
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		want := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   "Jame Bonds",
				WalletName: "Jame Wallet",
				WalletType: "Saving",
				Balance:    100.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "Jane Bonds",
				WalletName: "Jane Wallet",
				WalletType: "Saving",
				Balance:    500.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
		}
		w := New(StubStorer{wallets: want})

		w.WalletHandler(c)

		gotJson := res.Body.Bytes()
		var got []Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given user able to getting wallet by type should return correct list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		q := req.URL.Query()
		q.Add("wallet_type", "Saving")
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		body := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   "Jame Bonds",
				WalletName: "Jame Wallet",
				WalletType: "Saving",
				Balance:    100.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "Jane Bonds",
				WalletName: "Jane Wallet",
				WalletType: "Saving1",
				Balance:    500.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
		}
		want := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   "Jame Bonds",
				WalletName: "Jame Wallet",
				WalletType: "Saving",
				Balance:    100.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
		}
		w := New(StubStorer{wallets: body})

		w.WalletHandler(c)

		gotJson := res.Body.Bytes()
		var got []Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given user able to getting wallet by user id should return correct wallet", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		c.SetPath("/users/:id/wallets")
		c.SetParamNames("id")
		c.SetParamValues("2")
		body := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   "Jame Bonds",
				WalletName: "Jame Wallet",
				WalletType: "Saving",
				Balance:    100.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "Jane Bonds",
				WalletName: "Jane Wallet",
				WalletType: "Saving1",
				Balance:    500.00,
				CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
			},
		}
		want := Wallet{
			ID:         2,
			UserID:     2,
			UserName:   "Jane Bonds",
			WalletName: "Jane Wallet",
			WalletType: "Saving1",
			Balance:    500.00,
			CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
		}
		w := New(StubStorer{wallets: body})

		w.WalletHandlerByUser(c)

		gotJson := res.Body.Bytes()
		var got Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given user able to create wallet should return created wallet", func(t *testing.T) {
		createWallet := CreateWallet{
			UserID:     14,
			UserName:   "Jame",
			WalletName: "Jame Wallet",
			WalletType: "Savings",
			Balance:    1499.00,
		}
		body, err := json.Marshal(createWallet)
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		want := Wallet{
			ID:         1,
			UserID:     createWallet.UserID,
			UserName:   createWallet.UserName,
			WalletName: createWallet.WalletName,
			WalletType: createWallet.WalletType,
			Balance:    createWallet.Balance,
			CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
		}
		w := New(StubStorer{wallets: []Wallet{}})

		w.CreateWallet(c)

		gotJson := res.Body.Bytes()
		var got Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given user able to update wallet should return upated wallet", func(t *testing.T) {
		updateWallet := UpdateWallet{
			ID:         1,
			UserID:     14,
			UserName:   "Jame",
			WalletName: "Jame Wallet",
			WalletType: "Savings",
			Balance:    1499.00,
		}
		body, err := json.Marshal(updateWallet)
		if err != nil {
			t.Errorf("Unable to create body request, error: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)

		want := Wallet{
			ID:         updateWallet.ID,
			UserID:     updateWallet.UserID,
			UserName:   updateWallet.UserName,
			WalletName: updateWallet.WalletName,
			WalletType: updateWallet.WalletType,
			Balance:    updateWallet.Balance,
			CreatedAt:  time.Date(2024, 04, 12, 10, 45, 16, 0, time.UTC),
		}
		w := New(StubStorer{wallets: []Wallet{want}})

		w.UpdateWallet(c)

		gotJson := res.Body.Bytes()
		var got Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("Unable to unmarshal json: %v", err)
		}
		if got.ID != want.ID &&
			got.UserID != want.UserID &&
			got.UserName != want.UserName &&
			got.WalletName != want.WalletName &&
			got.WalletType != want.WalletType &&
			got.Balance != want.Balance {
			t.Errorf("expected %v but got %v", want, got)
		}
		if got.CreatedAt == want.CreatedAt {
			t.Errorf("CreatedAt is not changed old=%v, new=%v", want, got)
		}
	})
}
