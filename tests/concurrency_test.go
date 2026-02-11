package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestConcurrencySameWallet(t *testing.T) {
	ts := newTestServer(t)

	id := uuid.New()

	deposit(t, ts.URL, id, 1_000_000)

	const workers = 1000
	const withdrawAmount = 1

	var wg sync.WaitGroup
	wg.Add(workers)

	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			status := withdrawStatus(ts.URL, id, withdrawAmount)
			if status != 200 {
				errCh <- fmt.Errorf("unexpected status: %d", status)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func deposit(t *testing.T, base string, id uuid.UUID, amount int64) {
	payload := map[string]any{"walletId": id.String(), "operationType": "DEPOSIT", "amount": amount}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(base+"/api/v1/wallet", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("deposit status: %d", resp.StatusCode)
	}
}

func withdrawStatus(base string, id uuid.UUID, amount int64) int {
	payload := map[string]any{"walletId": id.String(), "operationType": "WITHDRAW", "amount": amount}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(base+"/api/v1/wallet", "application/json", bytes.NewReader(b))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
