package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/LedgerParity/ledger-parity-cli/examples"
	"github.com/LedgerParity/ledger-parity-cli/pkg/config"
	"github.com/LedgerParity/ledger-parity-cli/pkg/output"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/file"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/ingest"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const Version = "0.2.0-preview"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}

// Exit 0: complete match; 1: runtime/config error; 2: discrepancies; 3: unresolved evidence.
func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("ledger-parity", flag.ContinueOnError)
	flags.SetOutput(stderr)
	cfgPath := flags.String("config", "", "JSON configuration file")
	format := flags.String("format", "", "Override format: table, json, both")
	out := flags.String("out", "", "Override JSON report path; - writes JSON to stdout")
	demo := flags.Bool("demo", false, "Run deterministic offline fixtures")
	version := flags.Bool("version", false, "Show version")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	fail := func(err error) int { fmt.Fprintln(stderr, err); return 1 }
	if flags.NArg() != 0 {
		return fail(fmt.Errorf("unexpected positional arguments"))
	}
	if *version {
		fmt.Fprintln(stdout, Version)
		return 0
	}
	if (*cfgPath == "") == (!*demo) {
		return fail(fmt.Errorf("choose exactly one of --config or --demo"))
	}
	if *demo {
		dir, err := os.MkdirTemp("", "ledger-parity-demo-")
		if err != nil {
			return fail(err)
		}
		defer os.RemoveAll(dir)
		for _, name := range []string{"config.json", "internal.json", "onchain.json"} {
			data, err := examples.Files.ReadFile(name)
			if err != nil {
				return fail(err)
			}
			if err = os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
				return fail(err)
			}
		}
		*cfgPath = filepath.Join(dir, "config.json")
	}
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		return fail(err)
	}
	if *format != "" {
		cfg.Output.Format = *format
	}
	if *out != "" {
		cfg.Output.FilePath = *out
	} else if *demo {
		cfg.Output.FilePath = "discrepancy_report.json"
	}
	if err = cfg.Validate(); err != nil {
		return fail(err)
	}
	if cfg.Output.Format == "both" && cfg.Output.FilePath == "-" {
		return fail(fmt.Errorf("use --format json for JSON stdout"))
	}
	if cfg.Output.FilePath != "-" && cfg.Output.Format != "table" {
		dest, _ := filepath.Abs(cfg.Output.FilePath)
		for _, path := range []string{*cfgPath, cfg.TargetApp.SourcePath, cfg.Stellar.OnChainPath} {
			if path == "" {
				continue
			}
			src, _ := filepath.Abs(path)
			destInfo, destErr := os.Stat(dest)
			srcInfo, srcErr := os.Stat(src)
			if dest == src || (destErr == nil && srcErr == nil && os.SameFile(destInfo, srcInfo)) {
				return fail(fmt.Errorf("report path must not overwrite an input"))
			}
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	report, err := reconcile(ctx, cfg)
	if err != nil {
		return fail(err)
	}
	if *demo {
		report.GeneratedAt = cfg.Reconciliation.End
	}
	f := output.NewFormatter(stdout)
	if cfg.Output.Format == "table" || cfg.Output.Format == "both" {
		if err = f.RenderTerminalTable(report); err != nil {
			return fail(err)
		}
	}
	if cfg.Output.Format == "json" || cfg.Output.Format == "both" {
		if err = f.ExportJSON(report, cfg.Output.FilePath); err != nil {
			return fail(err)
		}
	}
	if report.TotalUnknown > 0 {
		return 3
	}
	if report.TotalDiscrepancies > 0 {
		return 2
	}
	return 0
}
func reconcile(ctx context.Context, cfg *config.Config) (*types.DiscrepancyReport, error) {
	start, end := cfg.Reconciliation.Start, cfg.Reconciliation.End
	ips, err := file.NewFileConnector(cfg.TargetApp.SourcePath, cfg.TargetApp.Format, cfg.TargetApp.Name).FetchInternalPayments(ctx, connector.Filter{TimeStart: start, TimeEnd: end})
	if err != nil {
		return nil, err
	}
	for _, p := range ips {
		if p.Network != cfg.Stellar.Network {
			return nil, fmt.Errorf("internal record network differs from config")
		}
		found := false
		for _, a := range cfg.Stellar.Accounts {
			if a == p.Sender || a == p.Recipient {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("internal record outside configured account scope")
		}
	}
	var observed ingest.FetchResult
	if cfg.Stellar.OnChainPath != "" {
		f, err := os.Open(cfg.Stellar.OnChainPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		dec.DisallowUnknownFields()
		if err = dec.Decode(&observed); err != nil {
			return nil, err
		}
		var extra any
		if dec.Decode(&extra) != io.EOF {
			return nil, fmt.Errorf("expected one offline observation object")
		}
		if observed.Payments == nil {
			return nil, fmt.Errorf("payments must be an array")
		}
		if observed.Coverage.Network != cfg.Stellar.Network {
			return nil, fmt.Errorf("offline coverage network mismatch")
		}
		for _, a := range cfg.Stellar.Accounts {
			found := false
			for _, b := range observed.Coverage.Accounts {
				if a == b {
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("offline fixture does not cover configured accounts")
			}
		}
		observed.Coverage.Source = "offline-file"
		observed.Coverage.Reason = "Offline caller assertion, not live verification: " + observed.Coverage.Reason
	} else {
		h := ingest.NewHorizonIngestor(cfg.Stellar.HorizonURL)
		h.Network = cfg.Stellar.Network
		tolerance := time.Duration(cfg.Reconciliation.TimeframeToleranceSec) * time.Second
		observed, err = h.Fetch(ctx, cfg.Stellar.Accounts, start.Add(-tolerance), end.Add(tolerance))
		if err != nil {
			return nil, fmt.Errorf("ingestion failed; no report produced: %w", err)
		}
	}
	scoped := []types.OnChainPayment{}
	for _, p := range observed.Payments {
		if err = types.ValidateOnChain(p); err != nil {
			return nil, err
		}
		if p.Network != cfg.Stellar.Network {
			return nil, fmt.Errorf("observed network mismatch")
		}
		found := false
		for _, a := range cfg.Stellar.Accounts {
			if a == p.Account || a == p.Destination {
				found = true
			}
		}
		if found {
			scoped = append(scoped, p)
		}
	}
	observed.Coverage.Accounts = cfg.Stellar.Accounts
	observed.Coverage.InternalComplete = cfg.TargetApp.Complete
	opts := engine.ReconcileOptions{TimeframeToleranceSec: cfg.Reconciliation.TimeframeToleranceSec, Coverage: observed.Coverage}
	return engine.NewReconciler(opts).Reconcile(cfg.TargetApp.Name, start, end, ips, scoped), nil
}
