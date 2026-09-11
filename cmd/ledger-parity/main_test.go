package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-cli/examples"
	"github.com/LedgerParity/ledger-parity-cli/pkg/config"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func prepare(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range []string{"config.json", "internal.json", "onchain.json"} {
		data, err := examples.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(dir, "config.json")
}
func execute(args ...string) (int, string, string) {
	var out, err bytes.Buffer
	code := run(context.Background(), args, &out, &err)
	return code, out.String(), err.String()
}
func TestDemoDeterministicAndMachineReadable(t *testing.T) {
	code, out, err := execute("--demo", "--format", "json", "--out", "-")
	if code != 3 || err != "" {
		t.Fatal(code, err)
	}
	var r types.DiscrepancyReport
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatal(err)
	}
	if r.TotalMatched != 1 || r.TotalDiscrepancies != 4 || r.TotalUnknown != 3 {
		t.Fatalf("%+v", r)
	}
	_, again, _ := execute("--demo", "--format", "json", "--out", "-")
	if again != out {
		t.Fatal("demo is not deterministic")
	}
}
func TestFlagsAndConfigOutputs(t *testing.T) {
	for _, args := range [][]string{nil, {"--wat"}, {"--demo", "--config", "x.json"}, {"--demo", "--format", "bogus"}, {"--demo", "--format", "both", "--out", "-"}} {
		if code, _, _ := execute(args...); code != 1 {
			t.Fatal(args, code)
		}
	}
	path := prepare(t)
	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Output.Format = "json"
	cfg.Output.FilePath = filepath.Join(t.TempDir(), "report.json")
	data, _ := json.Marshal(cfg)
	os.WriteFile(path, data, 0600)
	code, out, stderr := execute("--config", path)
	if code != 3 || out != "" || stderr != "" {
		t.Fatal(code, out, stderr)
	}
	if _, err = os.Stat(cfg.Output.FilePath); err != nil {
		t.Fatal("config output ignored", err)
	}
	code, out, _ = execute("--config", path, "--format", "table")
	if code != 3 || !bytes.Contains([]byte(out), []byte("Unknown: 3")) {
		t.Fatal(code, out)
	}
	if code, _, _ := execute("--config", path, "--out", cfg.TargetApp.SourcePath); code != 1 {
		t.Fatal("overwrote input")
	}
	if code, _, _ := execute("--config", path, "--out", filepath.Join(t.TempDir(), "missing", "out.json")); code != 1 {
		t.Fatal("output error hidden")
	}
}
func TestHorizonFailureDoesNotCreateMissingReport(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) }))
	defer s.Close()
	path := prepare(t)
	cfg, _ := config.LoadConfig(path)
	cfg.Stellar.OnChainPath = ""
	cfg.Stellar.HorizonURL = s.URL
	cfg.Output.Format = "json"
	cfg.Output.FilePath = filepath.Join(t.TempDir(), "report.json")
	data, _ := json.Marshal(cfg)
	os.WriteFile(path, data, 0600)
	code, out, err := execute("--config", path)
	if code != 1 || out != "" || err == "" {
		t.Fatal(code, out, err)
	}
	if _, err := os.Stat(cfg.Output.FilePath); !os.IsNotExist(err) {
		t.Fatal("failure wrote report")
	}
}
func TestEndToEndHorizonPayment(t *testing.T) {
	p := types.InternalPayment{ID: "1", Network: "test", OperationType: "payment", Sender: "A", Recipient: "B", Amount: "10", Asset: "XLM", AssetType: "native", Timestamp: mustTime(t, "2026-09-01T12:00:00Z"), Status: "completed", OperationID: "1"}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response any
		switch r.URL.Path {
		case "/":
			response = map[string]any{"network_passphrase": "test", "history_elder_ledger": 1, "history_latest_ledger": 2}
		case "/ledgers/1":
			response = map[string]string{"closed_at": "2026-08-31T12:00:00Z"}
		case "/ledgers/2":
			response = map[string]string{"closed_at": "2026-09-02T12:00:00Z"}
		default:
			rows := []map[string]any{}
			if r.URL.Query().Get("cursor") == "" {
				rows = append(rows, map[string]any{"id": "1", "paging_token": "1", "type": "payment", "transaction_successful": true, "created_at": p.Timestamp, "transaction_hash": "tx", "from": "A", "to": "B", "amount": "10.0000000", "asset_type": "native"})
			}
			response = map[string]any{"_embedded": map[string]any{"records": rows}}
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer s.Close()
	path := prepare(t)
	cfg, _ := config.LoadConfig(path)
	cfg.Stellar.OnChainPath = ""
	cfg.Stellar.HorizonURL = s.URL
	cfg.Stellar.Network = "test"
	cfg.Stellar.Accounts = []string{"A"}
	cfg.Output.Format = "json"
	cfg.Output.FilePath = "-"
	input, _ := json.Marshal([]types.InternalPayment{p})
	os.WriteFile(cfg.TargetApp.SourcePath, input, 0600)
	data, _ := json.Marshal(cfg)
	os.WriteFile(path, data, 0600)
	code, out, err := execute("--config", path)
	if code != 0 {
		t.Fatal(code, out, err)
	}
	var report types.DiscrepancyReport
	if json.Unmarshal([]byte(out), &report) != nil || report.TotalMatched != 1 || !report.Coverage.Complete {
		t.Fatal(out)
	}
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
