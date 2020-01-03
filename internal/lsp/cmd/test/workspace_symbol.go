package cmdtest

import (
	"testing"

	"golang.org/x/tools/internal/lsp/protocol"
)

func (r *runner) WorkspaceSymbols(t *testing.T, query string, expectedSymbols []protocol.SymbolInformation) {
	if query == "" {
		t.Skip("skipping empty query")
	}
	got, _ := r.NormalizeGoplsCmd(t, "workspace_symbols", query)
	expect := string(r.data.Golden("workspacesymbol", query, func() ([]byte, error) {
		return []byte(got), nil
	}))
	if expect != got {
		t.Errorf("workspace_symbols failed for %s expected:\n%s\ngot:\n%s", query, expect, got)
	}
}
