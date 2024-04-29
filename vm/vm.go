package vm

import "errors"

// StackVM 表示堆栈式虚拟机
type StackVM struct {
	stack []int
}

// Push 将值压入堆栈
func (vm *StackVM) Push(value int) {
	vm.stack = append(vm.stack, value)
}

// Pop 弹出堆栈顶部的值
func (vm *StackVM) Pop() (int, error) {
	if len(vm.stack) == 0 {
		return 0, errors.New("stack underflow")
	}
	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value, nil
}

// Add 从堆栈中弹出两个值，相加后将结果压入堆栈
func (vm *StackVM) Add() error {
	if len(vm.stack) < 2 {
		return errors.New("not enough operands on the stack for ADD operation")
	}
	a, _ := vm.Pop()
	b, _ := vm.Pop()
	vm.Push(a + b)
	return nil
}

// Sub 从堆栈中弹出两个值，第二个值减去第一个值，将结果压入堆栈
func (vm *StackVM) Sub() error {
	if len(vm.stack) < 2 {
		return errors.New("not enough operands on the stack for SUB operation")
	}
	a, _ := vm.Pop()
	b, _ := vm.Pop()
	vm.Push(b - a)
	return nil
}

// Mul 从堆栈中弹出两个值，相乘后将结果压入堆栈
func (vm *StackVM) Mul() error {
	if len(vm.stack) < 2 {
		return errors.New("not enough operands on the stack for MUL operation")
	}
	a, _ := vm.Pop()
	b, _ := vm.Pop()
	vm.Push(a * b)
	return nil
}

// Div 从堆栈中弹出两个值，第二个值除以第一个值，将结果压入堆栈
func (vm *StackVM) Div() error {
	if len(vm.stack) < 2 {
		return errors.New("not enough operands on the stack for DIV operation")
	}
	a, _ := vm.Pop()
	b, _ := vm.Pop()
	if a == 0 {
		return errors.New("division by zero")
	}
	vm.Push(b / a)
	return nil
}
