package lsp

import (
	"context"
	"go/ast"
	"strings"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/telemetry/trace"
)

func (s *Server) workspaceSymbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "lsp.Server.workspaceSymbol")
	defer done()

	var symbols []protocol.SymbolInformation
	for _, view := range s.session.Views() {
		for _, pkg := range view.Snapshot().KnownPackages(ctx) {
			for _, h := range pkg.CompiledGoFiles() {
				file, mapper, err1, err2 := h.Parse(ctx)
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
							r, err := mapper.Range(span)
							if err != nil {
								return false
							}
							symbols = append(symbols, protocol.SymbolInformation{
								Name: t.Name,
								Kind: protocol.File,
								Location: protocol.Location{
									URI:   protocol.DocumentURI("file://" + h.File().Identity().URI.Filename()),
									Range: r,
								},
							})
						}
						return false
					}
					return true
				})
			}
		}
	}

	return symbols, nil
}
