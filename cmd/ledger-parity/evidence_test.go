package main

import (
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-cli/pkg/evidence"
	"os"
	"path/filepath"
	"testing"
)

func TestEvidenceCLIReplay(t *testing.T) {
	path := prepare(t)
	bundle := filepath.Join(t.TempDir(), "evidence.json")
	code, out, stderr := execute("--config", path, "--format", "json", "--out", "-", "--bundle", bundle)
	if code != 3 || stderr != "" {
		t.Fatal(code, stderr)
	}
	before, err := evidence.Load(bundle)
	if err != nil {
		t.Fatal(err)
	}
	rawSource, err := os.ReadFile(filepath.Join(filepath.Dir(path), "internal.json"))
	if err != nil {
		t.Fatal(err)
	}
	if before.Data.SourceHashes["internal"] != evidence.Hash(rawSource) {
		t.Fatal("source bytes not captured exactly")
	}
	// Replay is independent of config and original input paths.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	again, replayed, stderr := execute("--replay", bundle)
	if again != code || replayed != out || stderr != "" {
		t.Fatal("replay differed", again, stderr)
	}
	if c, _, _ := execute("--replay", bundle, "--out", bundle); c != 1 {
		t.Fatal("overwrote evidence")
	}
	if c, _, _ := execute("--replay", bundle, "--demo"); c != 1 {
		t.Fatal("mixed modes")
	}
	b, err := evidence.Load(bundle)
	if err != nil {
		t.Fatal(err)
	}
	b.Data.Internal[0].Amount = "1"
	raw, _ := json.Marshal(b)
	os.WriteFile(bundle, raw, 0600)
	if c, _, _ := execute("--replay", bundle); c != 1 {
		t.Fatal("tampered evidence accepted")
	}
}
