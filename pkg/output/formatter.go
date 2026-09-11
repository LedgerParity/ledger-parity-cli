package output

import (
	"encoding/json"
	"fmt"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"os"
	"strings"
)

type Formatter struct{ Out io.Writer }

func NewFormatter(out io.Writer) *Formatter {
	if out == nil {
		out = os.Stdout
	}
	return &Formatter{out}
}
func (f *Formatter) RenderTerminalTable(r *types.DiscrepancyReport) error {
	if _, err := fmt.Fprintln(f.Out, r.Summary()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f.Out, "Coverage: complete=%t; internal_complete=%t; source=%s\n%s\n", r.Coverage.Complete, r.Coverage.InternalComplete, r.Coverage.Source, r.Coverage.Reason); err != nil {
		return err
	}
	for _, v := range r.Results {
		if len(v.CandidateOperationIDs) > 1 {
			v.Notes += "; candidate operations: " + strings.Join(v.CandidateOperationIDs, ", ")
		}
		id := "-"
		if v.InternalPayment != nil {
			id = v.InternalPayment.ID
		}
		op := "-"
		if v.OnChainPayment != nil {
			op = v.OnChainPayment.OperationID
		}
		if _, err := fmt.Fprintf(f.Out, "%s | internal=%s | op=%s | %s | %s\n", v.Status, id, op, v.Discrepancy, v.Notes); err != nil {
			return err
		}
	}
	return nil
}
func (f *Formatter) ExportJSON(r *types.DiscrepancyReport, path string) error {
	if path == "-" {
		enc := json.NewEncoder(f.Out)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}
