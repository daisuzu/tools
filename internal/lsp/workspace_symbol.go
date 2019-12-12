package lsp

import (
	"context"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/lsp/source"
	"golang.org/x/tools/internal/telemetry/trace"
)

func (s *Server) workspaceSymbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "lsp.Server.workspaceSymbol")
	defer done()

	return source.WorkspaceSymbols(ctx, s.session.Views()[0], params.Query)
}
