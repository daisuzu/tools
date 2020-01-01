package source

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/telemetry/log"
	"golang.org/x/tools/internal/telemetry/trace"
)

func WorkspaceSymbols(ctx context.Context, view View, query string) ([]protocol.SymbolInformation, error) {
	ctx, done := trace.StartSpan(ctx, "source.WorkspaceSymbols")
	defer done()

	var symbols []protocol.SymbolInformation
	for _, ph := range view.Snapshot().KnownPackages(ctx) {
		pkg, err := ph.Check(ctx)
		if err != nil {
			return nil, err
		}
		for _, handle := range pkg.CompiledGoFiles() {
			file, mapper, _, err := handle.Cached()
			if err != nil {
				return nil, err
			}
			for _, si := range searchSymbols(file.Decls, pkg.GetTypesInfo(), query) {
				rng, err := nodeToProtocolRange(ctx, view, mapper, si.node)
				if err != nil {
					log.Error(ctx, "Error getting range for node", err)
					continue
				}
				symbols = append(symbols, protocol.SymbolInformation{
					Name: si.name,
					Kind: si.kind,
					Location: protocol.Location{
						URI:   protocol.NewURI(handle.File().Identity().URI),
						Range: rng,
					},
				})
			}
		}
	}
	return symbols, nil
}

type symbolInformation struct {
	name string
	kind protocol.SymbolKind
	node ast.Node
}

func searchSymbols(decls []ast.Decl, info *types.Info, query string) []symbolInformation {
	var result []symbolInformation
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if strings.Contains(decl.Name.Name, query) {
				kind := protocol.Function
				if decl.Recv != nil {
					kind = protocol.Method
				}
				result = append(result, symbolInformation{
					name: decl.Name.Name,
					kind: kind,
					node: decl.Name,
				})
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if strings.Contains(spec.Name.Name, query) {
						result = append(result, symbolInformation{
							name: spec.Name.Name,
							kind: typeToKind(info.TypeOf(spec.Type)),
							node: spec.Name,
						})
					}
					if st, ok := spec.Type.(*ast.StructType); ok {
						for _, field := range st.Fields.List {
							if len(field.Names) == 0 {
								name := types.ExprString(field.Type)
								if strings.Contains(name, query) {
									result = append(result, symbolInformation{
										name: name,
										kind: protocol.Field,
										node: field,
									})
								}
								continue
							}
							for _, name := range field.Names {
								if strings.Contains(name.Name, query) {
									result = append(result, symbolInformation{
										name: name.Name,
										kind: protocol.Field,
										node: name,
									})
								}
							}
						}
					}
					if it, ok := spec.Type.(*ast.InterfaceType); ok {
						for _, field := range it.Methods.List {
							if len(field.Names) == 0 {
								name := types.ExprString(field.Type)
								if strings.Contains(name, query) {
									result = append(result, symbolInformation{
										name: name,
										kind: protocol.Interface,
										node: field,
									})
								}
								continue
							}
							for _, name := range field.Names {
								if strings.Contains(name.Name, query) {
									result = append(result, symbolInformation{
										name: name.Name,
										kind: protocol.Method,
										node: name,
									})
								}
							}
						}
					}
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						if strings.Contains(name.Name, query) {
							kind := protocol.Variable
							if decl.Tok == token.CONST {
								kind = protocol.Constant
							}
							result = append(result, symbolInformation{
								name: name.Name,
								kind: kind,
								node: name,
							})
						}
					}
				}
			}
		}
	}
	return result
}

func typeToKind(typ types.Type) protocol.SymbolKind {
	switch typ := typ.Underlying().(type) {
	case *types.Interface:
		return protocol.Interface
	case *types.Struct:
		return protocol.Struct
	case *types.Signature:
		if typ.Recv() != nil {
			return protocol.Method
		}
		return protocol.Function
	case *types.Named:
		return typeToKind(typ.Underlying())
	case *types.Basic:
		i := typ.Info()
		switch {
		case i&types.IsNumeric != 0:
			return protocol.Number
		case i&types.IsBoolean != 0:
			return protocol.Boolean
		case i&types.IsString != 0:
			return protocol.String
		}
	}
	return protocol.Variable
}
