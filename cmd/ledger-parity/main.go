package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/LedgerParity/ledger-parity-cli/examples"
	"github.com/LedgerParity/ledger-parity-cli/pkg/config"
	"github.com/LedgerParity/ledger-parity-cli/pkg/evidence"
	"github.com/LedgerParity/ledger-parity-cli/pkg/output"
	"github.com/LedgerParity/ledger-parity-cli/pkg/verify"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/file"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/sdp"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/ingest"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const Version = "0.3.1-preview"

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
	bundlePath := flags.String("bundle", "", "Create a new replayable evidence file")
	replayPath := flags.String("replay", "", "Verify and replay evidence offline; no configuration/network access")
	verifyPath := flags.String("verify", "", "After report, store SHA-256 hash (writes .proof.json)")
	verifyCheck := flags.String("verify-check", "", "Verify a report against a saved proof file")
	version := flags.Bool("version", false, "Show version")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	fail := func(err error) int { fmt.Fprintln(stderr, err); return 1 }
	if *verifyCheck == "" && flags.NArg() != 0 {
		return fail(fmt.Errorf("unexpected positional arguments"))
	}
	if *version && *verifyCheck == "" {
		fmt.Fprintln(stdout, Version)
		return 0
	}
	if *verifyCheck != "" {
		if *demo || *cfgPath != "" || *bundlePath != "" || *replayPath != "" || *verifyPath != "" || *format != "" || *out != "" || *version {
			return fail(fmt.Errorf("verify-check cannot be combined with other modes or output flags"))
		}
		if flags.NArg() != 1 {
			return fail(fmt.Errorf("supply exactly one report JSON path to verify"))
		}
		proof, err := verify.LoadProof(*verifyCheck)
		if err != nil {
			return fail(err)
		}
		hash, err := verify.HashReport(flags.Arg(0))
		if err != nil {
			return fail(err)
		}
		if hash == proof.Hash {
			fmt.Fprintf(stdout, "VERIFIED: report matches proof (hash=%s)\n", hash[:16])
			return 0
		}
		fmt.Fprintf(stderr, "MISMATCH: report hash %s != proof hash %s\n", hash[:16], proof.Hash[:16])
		return 1
	}
	if *replayPath != "" {
		if *verifyPath != "" {
			return fail(fmt.Errorf("replay cannot be combined with verify"))
		}
		if *demo || *cfgPath != "" || *bundlePath != "" {
			return fail(fmt.Errorf("replay cannot be combined with config, demo or bundle"))
		}
		b, err := evidence.Load(*replayPath)
		if err != nil {
			return fail(err)
		}
		r, err := evidence.Replay(b)
		if err != nil {
			return fail(err)
		}
		if *format == "" {
			*format = "json"
		}
		if *out == "" {
			*out = "-"
		}
		if *out != "-" && samePath(*out, *replayPath) {
			return fail(fmt.Errorf("report must not overwrite evidence"))
		}
		return render(r, *format, *out, stdout, stderr)
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
	if *verifyPath != "" {
		if cfg.Output.FilePath == "-" || cfg.Output.Format == "table" {
			return fail(fmt.Errorf("--verify requires a JSON report file (json or both format)"))
		}
		if *verifyPath == "auto" {
			*verifyPath = cfg.Output.FilePath + ".proof.json"
		}
		for _, path := range []string{*cfgPath, cfg.TargetApp.SourcePath, cfg.Stellar.OnChainPath, cfg.Output.FilePath, *bundlePath} {
			if path != "" && samePath(*verifyPath, path) {
				return fail(fmt.Errorf("proof path must differ from inputs, report and bundle"))
			}
		}
		if _, err := os.Lstat(*verifyPath); !os.IsNotExist(err) {
			return fail(fmt.Errorf("proof path must be a new file"))
		}
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
	if *bundlePath != "" && samePath(*bundlePath, cfg.Output.FilePath) {
		return fail(fmt.Errorf("bundle and report paths must differ"))
	}
	var captured evidence.Data
	report, err := reconcileCapture(ctx, cfg, &captured)
	if err != nil {
		return fail(err)
	}
	if *demo {
		report.GeneratedAt = cfg.Reconciliation.End
	}
	if *bundlePath != "" {
		b, err := evidence.Seal(captured)
		if err != nil {
			return fail(err)
		}
		if err = evidence.Write(*bundlePath, b); err != nil {
			return fail(err)
		}
		if samePath(*bundlePath, cfg.Output.FilePath) {
			return fail(fmt.Errorf("report resolves to evidence file; refusing overwrite"))
		}
	}
	code := render(report, cfg.Output.Format, cfg.Output.FilePath, stdout, stderr)
	if code == 1 {
		return code
	}
	if *verifyPath != "" {
		reportPath := cfg.Output.FilePath
		hash, err := verify.HashReport(reportPath)
		if err != nil {
			return fail(err)
		}
		proof := &verify.Proof{
			Hash:      hash,
			Owner:     "ledger-parity",
			Network:   cfg.Stellar.Network,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		proofPath := *verifyPath
		if err := verify.SaveProof(proofPath, proof); err != nil {
			return fail(err)
		}
		fmt.Fprintf(stderr, "Verification proof saved: %s (hash=%s)\n", proofPath, hash[:16])
	}
	return code
}

func samePath(a, b string) bool {
	x, _ := filepath.Abs(a)
	y, _ := filepath.Abs(b)
	xi, xe := os.Stat(x)
	yi, ye := os.Stat(y)
	return x == y || (runtime.GOOS == "windows" && strings.EqualFold(x, y)) || (xe == nil && ye == nil && os.SameFile(xi, yi))
}
func render(report *types.DiscrepancyReport, format, path string, stdout, stderr io.Writer) int {
	fail := func(err error) int { fmt.Fprintln(stderr, err); return 1 }
	if format != "table" && format != "json" && format != "both" {
		return fail(fmt.Errorf("invalid output format"))
	}
	if format == "both" && path == "-" {
		return fail(fmt.Errorf("use json format for stdout"))
	}
	f := output.NewFormatter(stdout)
	if format == "table" || format == "both" {
		if err := f.RenderTerminalTable(report); err != nil {
			return fail(err)
		}
	}
	if format == "json" || format == "both" {
		if err := f.ExportJSON(report, path); err != nil {
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
	return reconcileCapture(ctx, cfg, &evidence.Data{})
}
func reconcileCapture(ctx context.Context, cfg *config.Config, captured *evidence.Data) (*types.DiscrepancyReport, error) {
	start, end := cfg.Reconciliation.Start, cfg.Reconciliation.End
	input, err := evidence.ReadBytes(cfg.TargetApp.SourcePath)
	if err != nil {
		return nil, err
	}
	captured.SourceHashes = map[string]string{"internal": evidence.Hash(input)}
	var ips []types.InternalPayment
	if cfg.TargetApp.Format == "sdp-csv" {
		ips, err = sdp.Parse(bytes.NewReader(input), cfg.Stellar.Network, cfg.TargetApp.Name, *cfg.TargetApp.SDP)
		captured.Assertions = map[string]string{"sdp_release": sdp.Release, "sdp_revision": sdp.Revision, "scope": cfg.TargetApp.SDP.Assertion}
	} else {
		if cfg.TargetApp.Format == "json" {
			var checked []types.InternalPayment
			if err = evidence.Decode(input, &checked); err != nil {
				return nil, err
			}
		}
		ips, err = file.NewFileConnector(cfg.TargetApp.SourcePath, cfg.TargetApp.Format, cfg.TargetApp.Name).Parse(ctx, bytes.NewReader(input), connector.Filter{TimeStart: start, TimeEnd: end})
	}
	if err != nil {
		return nil, err
	}
	// Free-form metadata is not needed by matching and can contain secrets/PII.
	for i := range ips {
		ips[i].Metadata = nil
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
		raw, err := evidence.ReadBytes(cfg.Stellar.OnChainPath)
		if err != nil {
			return nil, err
		}
		captured.SourceHashes["observations"] = evidence.Hash(raw)
		if err = evidence.Decode(raw, &observed); err != nil {
			return nil, err
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
		// Keep provider identity without URL credentials, query keys or private paths.
		u, _ := url.Parse(cfg.Stellar.HorizonURL)
		observed.Coverage.Source = "horizon:" + u.Host
	}
	scoped := []types.OnChainPayment{}
	for _, p := range observed.Payments {
		p.Memo = ""
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
	report := engine.NewReconciler(opts).Reconcile(cfg.TargetApp.Name, start, end, ips, scoped)
	captured.Tool = Version
	captured.Build = evidence.BuildInfo()
	captured.App = cfg.TargetApp.Name
	captured.Start = start
	captured.End = end
	captured.Options = opts
	captured.Internal = ips
	captured.Observations = scoped
	captured.Report = report
	return report, nil
}
