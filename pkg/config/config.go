package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type TargetAppConfig struct {
	Name       string `json:"name"`
	Format     string `json:"format"`
	SourcePath string `json:"source_path"`
	Complete   bool   `json:"complete"`
}
type StellarConfig struct {
	HorizonURL  string   `json:"horizon_url"`
	Network     string   `json:"network"`
	Accounts    []string `json:"accounts"`
	OnChainPath string   `json:"on_chain_path"`
}
type ReconciliationConfig struct {
	TimeframeToleranceSec int64     `json:"timeframe_tolerance_sec"`
	Start                 time.Time `json:"start"`
	End                   time.Time `json:"end"`
}
type OutputConfig struct {
	Format   string `json:"format"`
	FilePath string `json:"file_path"`
}
type Config struct {
	TargetApp      TargetAppConfig      `json:"target_app"`
	Stellar        StellarConfig        `json:"stellar"`
	Reconciliation ReconciliationConfig `json:"reconciliation"`
	Output         OutputConfig         `json:"output"`
}

func DefaultConfig() Config {
	return Config{TargetApp: TargetAppConfig{Name: "file", Format: "json"}, Reconciliation: ReconciliationConfig{TimeframeToleranceSec: 600}, Output: OutputConfig{Format: "table", FilePath: "discrepancy_report.json"}}
}
func LoadConfig(path string) (*Config, error) {
	if filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("configuration must be a .json file; YAML is unsupported")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg := DefaultConfig()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err = dec.Decode(&cfg); err != nil {
		return nil, err
	}
	var extra any
	if dec.Decode(&extra) != io.EOF {
		return nil, fmt.Errorf("expected one JSON configuration object")
	}
	if err = cfg.Validate(); err != nil {
		return nil, err
	}
	base := filepath.Dir(path)
	resolve := func(p string) string {
		if p != "" && p != "-" && !filepath.IsAbs(p) {
			return filepath.Join(base, p)
		}
		return p
	}
	cfg.TargetApp.SourcePath = resolve(cfg.TargetApp.SourcePath)
	cfg.Stellar.OnChainPath = resolve(cfg.Stellar.OnChainPath)
	cfg.Output.FilePath = resolve(cfg.Output.FilePath)
	return &cfg, nil
}
func (c *Config) Validate() error {
	if c.TargetApp.SourcePath == "" || c.TargetApp.Name == "" || (c.TargetApp.Format != "json" && c.TargetApp.Format != "csv") {
		return fmt.Errorf("target_app needs name, source_path and json/csv format")
	}
	if c.Stellar.Network == "" || len(c.Stellar.Accounts) == 0 {
		return fmt.Errorf("stellar network passphrase and monitored accounts required")
	}
	for _, a := range c.Stellar.Accounts {
		if a == "" {
			return fmt.Errorf("empty monitored account")
		}
	}
	if (c.Stellar.HorizonURL == "") == (c.Stellar.OnChainPath == "") {
		return fmt.Errorf("choose exactly one of horizon_url or on_chain_path")
	}
	r := c.Reconciliation
	if r.Start.IsZero() || r.End.IsZero() || r.End.Before(r.Start) || r.TimeframeToleranceSec < 0 || r.TimeframeToleranceSec > 86400 {
		return fmt.Errorf("ordered start/end and tolerance 0..86400 seconds required")
	}
	if c.Output.Format != "table" && c.Output.Format != "json" && c.Output.Format != "both" {
		return fmt.Errorf("format must be table, json or both")
	}
	if c.Output.Format != "table" && c.Output.FilePath == "" {
		return fmt.Errorf("JSON output path required")
	}
	return nil
}
