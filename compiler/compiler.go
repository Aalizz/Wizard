package compiler

import (
	"fmt"
	"my.com/myfile/ast"
	"my.com/myfile/code"
	"my.com/myfile/object"
)

type Compiler struct {
	instructions code.Instructions // 指令
	constants    []object.Object   // 常量池
}

func New() *Compiler { // 返回Compiler结构指针
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error { // 编译器
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement: // 表达式
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop) // 每次执行表达式后执行一次弹栈操作清理栈
	case *ast.InfixExpression: // 中缀表达式
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral: // 整数字面量
		integer := &object.Integer{Value: node.Value}   // 求值
		c.emit(code.OpConstant, c.addConstant(integer)) // 生成opConstant指令
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}
	return nil
}

type Bytecode struct {
	Instructions code.Instructions // 字节码
	Constants    []object.Object   // 切片类型，常量池
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int { // 将求值结果添加到常量池中
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1 // 添加到编译器constants切片末尾，返回其在constants切片中的索引来为其提供标识符，用作opConstant指令的操作数
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	// 生成指令
	ins := code.Make(op, operands...)

	// 将指令添加到指令集中
	pos := c.addInstruction(ins)

	// 返回新指令的位置
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	// 获取新指令的起始位置
	posNewInstruction := len(c.instructions)

	// 将新指令添加到指令集中
	c.instructions = append(c.instructions, ins...)

	// 返回新指令的起始位置
	return posNewInstruction
}
