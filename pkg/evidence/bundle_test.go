package evidence

import (
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"path/filepath"
	"testing"
	"time"
)

func TestReplayAndIntegrity(t *testing.T) {
	a := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	d := Data{Tool: "test", App: "empty", Start: a, End: a.Add(time.Hour), Options: engine.DefaultOptions(), Internal: []types.InternalPayment{}, Observations: []types.OnChainPayment{}, SourceHashes: map[string]string{"internal": Hash([]byte("[]"))}}
	d.Report = engine.NewReconciler(d.Options).Reconcile(d.App, d.Start, d.End, d.Internal, d.Observations)
	b, err := Seal(d)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "bundle.json")
	if err = Write(path, b); err != nil {
		t.Fatal(err)
	}
	if err = Write(path, b); err == nil {
		t.Fatal("overwrote evidence")
	}
	b, err = Load(path)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Replay(b)
	if err != nil || r.TotalUnknown != 1 {
		t.Fatal("unknown lost", err)
	}
	b.Data.Tool = "changed"
	if _, err = Replay(b); err == nil {
		t.Fatal("tampering accepted")
	}
	b, _ = Seal(b.Data)
	b.Data.Report.TotalUnknown = 0
	b, _ = Seal(b.Data)
	if _, err = Replay(b); err == nil {
		t.Fatal("false report accepted")
	}
	if Hash([]byte("[]")) == Hash([]byte("[] ")) {
		t.Fatal("byte change hidden")
	}
}
func TestStrictEvidenceJSON(t *testing.T) {
	for _, s := range []string{`{"sha256":"a","sha256":"b"}`, `{} {}`, `{"unknown":1}`, `{"data":{"build":{"a":"1","a":"2"}}}`} {
		var b Bundle
		if err := Decode([]byte(s), &b); err == nil {
			t.Fatal("accepted", s)
		}
	}
	b, _ := json.Marshal(Bundle{})
	var out Bundle
	if err := Decode(b, &out); err != nil {
		t.Fatal(err)
	}
}
