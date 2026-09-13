package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LedgerParity/ledger-parity-cli/pkg/verify"
)

func TestVerificationFreshReportAndTampering(t *testing.T) {
	for _, stale := range []bool{false, true} {
		t.Run(map[bool]string{false: "new", true: "stale"}[stale], func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "report.json")
			if stale {
				if err := os.WriteFile(path, []byte("old report"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			code, _, msg := execute("--demo", "--format", "json", "--out", path, "--verify", "auto")
			if code != 3 {
				t.Fatal(code, msg)
			}
			proof, err := verify.LoadProof(path + ".proof.json")
			if err != nil {
				t.Fatal(err)
			}
			if proof.Timestamp == "" {
				t.Fatal("missing creation time")
			}
			if code, out, msg := execute("--verify-check", path+".proof.json", path); code != 0 || !strings.Contains(out, "VERIFIED") {
				t.Fatal(code, out, msg)
			}
			if err := os.WriteFile(path, []byte("changed"), 0600); err != nil {
				t.Fatal(err)
			}
			if code, _, msg := execute("--verify-check", path+".proof.json", path); code != 1 || !strings.Contains(msg, "MISMATCH") {
				t.Fatal(code, msg)
			}
		})
	}
}

func TestVerificationRejectsUnsafeOutputs(t *testing.T) {
	cfg := prepare(t)
	dir := filepath.Dir(cfg)
	for _, name := range []string{"config.json", "internal.json", "onchain.json", "report.json", "bundle.json"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			before, _ := os.ReadFile(path)
			code, _, msg := execute("--config", cfg, "--format", "json", "--out", filepath.Join(dir, "report.json"), "--bundle", filepath.Join(dir, "bundle.json"), "--verify", path)
			if code != 1 {
				t.Fatal(code, msg)
			}
			after, _ := os.ReadFile(path)
			if string(before) != string(after) {
				t.Fatal("overwritten file")
			}
		})
	}
	for _, args := range [][]string{
		{"--demo", "--format", "table", "--verify", "auto"},
		{"--demo", "--format", "json", "--out", "-", "--verify", "auto"},
		{"--replay", "missing", "--verify", "auto"},
		{"--verify-check", "missing"},
		{"--verify-check", "missing", "one", "two"},
		{"--verify-check", "missing", "--demo", "one"},
		{"--verify-check", "missing", "--version", "one"},
	} {
		if code, _, msg := execute(args...); code != 1 {
			t.Fatal(args, code, msg)
		}
	}
}

func TestVerificationMalformedProofAndWriteFailure(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.json")
	proof := filepath.Join(dir, "proof.json")
	if err := os.WriteFile(report, []byte("report"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{`{}`, `null`, `{"hash":"short"}`, `{"hash":"a","hash":"b"}`} {
		if err := os.WriteFile(proof, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		if code, _, _ := execute("--verify-check", proof, report); code != 1 {
			t.Fatal(raw, code)
		}
	}
	newProof := filepath.Join(dir, "new-proof.json")
	if code, _, msg := execute("--demo", "--format", "json", "--out", dir, "--verify", newProof); code != 1 {
		t.Fatal(code, msg)
	}
	if _, err := os.Stat(newProof); !os.IsNotExist(err) {
		t.Fatal("proof created after failed report write", err)
	}
	if err := verify.SaveProof(proof, &verify.Proof{Hash: strings.Repeat("a", 64)}); err == nil {
		t.Fatal("existing proof overwritten")
	}
}
