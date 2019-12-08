package lsp

import (
	"context"
	"fmt"
	"go/ast"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/telemetry/trace"
)

func (s *Server) workspaceSymbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "lsp.Server.workspaceSymbol")
	defer done()

	viewNames := make([]string, len(s.session.Views()))

	symbolsMap := map[string]protocol.SymbolInformation{}
	for i, view := range s.session.Views() {
		viewNames[i] = view.Folder().Filename()
		for _, pkg := range view.Snapshot().KnownPackages(ctx) {
			for _, handle := range pkg.CompiledGoFiles() {
				file, mapper, err1, err2 := handle.Parse(ctx)
				if err1 != nil || err2 != nil {
					continue
				}
				ast.Inspect(file, func(node ast.Node) bool {
					switch t := node.(type) {
					case *ast.Ident:
						if strings.Contains(t.Name, params.Query) {
							pos := view.Session().Cache().FileSet().Position(t.Pos())
							span, err := mapper.PointSpan(protocol.Position{Line: float64(pos.Line - 1), Character: float64(pos.Column)})
							if err != nil {
								return false
							}
							rng, err := mapper.Range(span)
							if err != nil {
								return false
							}
							location := protocol.Location{
								URI:   protocol.DocumentURI("file://" + handle.File().Identity().URI.Filename()),
								Range: rng,
							}
							symbolsMap[fmt.Sprintf("%v", location)] = protocol.SymbolInformation{
								Name:     t.Name,
								Kind:     protocol.File,
								Location: location,
							}
						}
						return false
					}
					return true
				})
			}
		}
	}

	symbols := make([]protocol.SymbolInformation, 0, len(symbolsMap))
	for _, v := range symbolsMap {
		symbols = append(symbols, v)
	}
	sort.Slice(symbols, func(i, j int) bool {
		li := strings.Index(symbols[i].Name, params.Query) + len(symbols[i].Name)
		lj := strings.Index(symbols[j].Name, params.Query) + len(symbols[j].Name)
		if li == lj {
			di := pkgDistance(symbols[i].Location.URI, viewNames)
			dj := pkgDistance(symbols[j].Location.URI, viewNames)
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

func pkgDistance(uri string, names []string) int {
	if len(names) == 0 {
		return -1
	}

	path := filepath.Dir(strings.ReplaceAll(uri, "file://", ""))
	if path == names[0] {
		return 0
	}

	dir := names[0]
	for i := 1; i <= 3; i++ {
		if strings.HasPrefix(path, dir) {
			return i
		}
		dir = filepath.Dir(dir)
	}

	return 100
}
