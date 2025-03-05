package compiler

type Symbol struct {
	Name  string
	Index int
}

type SymbolTable struct {
	store map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func (table *SymbolTable) Define(identifier string) Symbol {
	if symbol, exists := table.Lookup(identifier); exists {
		return symbol
	}
	symbol := Symbol{Name: identifier, Index: len(table.store)}
	table.store[identifier] = symbol
	return symbol
}

func (table *SymbolTable) Lookup(identifier string) (Symbol, bool) {
	if symbol, ok := table.store[identifier]; ok {
		return symbol, true
	}
	return Symbol{}, false
}
