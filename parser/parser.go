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
	writer  code.Writer
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

func NewParser(file *os.File, writer code.Writer) Parser {
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
	default:
		fmt.Println("unexpected command")
		os.Exit(1)
	}
}

func (p *Parser) CommandType() int {
	if arithmeticChecker.MatchString(p.Command) {
		return C_ARITHMETIC
	} else if pushalphaChecker.MatchString(p.Command) {
		return C_PUSH
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
