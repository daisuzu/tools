// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"flag"
	"fmt"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/span"
	"golang.org/x/tools/internal/tool"
)

// workspaceSymbols implements the references verb for gopls
type workspaceSymbols struct {
	app *Application
}

func (r *workspaceSymbols) Name() string      { return "workspace_symbols" }
func (r *workspaceSymbols) Usage() string     { return "<query>" }
func (r *workspaceSymbols) ShortHelp() string { return "display selected package's symbols" }
func (r *workspaceSymbols) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), `
Example:
  $ gopls workspace_symbols query
`)
	f.PrintDefaults()
}
func (r *workspaceSymbols) Run(ctx context.Context, args ...string) error {
	if len(args) != 1 {
		return tool.CommandLineErrorf("workspace_symbols expects 1 argument (query)")
	}

	conn, err := r.app.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.terminate(ctx)

	p := protocol.WorkspaceSymbolParams{
		Query: args[0],
	}

	symbols, err := conn.Symbol(ctx, &p)
	if err != nil {
		return err
	}
	for _, s := range symbols {
		fmt.Println(symbolInfoToString(s))
	}

	return nil
}

func symbolInfoToString(symbol protocol.SymbolInformation) string {
	r := symbol.Location.Range
	// convert ranges to user friendly 1-based positions
	position := fmt.Sprintf("%v:%v-%v:%v",
		r.Start.Line+1,
		r.Start.Character+1,
		r.End.Line+1,
		r.End.Character+1,
	)

	return fmt.Sprintf("%s %s %s:%s", symbol.Name, symbol.Kind, span.NewURI(symbol.Location.URI).Filename(), position)
}
