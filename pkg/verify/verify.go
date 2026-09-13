package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// HashReport computes the SHA-256 of a report JSON file.
func HashReport(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read report: %w", err)
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
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
	data, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadProof reads a verification proof from a JSON file.
func LoadProof(path string) (*Proof, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Proof
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
