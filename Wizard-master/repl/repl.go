package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"my.com/myfile/evaluator"
	"my.com/myfile/lexer"
	"my.com/myfile/object"
	"my.com/myfile/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in) //创建一个新的扫描器
	env := object.NewEnvironment()
	input := "" //初始化

	fmt.Fprintf(out, PROMPT)
	for scanner.Scan() { //
		line := scanner.Text() //当用户输入一行并按下Enter键时，scanner.Text()会获取该行文本。
		input += line + "\n"   //累加到input字符串，直到输入被视为完整

		if IsComplete(input) { // 当输入两个空行后，进入if语句，值得注意的是，虽然实在for循环中，一个完整的语句只经历一次if语句中的程序
			l := lexer.New(input) //得到一个lexer结构体指针
			p := parser.New(l)    //在parser.New()中得到一个Parser结构体的指针，将其中的成员l初始化为参数l，并且创建token类型与解析函数的映射，方便遇到特定类型的token时调用

			program := p.ParseProgram() //创建一个ast.program结构体，并且创建所有的ast.steatment，也就是创建一个抽象语法树
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors()) //错误会通过printParserErrors函数输出到out
			} else {
				evaluated := evaluator.Eval(program, env) //返回一个Object接口
				if evaluated != nil {
					inspectedValue := evaluated.Inspect()
					if inspectedValue != "null" {
						io.WriteString(out, inspectedValue) //使用Inspect计算
						io.WriteString(out, "\n")
					} else {
						io.WriteString(out, "\n")
					}
				}
			}

			input = "" // 清空输入以准备接受下一个完整的语句
			fmt.Fprintf(out, PROMPT)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(out, "error: %v\n", err)
	}
}
func printParserErrors(out io.Writer, errors []string) { //错误输出
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func IsComplete(input string) bool { //判断完整性，但是你在程序中间不能空两行，因为那样会被视为结束
	// 将输入字符串拆分成行
	lines := strings.Split(input, "\n")

	// 获取最后两行
	lastLine := lines[len(lines)-1]
	secondLastLine := lines[len(lines)-2]

	// 如果最后两行都是空字符串，则认为是完整的语句
	return strings.TrimSpace(lastLine) == "" && strings.TrimSpace(secondLastLine) == ""
}
