package main

import (
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-cli/examples"
	"github.com/LedgerParity/ledger-parity-cli/pkg/config"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSDPEndToEnd(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"config.json", "payments.csv", "onchain.json"} {
		b, err := examples.Files.ReadFile("sdp/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(dir, "config.json")
	bundle := filepath.Join(dir, "bundle.json")
	code, out, stderr := execute("--config", path, "--bundle", bundle)
	if code != 3 || stderr != "" {
		t.Fatal(code, stderr)
	}
	var r types.DiscrepancyReport
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatal(err)
	}
	if r.TotalMatched != 1 || r.TotalDiscrepancies != 1 || r.TotalUnknown != 1 || r.Coverage.InternalComplete {
		t.Fatalf("%+v", r)
	}
	b, _ := os.ReadFile(bundle)
	if strings.Contains(string(b), "private@example") || strings.Contains(out, "private@example") {
		t.Fatal("contact field leaked")
	}
	if c, replayed, e := execute("--replay", bundle); c != code || e != "" || replayed != out {
		t.Fatal("SDP replay", c, e)
	}
	cfg, err := config.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg.TargetApp.Complete = true
	if cfg.Validate() == nil {
		t.Fatal("complete SDP export accepted")
	}
	cfg.TargetApp.Complete = false
	cfg.TargetApp.SDP.Sender = "other"
	if cfg.Validate() == nil {
		t.Fatal("unmonitored sender")
	}
	// A later-page/HTTP failure must not create an evidence file either.
	cfg.TargetApp.SDP.Sender = "GDEMO_SENDER"
	cfg.Stellar.OnChainPath = "missing.json"
	raw, _ := json.Marshal(cfg)
	os.WriteFile(path, raw, 0600)
	missingBundle := filepath.Join(dir, "failure.json")
	if c, _, _ := execute("--config", path, "--bundle", missingBundle); c != 1 {
		t.Fatal("missing observation succeeded")
	}
	if _, err = os.Stat(missingBundle); !os.IsNotExist(err) {
		t.Fatal("failure wrote evidence")
	}
}
