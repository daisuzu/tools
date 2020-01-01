package tests

import (
	"path"

	"golang.org/x/tools/internal/lsp/protocol"
)

func FilterWorkspaceSymbols(got, want []protocol.SymbolInformation) []protocol.SymbolInformation {
	dirs := make(map[string]struct{})
	for _, sym := range want {
		dirs[path.Dir(sym.Location.URI)] = struct{}{}
	}

	var result []protocol.SymbolInformation
	for _, sym := range got {
		if _, ok := dirs[path.Dir(sym.Location.URI)]; ok {
			result = append(result, sym)
		}
	}
	return result
}
