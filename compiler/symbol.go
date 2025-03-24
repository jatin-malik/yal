package compiler

type SymbolScope string

const (
	GLOBAL  SymbolScope = "global"
	LOCAL   SymbolScope = "local"
	BUILTIN SymbolScope = "builtin"
)

type Symbol struct {
	Name  string
	Index int
	Scope SymbolScope
}

type SymbolTable struct {
	store map[string]Symbol
	outer *SymbolTable
}

var builtInSymbols = map[string]Symbol{
	"len": {
		Name:  "len",
		Index: 0,
		Scope: BUILTIN,
	},
	"first": {
		Name:  "first",
		Index: 1,
		Scope: BUILTIN,
	},
	"last": {
		Name:  "last",
		Index: 2,
		Scope: BUILTIN,
	},
	"rest": {
		Name:  "rest",
		Index: 3,
		Scope: BUILTIN,
	},
	"push": {
		Name:  "push",
		Index: 4,
		Scope: BUILTIN,
	},
	"puts": {
		Name:  "puts",
		Index: 5,
		Scope: BUILTIN,
	},
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), outer: outer}
}

func (table *SymbolTable) Define(identifier string) Symbol {
	if symbol, exists := table.store[identifier]; exists {
		return symbol
	}
	symbol := Symbol{Name: identifier, Index: len(table.store)}
	if table.outer == nil {
		symbol.Scope = GLOBAL
	} else {
		symbol.Scope = LOCAL
	}

	table.store[identifier] = symbol
	return symbol
}

func (table *SymbolTable) Lookup(identifier string) (Symbol, bool) {
	if symbol, ok := table.store[identifier]; ok {
		return symbol, true
	}
	if table.outer != nil {
		return table.outer.Lookup(identifier)
	}

	if symbol, exists := builtInSymbols[identifier]; exists {
		return symbol, true
	}

	return Symbol{}, false
}

func (table *SymbolTable) len() int {
	return len(table.store)
}
