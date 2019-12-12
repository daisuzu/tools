package source

import (
	"context"
	"go/ast"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/telemetry/trace"
)

func WorkspaceSymbols(ctx context.Context, view View, query string) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "source.WorkspaceSymbols")
	defer done()

	pkgsMap := make(map[string]struct{})
	var symbols []protocol.SymbolInformation
	for _, pkg := range view.Snapshot().KnownPackages(ctx) {
		if _, ok := pkgsMap[pkg.PkgPath()]; ok {
			continue
		}
		pkgsMap[pkg.PkgPath()] = struct{}{}
		for _, handle := range pkg.CompiledGoFiles() {
			file, mapper, _, err := handle.Cached()
			if err != nil {
				continue
			}
			ast.Inspect(file, func(node ast.Node) bool {
				switch t := node.(type) {
				case *ast.Ident:
					if strings.Contains(t.Name, query) {
						rng, err := nodeToProtocolRange(ctx, view, mapper, t)
						if err != nil {
							return false
						}
						symbols = append(symbols, protocol.SymbolInformation{
							Name: t.Name,
							Kind: protocol.File,
							Location: protocol.Location{
								URI:   protocol.NewURI(handle.File().Identity().URI),
								Range: rng,
							},
						})
					}
					return false
				}
				return true
			})
		}
	}

	if len(symbols) == 0 {
		return nil, nil
	}

	sort.Slice(symbols, func(i, j int) bool {
		li := strings.Index(symbols[i].Name, query) + len(symbols[i].Name)
		lj := strings.Index(symbols[j].Name, query) + len(symbols[j].Name)
		if li == lj {
			dir := view.Folder().Filename()
			di := pkgDistance(symbols[i].Location.URI, dir)
			dj := pkgDistance(symbols[j].Location.URI, dir)
			if di == dj {
				if symbols[i].Location.URI == symbols[j].Location.URI {
					return symbols[i].Location.Range.Start.Line < symbols[j].Location.Range.Start.Line
				}
				return symbols[i].Location.URI < symbols[j].Location.URI
			}
			return di < dj
		}
		return li < lj
	})
	return symbols, nil
}

func pkgDistance(uri, dir string) int {
	path := filepath.Dir(strings.ReplaceAll(uri, "file://", ""))
	if path == dir {
		return 0
	}

	for i := 1; i <= 3; i++ {
		if strings.HasPrefix(path, dir) {
			return i
		}
		dir = filepath.Dir(dir)
	}

	return 100
}
