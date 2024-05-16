package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte // 字节集合

type Opcode byte // 操作码

const (
	OpConstant Opcode = iota // 以操作数为索引检索常量并压栈
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
)

type Definition struct {
	Name          string // 操作数名车给
	OperandWidths []int  // 字节宽度
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // OpConstant占两字节，操作码为2
	OpAdd:      {"opAdd", []int{}},       // 空的整数切片，不需要操作数
	OpPop:      {"OpPop", []int{}},       // 弹栈，用于清理栈，每个表达式语句执行后都要执行这个操作码
	OpSub:      {"OpSub", []int{}},       // 减法操作
	OpMul:      {"OpMul", []int{}},       // 乘法
	OpDiv:      {"OpDiv", []int{}},       // 除法
	OpTrue:     {"OpTrue", []int{}},      // 布尔真
	OpFalse:    {"OpFalse", []int{}},     // 布尔假
}

func Lookup(op byte) (*Definition, error) { // 查找操作码
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte { // 创建包含操作码和可选操作数的指令，...代表可选
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	// 将一个 uint16 类型的整数值（o）以大端字节序的形式写入到一个字节切片（instruction）的指定位置（从 offset 开始）
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) { // Make的逆向操作，返回解码后的操作数
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 { // 从给定的 Instructions 类型变量中读取一个 16 位无符号整数，并使用大端序解释数据
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) String() string { // 返回指令的可读形式
	var out bytes.Buffer

	i := 0
	for i <= len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string { // 打印操作数
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match definition %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0: // 如果操作数数量为 0，直接返回指令的名称。
		return def.Name
	case 1: //  如果操作数数量为 1，返回指令名称和第一个操作数的字符串表示形式。
		return fmt.Sprintf("%s, %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
