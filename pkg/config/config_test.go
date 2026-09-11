package config

import (
	"github.com/LedgerParity/ledger-parity-cli/examples"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStrictConfiguration(t *testing.T) {
	raw, _ := examples.Files.ReadFile("config.json")
	for _, bad := range []string{`{}`, `null`, string(raw) + `{}`, strings.Replace(string(raw), `"target_app"`, `"unknown"`, 1), strings.Replace(string(raw), `"start": "2026-09-01T11:00:00Z"`, `"start": "bad"`, 1)} {
		path := filepath.Join(t.TempDir(), "config.json")
		os.WriteFile(path, []byte(bad), 0600)
		if _, err := LoadConfig(path); err == nil {
			t.Fatal("accepted", bad)
		}
	}
	if _, err := LoadConfig("config.yaml"); err == nil {
		t.Fatal("YAML claim")
	}
}
