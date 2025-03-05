package compiler

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/object"
)

type Compiler struct {
	instructions       bytecode.Instructions
	constantPool       []object.Object
	lastAddedInsOffset int
	symbolTable        *SymbolTable
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

func New(options ...Option) *Compiler {
	compiler := &Compiler{
		instructions: bytecode.Instructions{},
		constantPool: []object.Object{},
		symbolTable:  NewSymbolTable(),
	}

	// Apply provided options
	for _, option := range options {
		option(compiler)
	}

	return compiler
}

// Compile walks through the input AST and generates bytecode. It also populates the constant pool as it evaluates self
// evaluating literals in the AST. It returns an error in case compilation fails.
func (compiler *Compiler) Compile(node ast.Node) error {
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
		err := compiler.Compile(n.Right)
		if err != nil {
			return err
		}

		symbol := compiler.symbolTable.Define(n.Name.Value)
		compiler.emit(bytecode.OpSetGlobal, symbol.Index)
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

		conditionalJumpOffset := compiler.lastAddedInsOffset

		err = compiler.Compile(n.Consequence)
		if err != nil {
			return err
		}

		compiler.emit(bytecode.OpJump, 9999)
		jumpOffset := compiler.lastAddedInsOffset

		// Back-patch conditional jump
		newConditionalJumpIns, _ := bytecode.Make(bytecode.OpJumpIfFalse, len(compiler.instructions))
		compiler.modifyInstruction(conditionalJumpOffset, newConditionalJumpIns)

		if n.Alternative != nil {
			err = compiler.Compile(n.Alternative)
			if err != nil {
				return err
			}
		} else {
			compiler.emit(bytecode.OpPushNull)
		}

		newJumpIns, _ := bytecode.Make(bytecode.OpJump, len(compiler.instructions))
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
		compiler.emit(bytecode.OpGetGlobal, symbol.Index)
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
	insertPos := len(compiler.instructions)
	compiler.instructions = append(compiler.instructions, ins...)
	compiler.lastAddedInsOffset = insertPos
}

func (compiler *Compiler) modifyInstruction(offset int, newInstruction []byte) {
	copy(compiler.instructions[offset:], newInstruction)
}

// Output wraps compiler output in ByteCode struct and returns it
func (compiler *Compiler) Output() ByteCode {
	return ByteCode{
		Instructions: compiler.instructions,
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
