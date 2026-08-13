package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// Formatter renders discrepancy reports to terminal tables or JSON files.
type Formatter struct {
	Out io.Writer
}

func NewFormatter(out io.Writer) *Formatter {
	if out == nil {
		out = os.Stdout
	}
	return &Formatter{Out: out}
}

// RenderTerminalTable formats and prints a high-legibility terminal report.
func (f *Formatter) RenderTerminalTable(report *types.DiscrepancyReport) {
	banner := strings.Repeat("═", 84)
	subBanner := strings.Repeat("─", 84)

	fmt.Fprintln(f.Out, banner)
	fmt.Fprintf(f.Out, "  LEDGERPARITY RECONCILIATION REPORT — %s\n", strings.ToUpper(report.TargetApp))
	fmt.Fprintln(f.Out, banner)
	fmt.Fprintf(f.Out, "  Generated At:   %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(f.Out, "  Time Window:    %s to %s\n",
		report.TimeWindowStart.Format("2006-01-02 15:04"),
		report.TimeWindowEnd.Format("2006-01-02 15:04"),
	)
	fmt.Fprintln(f.Out, subBanner)
	fmt.Fprintf(f.Out, "  TOTAL INTERNAL: %-6d | TOTAL ON-CHAIN: %-6d | MATCHED: %-6d | DISCREPANCIES: %-6d\n",
		report.TotalInternal, report.TotalOnChain, report.TotalMatched, report.TotalDiscrepancies,
	)
	fmt.Fprintln(f.Out, subBanner)

	if len(report.DiscrepancyCounts) > 0 {
		fmt.Fprintln(f.Out, "  DISCREPANCY BREAKDOWN:")
		for discType, count := range report.DiscrepancyCounts {
			fmt.Fprintf(f.Out, "    • %-22s : %d\n", discType, count)
		}
		fmt.Fprintln(f.Out, subBanner)
	}

	fmt.Fprintln(f.Out, "  RECONCILIATION DETAILS:")
	fmt.Fprintf(f.Out, "  %-12s | %-12s | %-16s | %-12s | %-20s\n",
		"STATUS", "INTERNAL ID", "RECIPIENT", "AMOUNT", "DISCREPANCY / NOTES",
	)
	fmt.Fprintln(f.Out, subBanner)

	for _, res := range report.Results {
		statusStr := fmt.Sprintf("[%s]", res.Status)
		intID := "-"
		recipient := "-"
		amount := "-"

		if res.InternalPayment != nil {
			intID = truncateStr(res.InternalPayment.ID, 12)
			recipient = truncateStr(res.InternalPayment.Recipient, 16)
			amount = res.InternalPayment.Amount + " " + res.InternalPayment.Asset
		} else if res.OnChainPayment != nil {
			recipient = truncateStr(res.OnChainPayment.Destination, 16)
			amount = res.OnChainPayment.Amount + " " + res.OnChainPayment.AssetCode
		}

		noteStr := string(res.Discrepancy)
		if noteStr == "NONE" {
			noteStr = res.Notes
		}

		fmt.Fprintf(f.Out, "  %-12s | %-12s | %-16s | %-12s | %s\n",
			statusStr, intID, recipient, amount, noteStr,
		)
	}

	fmt.Fprintln(f.Out, banner)
}

// ExportJSON writes the full DiscrepancyReport as pretty-printed JSON to a file path.
func (f *Formatter) ExportJSON(report *types.DiscrepancyReport, filePath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshaling json report: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed writing json report to %s: %w", filePath, err)
	}

	fmt.Fprintf(f.Out, "✓ Exported JSON discrepancy report to %s\n", filePath)
	return nil
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
