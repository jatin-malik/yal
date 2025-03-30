package compiler

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/object"
)

type CompilationScope struct {
	instructions       bytecode.Instructions
	lastAddedInsOffset int
}

func NewCompilationScope() *CompilationScope {
	return &CompilationScope{
		instructions: bytecode.Instructions{},
	}
}

type Compiler struct {
	scopes         []*CompilationScope
	activeScopeIdx int
	constantPool   []object.Object
	symbolTable    *SymbolTable
}

// ByteCode encloses the output of the compiler
type ByteCode struct {
	Instructions bytecode.Instructions
	ConstantPool []object.Object
}

type Option func(*Compiler)

// WithSymbolTable allows setting a custom symbol table.
func WithSymbolTable(symTable *SymbolTable) Option {
	return func(c *Compiler) {
		c.symbolTable = symTable
	}
}

func WithConstantPool(constantPool []object.Object) Option {
	return func(c *Compiler) {
		c.constantPool = constantPool
	}
}

func New(options ...Option) *Compiler {
	var scopes []*CompilationScope
	scopes = append(scopes, NewCompilationScope())
	compiler := &Compiler{
		scopes:       scopes,
		constantPool: []object.Object{},
		symbolTable:  NewSymbolTable(nil),
	}

	// Apply provided options
	for _, option := range options {
		option(compiler)
	}

	return compiler
}

// Compile walks through the input AST and generates bytecode. It also populates the constant pool as it evaluates
// constant literals in the AST. It returns an error in case compilation fails.
func (compiler *Compiler) Compile(node ast.Node) error {
	activeScope := compiler.scopes[compiler.activeScopeIdx]
	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			err := compiler.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			err := compiler.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		var symbol Symbol
		if fl, ok := n.Right.(*ast.FunctionLiteral); ok {
			// Register function name first to allow recursive functions
			symbol = compiler.symbolTable.Define(n.Name.Value)

			fl.Name = n.Name.Value // assign function literal its name

			err := compiler.Compile(fl)
			if err != nil {
				return err
			}
		} else {
			err := compiler.Compile(n.Right)
			if err != nil {
				return err
			}
			symbol = compiler.symbolTable.Define(n.Name.Value)
		}

		compiler.storeSymbol(symbol)

	case *ast.ReturnStatement:
		err := compiler.Compile(n.Value)
		if err != nil {
			return err
		}
		compiler.emit(bytecode.OpReturnValue)
	case *ast.ExpressionStatement:
		err := compiler.Compile(n.Expr)
		if err != nil {
			return err
		}
	case *ast.IfElseConditional:
		err := compiler.Compile(n.Condition)
		if err != nil {
			return err
		}

		compiler.emit(bytecode.OpJumpIfFalse, 9999)

		conditionalJumpOffset := activeScope.lastAddedInsOffset

		err = compiler.Compile(n.Consequence)
		if err != nil {
			return err
		}

		compiler.emit(bytecode.OpJump, 9999)
		jumpOffset := activeScope.lastAddedInsOffset

		// Back-patch conditional jump
		newConditionalJumpIns, _ := bytecode.Make(bytecode.OpJumpIfFalse, len(activeScope.instructions))
		compiler.modifyInstruction(conditionalJumpOffset, newConditionalJumpIns)

		if n.Alternative != nil {
			err = compiler.Compile(n.Alternative)
			if err != nil {
				return err
			}
		} else {
			compiler.emit(bytecode.OpPushNull)
		}

		newJumpIns, _ := bytecode.Make(bytecode.OpJump, len(activeScope.instructions))
		compiler.modifyInstruction(jumpOffset, newJumpIns)

	case *ast.ArrayLiteral:
		for _, element := range n.Elements {
			err := compiler.Compile(element)
			if err != nil {
				return err
			}
		}

		compiler.emit(bytecode.OpArray, len(n.Elements))
	case *ast.HashLiteral:
		for k, v := range n.Pairs {
			err := compiler.Compile(v)
			if err != nil {
				return err
			}
			err = compiler.Compile(k)
			if err != nil {
				return err
			}
		}

		compiler.emit(bytecode.OpHash, len(n.Pairs))
	case *ast.IndexExpression:
		err := compiler.Compile(n.Left)
		if err != nil {
			return err
		}
		err = compiler.Compile(n.Index)
		if err != nil {
			return err
		}

		compiler.emit(bytecode.OpIndex)
	case *ast.FunctionLiteral:
		compiler.enterScope()
		activeScope := compiler.scopes[compiler.activeScopeIdx]

		localSymbolTable := NewSymbolTable(compiler.symbolTable)

		compiler.symbolTable = localSymbolTable

		for _, param := range n.Parameters {
			compiler.symbolTable.Define(param.Value)
		}

		if n.Name != "" {
			localSymbolTable.DefineFunctionSymbol(n.Name)
		}

		err := compiler.Compile(n.Body)
		if err != nil {
			return err
		}

		if len(activeScope.instructions) == 0 {
			// empty function body
			compiler.emit(bytecode.OpPushNull)
			compiler.emit(bytecode.OpReturnValue)
		}

		if bytecode.OpCode(activeScope.instructions[activeScope.lastAddedInsOffset]) != bytecode.OpReturnValue {
			// implicit return
			compiler.emit(bytecode.OpReturnValue)
		}
		compiledInstructions := activeScope.instructions

		compiler.symbolTable = localSymbolTable.outer
		compiler.exitScope()

		compiledFunctionObj := &object.CompiledFunction{
			Instructions: compiledInstructions,
			NumLocals:    localSymbolTable.len(),
		}
		idx := compiler.addConstant(compiledFunctionObj)

		// Push free variables on stack
		for _, symbol := range localSymbolTable.freeSymbols {
			compiler.loadSymbol(symbol)
		}

		compiler.emit(bytecode.OpClosure, idx, len(localSymbolTable.freeSymbols))

	case *ast.CallExpression:
		err := compiler.Compile(n.Function)
		if err != nil {
			return err
		}

		for _, arg := range n.Arguments {
			err := compiler.Compile(arg)
			if err != nil {
				return err
			}
		}

		compiler.emit(bytecode.OpCall, len(n.Arguments))
	case *ast.PrefixExpression:
		err := compiler.Compile(n.Right)
		if err != nil {
			return err
		}

		switch n.Operator {
		case "!":
			compiler.emit(bytecode.OpNegateBoolean)
		case "-":
			compiler.emit(bytecode.OpNegateNumber)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *ast.InfixExpression:

		if n.Operator == "<" {
			err := compiler.Compile(n.Right)
			if err != nil {
				return err
			}

			err = compiler.Compile(n.Left)
			if err != nil {
				return err
			}

			compiler.emit(bytecode.OpGT)

		} else {
			err := compiler.Compile(n.Left)
			if err != nil {
				return err
			}

			err = compiler.Compile(n.Right)
			if err != nil {
				return err
			}

			switch n.Operator {
			case "+":
				compiler.emit(bytecode.OpAdd)
			case "-":
				compiler.emit(bytecode.OpSub)
			case "*":
				compiler.emit(bytecode.OpMul)
			case "/":
				compiler.emit(bytecode.OpDiv)
			case "==":
				compiler.emit(bytecode.OpEqual)
			case "!=":
				compiler.emit(bytecode.OpNotEqual)
			case ">":
				compiler.emit(bytecode.OpGT)
			default:
				return fmt.Errorf("unknown operator %s", n.Operator)
			}
		}
	case *ast.Identifier:
		symbol, exists := compiler.symbolTable.Lookup(n.Value)
		if !exists {
			return fmt.Errorf("unknown identifier %s", n.Value)
		}
		compiler.loadSymbol(symbol)

	case *ast.IntegerLiteral:
		obj := &object.Integer{Value: n.Value}
		idx := compiler.addConstant(obj)
		compiler.emit(bytecode.OpPush, idx)
	case *ast.StringLiteral:
		obj := &object.String{Value: n.Value}
		idx := compiler.addConstant(obj)
		compiler.emit(bytecode.OpPush, idx)
	case *ast.BooleanLiteral:
		if n.Value {
			compiler.emit(bytecode.OpPushTrue)
		} else {
			compiler.emit(bytecode.OpPushFalse)
		}
	}
	return nil
}

