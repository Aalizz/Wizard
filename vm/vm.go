package vm

// 定义虚拟机

import (
	"fmt"
	"my.com/myfile/code"
	"my.com/myfile/compiler"
	"my.com/myfile/object"
)

// 定义全局变量True,False
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

const StackSize = 2048

type VM struct {
	constants    []object.Object   // 常量池
	instructions code.Instructions // 指令

	stack []object.Object // 栈
	sp    int             // 指向栈里下一个空闲槽，栈顶的值是stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM { // 创建栈
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) StackTop() object.Object { // 访问栈顶元素
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error { // 运行虚拟机
	for ip := 0; ip <= len(vm.instructions); ip++ { // 取出指令
		op := code.Opcode(vm.instructions[ip]) // 也可以使用code.Lookup代替，但是速度慢

		switch op {
		case code.OpConstant: // 从虚拟机的常量池中取出一个常量值，并将其压入到虚拟机的栈中
			constIndex := code.ReadUint16(vm.instructions[ip+1:]) // 这里也可以用code.ReadOperands代替ReadUint16，但是速度慢
			ip += 2                                               // opConstant宽度为2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv: // 加法
			//right := vm.pop()
			//left := vm.pop()
			//leftValue := left.(*object.Integer).Value
			//rightValue := right.(*object.Integer).Value
			//result := leftValue + rightValue
			//vm.push(&object.Integer{Value: result})
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error { // 压栈操作
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overFlow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object { // 出栈操作
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) LastPoppedStackElem() object.Object { // 返回最近弹栈的元素
	return vm.stack[vm.sp]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error { // 做类型断言
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operator: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error { // 处理整数操作
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unkoown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}
