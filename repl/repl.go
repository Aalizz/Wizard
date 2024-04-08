package repl

import (
	"bufio"
	"fmt"
	"io"
	"my.com/myfile/evaluator"
	"my.com/myfile/lexer"
	"my.com/myfile/object"
	"my.com/myfile/parser"
	"strings"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	input := ""

	fmt.Fprintf(out, PROMPT)
	for scanner.Scan() {
		line := scanner.Text()
		input += line + "\n"

		// 检查输入是否完整

		if IsComplete(input) {
			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil {
					io.WriteString(out, evaluated.Inspect())
					io.WriteString(out, "\n")
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

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func IsComplete(input string) bool {
	// 将输入字符串拆分成行
	lines := strings.Split(input, "\n")

	// 获取最后两行
	lastLine := lines[len(lines)-1]
	secondLastLine := lines[len(lines)-2]

	// 如果最后两行都是空字符串，则认为是完整的语句
	return strings.TrimSpace(lastLine) == "" && strings.TrimSpace(secondLastLine) == ""
}
