package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/parser"
)

const (
	PROMPT = ">> "
	EXIT   = "exit"
	QUIT   = "q"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		if line == EXIT || line == QUIT {
			io.WriteString(out, "\nSayonara...(_ _)m")
			break
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
