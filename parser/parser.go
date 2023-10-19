package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Mdaiki0730/hackvm/code"
)

type Parser struct {
	file    *os.File
	writer  *code.Writer
	scanner bufio.Scanner
	Command string
}

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

var arithmeticChecker = regexp.MustCompile(`^(add|sub|neg|eq|gt|lt|and|or|not)$`)
var pushalphaChecker = regexp.MustCompile(`^push constant \d+$`)
var pushbetaChecker = regexp.MustCompile(`^push (local|argument|this|that|pointer|temp|constant|static) \d+$`)
var popChecker = regexp.MustCompile(`^pop (local|argument|this|that|pointer|temp|constant|static) \d+$`)
var labelChecker = regexp.MustCompile(`^label ([a-z]|[A-Z]|[0-9]|_|.|:)+$`)
var gotoChecker = regexp.MustCompile(`^goto ([a-z]|[A-Z]|[0-9]|_|.|:)+$`)
var ifgotoChecker = regexp.MustCompile(`^if-goto ([a-z]|[A-Z]|[0-9]|_|.|:)+$`)
var functionChecker = regexp.MustCompile(`^function ([a-z]|[A-Z]|[0-9]|_|.|:)+ \d+$`)
var callChecker = regexp.MustCompile(`^call ([a-z]|[A-Z]|[0-9]|_|.|:)+ \d+$`)
var returnChecker = regexp.MustCompile(`^return$`)

func NewParser(file *os.File, writer *code.Writer) Parser {
	return Parser{
		file:    file,
		writer:  writer,
		scanner: *bufio.NewScanner(file),
	}
}

func (p *Parser) HasMoreCommands() bool {
	if !p.scanner.Scan() {
		return false
	}
	return true
}

func (p *Parser) Advance() {
	// ignore comment out
	line := p.scanner.Text()
	index := strings.Index(line, "//")
	trimmedCommentString := line
	if index != -1 {
		trimmedCommentString = line[:index]
	}
	p.Command = strings.TrimSpace(trimmedCommentString)
	if p.Command == "" || p.Command == "\n" {
		return
	}

	switch p.CommandType() {
	case C_ARITHMETIC:
		p.writer.WriteArithmetic(p.Command)
	case C_PUSH:
		p.writer.WritePushPop("push", p.arg1(), p.arg2())
	case C_POP:
		p.writer.WritePushPop("pop", p.arg1(), p.arg2())
	case C_LABEL:
		p.writer.WriteLabel(p.arg1())
	case C_GOTO:
		p.writer.WriteGoto(p.arg1())
	case C_IF:
		p.writer.WriteIf(p.arg1())
	case C_FUNCTION:
		p.writer.WriteFunction(p.arg1(), p.arg2())
	case C_CALL:
		p.writer.WriteCall(p.arg1(), p.arg2())
	case C_RETURN:
		p.writer.WriteReturn()
	default:
		fmt.Println("unexpected command")
		os.Exit(1)
	}
}

func (p *Parser) CommandType() int {
	if arithmeticChecker.MatchString(p.Command) {
		return C_ARITHMETIC
	} else if pushbetaChecker.MatchString(p.Command) {
		return C_PUSH
	} else if popChecker.MatchString(p.Command) {
		return C_POP
	} else if labelChecker.MatchString(p.Command) {
		return C_LABEL
	} else if gotoChecker.MatchString(p.Command) {
		return C_GOTO
	} else if ifgotoChecker.MatchString(p.Command) {
		return C_IF
	} else if functionChecker.MatchString(p.Command) {
		return C_FUNCTION
	} else if callChecker.MatchString(p.Command) {
		return C_CALL
	} else if returnChecker.MatchString(p.Command) {
		return C_RETURN
	}
	fmt.Println("unexpected token")
	os.Exit(1)
	return 0
}

func (p *Parser) arg1() string {
	return strings.Split(p.Command, " ")[1]
}

func (p *Parser) arg2() int {
	strArg2 := strings.Split(p.Command, " ")[2]
	arg2, err := strconv.Atoi(strArg2)
	if err != nil {
		fmt.Println("can't get arg2")
		os.Exit(1)
	}
	return arg2
}
