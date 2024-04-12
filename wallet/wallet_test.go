package wallet

import (
	"encoding/json"
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
}
