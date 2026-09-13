package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/LedgerParity/ledger-parity-cli/pkg/evidence"
)

// HashReport computes the SHA-256 of a report JSON file.
func HashReport(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("read report: %w", err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash report: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Proof represents a stored verification proof.
type Proof struct {
	Hash      string `json:"hash"`
	Owner     string `json:"owner"`
	Network   string `json:"network"`
	Timestamp string `json:"timestamp"`
	Contract  string `json:"contract,omitempty"`
}

// SaveProof writes a verification proof to a JSON file.
func SaveProof(path string, proof *Proof) error {
	if err := validate(proof); err != nil {
		return err
	}
	data, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return err
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

// LoadProof reads a verification proof from a JSON file.
func LoadProof(path string) (*Proof, error) {
	data, err := evidence.ReadBytes(path)
	if err != nil {
		return nil, err
	}
	var p Proof
	if err := evidence.Decode(data, &p); err != nil {
		return nil, err
	}
	if err := validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func validate(p *Proof) error {
	if p == nil {
		return fmt.Errorf("proof is required")
	}
	decoded, err := hex.DecodeString(p.Hash)
	if err != nil || len(decoded) != sha256.Size || p.Hash != strings.ToLower(p.Hash) {
		return fmt.Errorf("proof hash must be 64 lowercase hexadecimal characters")
	}
	return nil
}
