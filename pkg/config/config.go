package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TargetAppConfig defines internal payment source details.
type TargetAppConfig struct {
	Name       string `json:"name" yaml:"name"`                 // e.g. "stellopay", "facilpay", "file"
	Format     string `json:"format" yaml:"format"`             // "json" or "csv"
	SourcePath string `json:"source_path" yaml:"source_path"` // Path to export file or DB connection
}

// StellarConfig defines Horizon / Soroban RPC connection settings.
type StellarConfig struct {
	HorizonURL string   `json:"horizon_url" yaml:"horizon_url"`
	Accounts   []string `json:"accounts" yaml:"accounts"`
}

// ReconciliationConfig defines tolerance parameters.
type ReconciliationConfig struct {
	TimeframeToleranceSec int64 `json:"timeframe_tolerance_sec" yaml:"timeframe_tolerance_sec"`
	IgnoreFailedOnChain   bool  `json:"ignore_failed_on_chain" yaml:"ignore_failed_on_chain"`
}

// OutputConfig defines report formatting & destination.
type OutputConfig struct {
	Format   string `json:"format" yaml:"format"`       // "table", "json", "both"
	FilePath string `json:"file_path" yaml:"file_path"` // Export JSON report path
}

// Config represents the complete CLI configuration.
type Config struct {
	TargetApp      TargetAppConfig      `json:"target_app" yaml:"target_app"`
	Stellar        StellarConfig        `json:"stellar" yaml:"stellar"`
	Reconciliation ReconciliationConfig `json:"reconciliation" yaml:"reconciliation"`
	Output         OutputConfig         `json:"output" yaml:"output"`
}

// DefaultConfig returns default runtime configuration.
func DefaultConfig() Config {
	return Config{
		TargetApp: TargetAppConfig{
			Name:       "stellopay",
			Format:     "json",
			SourcePath: "stellopay_export.json",
		},
		Stellar: StellarConfig{
			HorizonURL: "https://horizon-testnet.stellar.org",
			Accounts:   []string{},
		},
		Reconciliation: ReconciliationConfig{
			TimeframeToleranceSec: 600,
			IgnoreFailedOnChain:   true,
		},
		Output: OutputConfig{
			Format:   "table",
			FilePath: "discrepancy_report.json",
		},
	}
}

// LoadConfig reads configuration from a JSON or YAML file.
func LoadConfig(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath == "" {
		return &cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	ext := strings.ToLower(filepath.Ext(configPath))
	if ext == ".json" || ext == "" {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed parsing json config: %w", err)
		}
	} else {
		// Basic key-value / json fallback for zero external dependency
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed parsing config: %w", err)
		}
	}

	return &cfg, nil
}
