package compiler

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/object"
)

type Compiler struct {
	instructions bytecode.Instructions
	constantPool []object.Object
}

// ByteCode encloses the output of the compiler
type ByteCode struct {
	Instructions bytecode.Instructions
	ConstantPool []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constantPool: []object.Object{},
	}
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
	case *ast.ExpressionStatement:
		err := compiler.Compile(n.Expr)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
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
			ins, err := bytecode.Make(bytecode.OpAdd)
			if err != nil {
				return err
			}
			compiler.addInstruction(ins)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *ast.IntegerLiteral:
		obj := &object.Integer{Value: n.Value}
		idx := compiler.addConstant(obj)
		ins, err := bytecode.Make(bytecode.OpPush, idx)
		if err != nil {
			return err
		}
		compiler.addInstruction(ins)
	}
	return nil
}

// addConstant adds the constant to the constant pool and returns the index where it is stored
func (compiler *Compiler) addConstant(obj object.Object) int {
	compiler.constantPool = append(compiler.constantPool, obj)
	return len(compiler.constantPool) - 1
}

// addInstruction adds input instruction to the compiler instructions
func (compiler *Compiler) addInstruction(ins []byte) {
	compiler.instructions = append(compiler.instructions, ins...)
}

// Emit wraps compiler output in ByteCode struct and returns it
func (compiler *Compiler) Emit() ByteCode {
	return ByteCode{
		Instructions: compiler.instructions,
		ConstantPool: compiler.constantPool,
	}
}
