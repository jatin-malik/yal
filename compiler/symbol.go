package compiler

type SymbolScope string

const (
	GLOBAL   SymbolScope = "global"
	LOCAL    SymbolScope = "local"
	BUILTIN  SymbolScope = "builtin"
	FREE     SymbolScope = "free"
	FUNCTION SymbolScope = "function"
)

type Symbol struct {
	Name  string
	Index int
	Scope SymbolScope
}

type SymbolTable struct {
	store       map[string]Symbol
	outer       *SymbolTable
	freeSymbols []Symbol
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

func (table *SymbolTable) DefineFunctionSymbol(name string) Symbol {
	symbol := Symbol{Name: name, Index: len(table.store), Scope: FUNCTION}
	table.store[name] = symbol
	return symbol
}

func (table *SymbolTable) defineFree(symbol Symbol) Symbol {
	table.freeSymbols = append(table.freeSymbols, symbol)
	freeSymbol := Symbol{Name: symbol.Name, Index: len(table.freeSymbols) - 1, Scope: FREE}
	table.store[freeSymbol.Name] = freeSymbol
	return freeSymbol
}

func (table *SymbolTable) Lookup(identifier string) (Symbol, bool) {
	if symbol, ok := table.store[identifier]; ok {
		return symbol, true
	}
	if table.outer != nil {
		symbol, ok := table.outer.Lookup(identifier)
		if symbol.Scope == LOCAL || symbol.Scope == FREE {
			// this is a free variable
			return table.defineFree(symbol), ok
		}
		return symbol, ok
	}

	if symbol, exists := builtInSymbols[identifier]; exists {
		return symbol, true
	}

	return Symbol{}, false
}

func (table *SymbolTable) len() int {
	return len(table.store)
}