// addConstant adds the constant to the constant pool and returns the index where it is stored
func (compiler *Compiler) addConstant(obj object.Object) int {
	compiler.constantPool = append(compiler.constantPool, obj)
	return len(compiler.constantPool) - 1
}

// addInstruction appends input instruction to the compiler instructions and returns the insert offset.
func (compiler *Compiler) addInstruction(ins []byte) {
	activeScope := compiler.scopes[compiler.activeScopeIdx]
	insertPos := len(activeScope.instructions)
	activeScope.instructions = append(activeScope.instructions, ins...)
	activeScope.lastAddedInsOffset = insertPos
}

func (compiler *Compiler) modifyInstruction(offset int, newInstruction []byte) {
	copy(compiler.scopes[compiler.activeScopeIdx].instructions[offset:], newInstruction)
}

// Output wraps compiler output in ByteCode struct and returns it
func (compiler *Compiler) Output() ByteCode {
	return ByteCode{
		Instructions: compiler.scopes[compiler.activeScopeIdx].instructions,
		ConstantPool: compiler.constantPool,
	}
}

func (compiler *Compiler) emit(op bytecode.OpCode, operands ...int) error {
	ins, err := bytecode.Make(op, operands...)
	if err != nil {
		return err
	}
	compiler.addInstruction(ins)
	return nil
}

func (compiler *Compiler) enterScope() {
	scope := NewCompilationScope()
	compiler.scopes = append(compiler.scopes, scope)
	compiler.activeScopeIdx++
}

func (compiler *Compiler) exitScope() {
	compiler.scopes = compiler.scopes[:compiler.activeScopeIdx]
	compiler.activeScopeIdx--
}

func (compiler *Compiler) loadSymbol(symbol Symbol) {
	switch symbol.Scope {
	case GLOBAL:
		compiler.emit(bytecode.OpGetGlobal, symbol.Index)
	case LOCAL:
		compiler.emit(bytecode.OpGetLocal, symbol.Index)
	case BUILTIN:
		compiler.emit(bytecode.OpGetBuiltIn, symbol.Index)
	case FREE:
		compiler.emit(bytecode.OpGetFree, symbol.Index)
	case FUNCTION:
		compiler.emit(bytecode.OpGetCurrentClosure)
	}
}

func (compiler *Compiler) storeSymbol(symbol Symbol) {
	if symbol.Scope == GLOBAL {
		compiler.emit(bytecode.OpSetGlobal, symbol.Index)
	} else {
		compiler.emit(bytecode.OpSetLocal, symbol.Index)
	}
}
