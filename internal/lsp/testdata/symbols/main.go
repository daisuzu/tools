package main

import (
	"io"
)

var x = 42 //@symbol("x", "x", "Variable", ""), workspacesymbol("x", "x", "Variable")

const y = 43 //@symbol("y", "y", "Constant", ""), workspacesymbol("y", "y", "Constant")

type Number int //@symbol("Number", "Number", "Number", ""), workspacesymbol("Number", "Number", "Number")

type Alias = string //@symbol("Alias", "Alias", "String", ""), workspacesymbol("Alias", "Alias", "String")

type NumberAlias = Number //@symbol("NumberAlias", "NumberAlias", "Number", ""), workspacesymbol("NumberAlias", "NumberAlias", "Number")

type (
	Boolean   bool   //@symbol("Boolean", "Boolean", "Boolean", ""), workspacesymbol("Boolean", "Boolean", "Boolean")
	BoolAlias = bool //@symbol("BoolAlias", "BoolAlias", "Boolean", ""), workspacesymbol("BoolAlias", "BoolAlias", "Boolean")
)

type Foo struct { //@symbol("Foo", "Foo", "Struct", ""), workspacesymbol("Foo", "Foo", "Struct")
	Quux           //@symbol("Quux", "Quux", "Field", "Foo"), workspacesymbol("Quux", "Quux", "Field")
	W    io.Writer //@symbol("W" , "W", "Field", "Foo"), workspacesymbol("W" , "W", "Field")
	Bar  int       //@symbol("Bar", "Bar", "Field", "Foo"), workspacesymbol("Bar", "Bar", "Field")
	baz  string    //@symbol("baz", "baz", "Field", "Foo"), workspacesymbol("baz", "baz", "Field")
}

type Quux struct { //@symbol("Quux", "Quux", "Struct", ""), workspacesymbol("Quux", "Quux", "Struct")
	X, Y float64 //@symbol("X", "X", "Field", "Quux"), symbol("Y", "Y", "Field", "Quux"), workspacesymbol("X", "X", "Field"), workspacesymbol("Y", "Y", "Field")
}

func (f Foo) Baz() string { //@symbol("Baz", "Baz", "Method", "Foo"), workspacesymbol("Baz", "Baz", "Method")
	return f.baz
}

func (q *Quux) Do() {} //@symbol("Do", "Do", "Method", "Quux"), workspacesymbol("Do", "Do", "Method")

func main() { //@symbol("main", "main", "Function", ""), workspacesymbol("main", "main", "Function")

}

type Stringer interface { //@symbol("Stringer", "Stringer", "Interface", ""), workspacesymbol("Stringer", "Stringer", "Interface")
	String() string //@symbol("String", "String", "Method", "Stringer"), workspacesymbol("String", "String", "Method")
}

type ABer interface { //@symbol("ABer", "ABer", "Interface", ""), workspacesymbol("ABer", "ABer", "Interface")
	B()        //@symbol("B", "B", "Method", "ABer"), workspacesymbol("B", "B", "Method")
	A() string //@symbol("A", "A", "Method", "ABer"), workspacesymbol("A", "A", "Method")
}

type WithEmbeddeds interface { //@symbol("WithEmbeddeds", "WithEmbeddeds", "Interface", ""), workspacesymbol("WithEmbeddeds", "WithEmbeddeds", "Interface")
	Do()      //@symbol("Do", "Do", "Method", "WithEmbeddeds"), workspacesymbol("Do", "Do", "Method")
	ABer      //@symbol("ABer", "ABer", "Interface", "WithEmbeddeds"), workspacesymbol("ABer", "ABer", "Interface")
	io.Writer //@symbol("io.Writer", "io.Writer", "Interface", "WithEmbeddeds"), workspacesymbol("io.Writer", "io.Writer", "Interface")
}
