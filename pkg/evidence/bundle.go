// Package evidence stores replay inputs, not authenticated ledger proofs.
package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"time"
)

const Schema = "ledgerparity-evidence/v1"
const MaxBytes = 32 << 20

type Data struct {
	Schema       string                   `json:"schema"`
	Tool         string                   `json:"tool"`
	Build        map[string]string        `json:"build"`
	SourceHashes map[string]string        `json:"source_sha256"`
	Assertions   map[string]string        `json:"assertions,omitempty"`
	App          string                   `json:"app"`
	Start        time.Time                `json:"start"`
	End          time.Time                `json:"end"`
	Options      engine.ReconcileOptions  `json:"options"`
	Internal     []types.InternalPayment  `json:"internal"`
	Observations []types.OnChainPayment   `json:"observations"`
	Report       *types.DiscrepancyReport `json:"report"`
}
type Bundle struct {
	SHA256 string `json:"sha256"`
	Data   Data   `json:"data"`
}

func Hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func BuildInfo() map[string]string {
	out := map[string]string{}
	if b, ok := debug.ReadBuildInfo(); ok {
		out["go"] = b.GoVersion
		out[b.Main.Path] = b.Main.Version
		for _, d := range b.Deps {
			v := d.Version
			if d.Replace != nil {
				v = "local/replaced build"
			}
			out[d.Path] = v
		}
		for _, s := range b.Settings {
			if s.Key == "vcs.revision" || s.Key == "vcs.modified" {
				out[s.Key] = s.Value
			}
		}
	}
	return out
}
func Seal(d Data) (Bundle, error) {
	d.Schema = Schema
	raw, err := json.Marshal(d)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{SHA256: Hash(raw), Data: d}, nil
}

// Write creates a new file only, so no evidence/input/report can be overwritten.
func Write(path string, b Bundle) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	if len(data) > MaxBytes {
		return fmt.Errorf("evidence exceeds %d bytes", MaxBytes)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}
func ReadBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, MaxBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > MaxBytes {
		return nil, fmt.Errorf("input exceeds %d bytes", MaxBytes)
	}
	return b, nil
}

// Decode rejects duplicate JSON keys, trailing values, unknown fields and deep
// nesting so integrity checks and consumers agree on the document's meaning.
func Decode(raw []byte, v any) error {
	if len(raw) > MaxBytes {
		return fmt.Errorf("input too large")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var walk func(int) error
	walk = func(depth int) error {
		if depth > 64 {
			return fmt.Errorf("JSON nesting too deep")
		}
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		if d, ok := tok.(json.Delim); ok {
			switch d {
			case '{':
				seen := map[string]bool{}
				for dec.More() {
					k, err := dec.Token()
					if err != nil {
						return err
					}
					key, ok := k.(string)
					if !ok || seen[key] {
						return fmt.Errorf("duplicate/invalid JSON key")
					}
					seen[key] = true
					if err = walk(depth + 1); err != nil {
						return err
					}
				}
			case '[':
				for dec.More() {
					if err := walk(depth + 1); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("invalid JSON delimiter")
			}
			_, err = dec.Token()
			return err
		}
		return nil
	}
	if err := walk(0); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return fmt.Errorf("expected one JSON value")
	}
	dec = json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
func Load(path string) (Bundle, error) {
	raw, err := ReadBytes(path)
	if err != nil {
		return Bundle{}, err
	}
	var b Bundle
	if err = Decode(raw, &b); err != nil {
		return Bundle{}, err
	}
	return b, nil
}
func Replay(b Bundle) (*types.DiscrepancyReport, error) {
	d := b.Data
	if d.Schema != Schema || d.Report == nil || d.Internal == nil || d.Observations == nil {
		return nil, fmt.Errorf("unsupported or incomplete evidence")
	}
	sealed, err := Seal(d)
	if err != nil {
		return nil, err
	}
	if sealed.SHA256 != b.SHA256 {
		return nil, fmt.Errorf("evidence integrity check failed")
	}
	for _, p := range d.Internal {
		if err := types.ValidateInternal(p); err != nil {
			return nil, err
		}
	}
	for _, p := range d.Observations {
		if err := types.ValidateOnChain(p); err != nil {
			return nil, err
		}
	}
	r := engine.NewReconciler(d.Options).Reconcile(d.App, d.Start, d.End, d.Internal, d.Observations)
	r.GeneratedAt = d.Report.GeneratedAt
	if !reflect.DeepEqual(r, d.Report) {
		return nil, fmt.Errorf("replayed report differs; engine revision or evidence changed")
	}
	return r, nil
}
