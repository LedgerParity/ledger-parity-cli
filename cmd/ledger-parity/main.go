package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/LedgerParity/ledger-parity-cli/pkg/config"
	"github.com/LedgerParity/ledger-parity-cli/pkg/output"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/file"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/stellopay"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/ingest"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

const Version = "1.0.0"

func main() {
	configPath := flag.String("config", "", "Path to YAML/JSON configuration file")
	outputFormat := flag.String("format", "table", "Output format: table, json, or both")
	outputPath := flag.String("out", "discrepancy_report.json", "Output file path for JSON report")
	demoMode := flag.Bool("demo", false, "Run standalone demonstration with seeded test vectors")
	showVersion := flag.Bool("version", false, "Show LedgerParity version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("LedgerParity CLI v%s — Stellar Payment Reconciliation Service\n", Version)
		os.Exit(0)
	}

	fmt.Println("🔍 LedgerParity Payment Reconciliation Engine v" + Version)

	if *demoMode || *configPath == "" {
		fmt.Println("ℹ️  Running in standalone demonstration mode (seeded discrepancy test vectors)...")
		runDemo(*outputFormat, *outputPath)
		return
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	runReconciliation(cfg, *outputFormat, *outputPath)
}

func runDemo(outputFormat, outputPath string) {
	now := time.Now()
	windowStart := now.Add(-24 * time.Hour)
	windowEnd := now

	// Seeded Internal Payments (Target App: Stellopay)
	internals := []types.InternalPayment{
		{
			ID:          "stell_101",
			SourceApp:   "stellopay",
			ReferenceID: "tx_hash_001",
			Sender:      "GEMPLOYER1111111111111111111111111111111111111111111",
			Recipient:   "GEMPLOYEE2222222222222222222222222222222222222222222",
			Amount:      "150.0000000",
			Asset:       "XLM",
			Timestamp:   now.Add(-10 * time.Minute),
			Status:      "COMPLETED",
		},
		{
			ID:          "stell_102",
			SourceApp:   "stellopay",
			ReferenceID: "tx_hash_002",
			Sender:      "GEMPLOYER1111111111111111111111111111111111111111111",
			Recipient:   "GEMPLOYEE3333333333333333333333333333333333333333333",
			Amount:      "300.0000000",
			Asset:       "USDC",
			Timestamp:   now.Add(-30 * time.Minute),
			Status:      "COMPLETED",
		},
		{
			ID:          "stell_103",
			SourceApp:   "stellopay",
			Sender:      "GEMPLOYER1111111111111111111111111111111111111111111",
			Recipient:   "GEMPLOYEE4444444444444444444444444444444444444444444",
			Amount:      "500.0000000",
			Asset:       "XLM",
			Timestamp:   now.Add(-5 * time.Minute),
			Status:      "PENDING_SETTLEMENT", // Missing on-chain settlement
		},
		{
			ID:          "stell_104_dup",
			SourceApp:   "stellopay",
			Sender:      "GEMPLOYER1111111111111111111111111111111111111111111",
			Recipient:   "GEMPLOYEE2222222222222222222222222222222222222222222",
			Amount:      "150.0000000",
			Asset:       "XLM",
			Timestamp:   now.Add(-10 * time.Minute),
			Status:      "COMPLETED", // Duplicate internal record
		},
	}

	// Seeded On-Chain Stellar Payments
	onChains := []types.OnChainPayment{
		{
			TransactionHash: "tx_hash_001",
			OperationID:     "op_9001",
			Account:         "GEMPLOYER1111111111111111111111111111111111111111111",
			Destination:     "GEMPLOYEE2222222222222222222222222222222222222222222",
			Amount:          "150.0000000",
			AssetCode:       "XLM",
			Timestamp:       now.Add(-10 * time.Minute),
			Successful:      true,
		},
		{
			TransactionHash: "tx_hash_002",
			OperationID:     "op_9002",
			Account:         "GEMPLOYER1111111111111111111111111111111111111111111",
			Destination:     "GEMPLOYEE3333333333333333333333333333333333333333333",
			Amount:          "250.0000000", // Amount mismatch (Internal expected 300, on-chain settled 250)
			AssetCode:       "USDC",
			Timestamp:       now.Add(-30 * time.Minute),
			Successful:      true,
		},
		{
			TransactionHash: "tx_hash_orphan_99",
			OperationID:     "op_9099",
			Account:         "GUNKNOWN9999999999999999999999999999999999999999999",
			Destination:     "GEMPLOYEE2222222222222222222222222222222222222222222",
			Amount:          "1000.0000000",
			AssetCode:       "XLM",
			Timestamp:       now.Add(-1 * time.Hour),
			Successful:      true, // Orphaned on-chain payment
		},
	}

	rec := engine.NewReconciler()
	report := rec.Reconcile("stellopay_demo", windowStart, windowEnd, internals, onChains)

	fmtter := output.NewFormatter(os.Stdout)
	if outputFormat == "table" || outputFormat == "both" {
		fmtter.RenderTerminalTable(report)
	}

	if outputFormat == "json" || outputFormat == "both" {
		if err := fmtter.ExportJSON(report, outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error exporting JSON: %v\n", err)
		}
	}
}

func runReconciliation(cfg *config.Config, outputFormat, outputPath string) {
	ctx := context.Background()

	var conn connector.Connector
	if cfg.TargetApp.Name == "stellopay" {
		conn = stellopay.NewStellopayConnector(cfg.TargetApp.SourcePath)
	} else {
		conn = file.NewFileConnector(cfg.TargetApp.SourcePath, cfg.TargetApp.Format, cfg.TargetApp.Name)
	}

	filter := connector.Filter{Limit: 1000}
	internals, err := conn.FetchInternalPayments(ctx, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error fetching internal payment records: %v\n", err)
		os.Exit(1)
	}

	horizonClient := ingest.NewHorizonIngestor(cfg.Stellar.HorizonURL)
	now := time.Now()
	windowStart := now.Add(-24 * time.Hour)
	onChains, err := horizonClient.FetchOnChainPayments(ctx, cfg.Stellar.Accounts, windowStart, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ Warning querying Horizon: %v (falling back to empty on-chain list)\n", err)
		onChains = []types.OnChainPayment{}
	}

	recOpts := engine.ReconcileOptions{
		TimeframeToleranceSec: cfg.Reconciliation.TimeframeToleranceSec,
		IgnoreFailedOnChain:   cfg.Reconciliation.IgnoreFailedOnChain,
	}

	rec := engine.NewReconciler(recOpts)
	report := rec.Reconcile(cfg.TargetApp.Name, windowStart, now, internals, onChains)

	fmtter := output.NewFormatter(os.Stdout)
	if outputFormat == "table" || outputFormat == "both" {
		fmtter.RenderTerminalTable(report)
	}

	if outputFormat == "json" || outputFormat == "both" {
		if err := fmtter.ExportJSON(report, outputPath); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error exporting JSON: %v\n", err)
		}
	}
}
