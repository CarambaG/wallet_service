package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"TestProject_itk/internal/api"
	"TestProject_itk/internal/db"
	"TestProject_itk/internal/repository"
	"TestProject_itk/internal/service"

	"github.com/google/uuid"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	dsn := getenv("POSTGRES_DSN", "")
	if dsn == "" {
		t.Skip("POSTGRES_DSN is not set")
	}

	pool, err := db.NewPostgresPool(context.Background(), dsn)
	if err != nil {
		t.Fatalf("db: %v", err)
	}

	repo := repository.NewWalletRepo(pool)
	svc := service.NewWalletService(repo)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	ts := httptest.NewServer(r)
	t.Cleanup(func() {
		ts.Close()
		pool.Close()
	})

	return ts
}

func TestDepositAndGet(t *testing.T) {
	ts := newTestServer(t)

	id := uuid.New()

	body := map[string]any{
		"walletId":      id.String(),
		"operationType": "DEPOSIT",
		"amount":        1000,
	}
	b, _ := json.Marshal(body)

	resp, err := http.Post(ts.URL+"/api/v1/wallet", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %d", resp.StatusCode)
	}

	resp2, err := http.Get(ts.URL + "/api/v1/wallets/" + id.String())
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status: %d", resp2.StatusCode)
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
