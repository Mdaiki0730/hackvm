package code

import (
	"fmt"
	"os"
	"strings"
)

const (
	SP_MEMORY_DATA_TO_DREGISTA = "@SP\nA=M-1\nD=M\n" // AレジスタにSPのdata, DレジスタにSPが指すdataが保持されます
	SP_SUB_1                   = "@SP\nM=M-1\n"
	SP_ADD_1                   = "@SP\nM=M+1\n"
	SP_INITIALIZE              = "@256\nD=A\n@SP\nM=D\n"
	INSERT_DREGISTA_TO_SP      = "@SP\nA=M\nM=D\n"
)

type Writer struct {
	outfile                  *os.File
	filename                 string
	comparisonOperationIndex int
	callIndex                int
}

func NewWriter(of *os.File) Writer {
	// スタックの物理領域が256-2047なので、SP=256で初期化する
	w := Writer{outfile: of, comparisonOperationIndex: 0, callIndex: 0}
	w.WriteInit()
	return w
}

func (w *Writer) SetFileName(filename string) {
	w.filename = strings.TrimRight(filename, ".vm")
	w.write(fmt.Sprintf("// start %s\n", w.filename))
}

func (w *Writer) WriteInit() {
	w.write(SP_INITIALIZE)
	w.WriteCall("Sys.init", 0)
}

func (w *Writer) WriteArithmetic(command string) {
	switch command {
	case "add":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("M=D+M\n")
		w.write(SP_SUB_1)
	case "sub":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("M=M-D\n")
		w.write(SP_SUB_1)
	case "neg":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("M=-D\n")
	case "eq":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("D=M-D\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_TRUE%v\n", w.comparisonOperationIndex))
		w.write("D;JEQ\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_FALSE%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_TRUE%v)\n", w.comparisonOperationIndex))
		w.write("D=-1\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_FALSE%v)\n", w.comparisonOperationIndex))
		w.write("D=0\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_RESULT%v)\n", w.comparisonOperationIndex))
		w.write("@SP\n")
		w.write("A=M-1\n")
		w.write("A=A-1\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
		w.comparisonOperationIndex++
	case "gt":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("D=M-D\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_TRUE%v\n", w.comparisonOperationIndex))
		w.write("D;JGT\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_FALSE%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_TRUE%v)\n", w.comparisonOperationIndex))
		w.write("D=-1\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_FALSE%v)\n", w.comparisonOperationIndex))
		w.write("D=0\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_RESULT%v)\n", w.comparisonOperationIndex))
		w.write("@SP\n")
		w.write("A=M-1\n")
		w.write("A=A-1\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
		w.comparisonOperationIndex++
	case "lt":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("D=M-D\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_TRUE%v\n", w.comparisonOperationIndex))
		w.write("D;JLT\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_FALSE%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_TRUE%v)\n", w.comparisonOperationIndex))
		w.write("D=-1\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_FALSE%v)\n", w.comparisonOperationIndex))
		w.write("D=0\n")
		w.write(fmt.Sprintf("@COMPARISON_OPERATION_RESULT%v\n", w.comparisonOperationIndex))
		w.write("0;JMP\n")
		w.write(fmt.Sprintf("(COMPARISON_OPERATION_RESULT%v)\n", w.comparisonOperationIndex))
		w.write("@SP\n")
		w.write("A=M-1\n")
		w.write("A=A-1\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
		w.comparisonOperationIndex++
	case "and":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("M=D&M\n")
		w.write(SP_SUB_1)
	case "or":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("A=A-1\n")
		w.write("M=D|M\n")
		w.write(SP_SUB_1)
	case "not":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("M=!D\n")
	default:
		fmt.Println("unexpected arith")
		os.Exit(1)
	}
}

func (w *Writer) WritePushPop(command, segment string, index int) {
	switch command {
	case "push":
		w.WritePush(segment, index)
	case "pop":
		w.WritePop(segment, index)
	default:
		fmt.Println("unexpected vm code")
		os.Exit(1)
	}
}

func (w *Writer) WriteLabel(label string) {
	w.write(fmt.Sprintf("(%s)\n", label))
}

func (w *Writer) WriteGoto(label string) {
	w.write(fmt.Sprintf("@%s\n", label))
	w.write("0;JMP\n")
}

func (w *Writer) WriteIf(label string) {
	w.write(SP_MEMORY_DATA_TO_DREGISTA)
	w.write(SP_SUB_1)
	w.write(fmt.Sprintf("@%s\n", label))
	w.write("D;JNE\n")
}

func (w *Writer) WriteFunction(functionName string, numLocals int) {
	w.write(fmt.Sprintf("(%s)\n", functionName))
	for i := 0; i < numLocals; i++ {
		w.WritePush("constant", 0)
	}
}

func (w *Writer) WriteCall(functionName string, numArgs int) {
	w.callIndex++
	// push return-address
	w.write(fmt.Sprintf("@return-address%v\n", w.callIndex))
	w.write("D=A\n")
	w.write(INSERT_DREGISTA_TO_SP)
	w.write(SP_ADD_1)

	// push lcl
	w.write("@LCL\n")
	w.write("D=M\n")
	w.write(INSERT_DREGISTA_TO_SP)
	w.write(SP_ADD_1)

	// push arg
	w.write("@ARG\n")
	w.write("D=M\n")
	w.write(INSERT_DREGISTA_TO_SP)
	w.write(SP_ADD_1)

	// push this
	w.write("@THIS\n")
	w.write("D=M\n")
	w.write(INSERT_DREGISTA_TO_SP)
	w.write(SP_ADD_1)

	// push that
	w.write("@THAT\n")
	w.write("D=M\n")
	w.write(INSERT_DREGISTA_TO_SP)
	w.write(SP_ADD_1)

	// ARG = SP-n-5
	w.write("@SP\n")
	w.write("D=M\n")
	w.write(fmt.Sprintf("@%v\n", numArgs))
	w.write("D=D-A\n")
	w.write("@5\n")
	w.write("D=D-A\n")
	w.write("@ARG\n")
	w.write("M=D\n")

	// LCL = SP
	w.write("@SP\n")
	w.write("D=M\n")
	w.write("@LCL\n")
	w.write("M=D\n")

	// goto f
	w.WriteGoto(functionName)

	// def return-address
	w.WriteLabel(fmt.Sprintf("return-address%v", w.callIndex))
}

func (w *Writer) WriteReturn() {
	// store local to r13
	w.write("@LCL\n")
	w.write("D=M\n")
	w.write("@R13\n")
	w.write("M=D\n")

	// store ret address to r14
	w.write("@5\n")
	w.write("A=D-A\n")
	w.write("D=M\n")
	w.write("@R14\n")
	w.write("M=D\n")

	// pop retval and insert stack top
	w.write(SP_MEMORY_DATA_TO_DREGISTA)
	w.write(SP_SUB_1)
	w.write("@ARG\n")
	w.write("A=M\n")
	w.write("M=D\n")

	// restore sp
	w.write("@ARG\n")
	w.write("D=M+1\n")
	w.write("@SP\n")
	w.write("M=D\n")

	// restore that
	w.write("@R13\n")
	w.write("D=M\n")
	w.write("@1\n")
	w.write("A=D-A\n")
	w.write("D=M\n")
	w.write("@THAT\n")
	w.write("M=D\n")

	// restore this
	w.write("@R13\n")
	w.write("D=M\n")
	w.write("@2\n")
	w.write("A=D-A\n")
	w.write("D=M\n")
	w.write("@THIS\n")
	w.write("M=D\n")

	// restore arg
	w.write("@R13\n")
	w.write("D=M\n")
	w.write("@3\n")
	w.write("A=D-A\n")
	w.write("D=M\n")
	w.write("@ARG\n")
	w.write("M=D\n")

	// restore local
	w.write("@R13\n")
	w.write("D=M\n")
	w.write("@4\n")
	w.write("A=D-A\n")
	w.write("D=M\n")
	w.write("@LCL\n")
	w.write("M=D\n")

	// goto ret
	w.write("@R14\n")
	w.write("A=M\n")
	w.write("0;JMP\n")
}

func (w *Writer) WritePush(segment string, index int) {
	switch segment {
	case "constant":
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=A\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "local":
		w.write("@LCL\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "argument":
		w.write("@ARG\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "this":
		w.write("@THIS\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "that":
		w.write("@THAT\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "pointer":
		w.validatePointerIndex(index)
		w.write("@R3\n")
		w.write("D=A\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "temp":
		w.validateTempIndex(index)
		w.write("@R5\n")
		w.write("D=A\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("A=D+A\n")
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	case "static":
		w.write(fmt.Sprintf("@%s.%v\n", w.filename, index))
		w.write("D=M\n")
		w.write(INSERT_DREGISTA_TO_SP)
		w.write(SP_ADD_1)
	default:
		fmt.Println("unexpected segment")
		os.Exit(1)
	}
}

func (w *Writer) WritePop(segment string, index int) {
	switch segment {
	case "constant":
		fmt.Println("sorry, I can't understand how to implement")
		os.Exit(1)
	case "local":
		w.write("@LCL\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "argument":
		w.write("@ARG\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "this":
		w.write("@THIS\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "that":
		w.write("@THAT\n")
		w.write("D=M\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "pointer":
		w.validatePointerIndex(index)
		w.write("@R3\n")
		w.write("D=A\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "temp":
		w.validateTempIndex(index)
		w.write("@R5\n")
		w.write("D=A\n")
		w.write(fmt.Sprintf("@%v\n", index))
		w.write("D=D+A\n")
		w.write("@R13\n")
		w.write("M=D\n")
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write("@R13\n")
		w.write("A=M\n")
		w.write("M=D\n")
		w.write(SP_SUB_1)
	case "static":
		w.write(SP_MEMORY_DATA_TO_DREGISTA)
		w.write(fmt.Sprintf("@%s.%v\n", w.filename, index))
		w.write("M=D\n")
		w.write(SP_SUB_1)
	default:
		fmt.Println("unexpected segment")
		os.Exit(1)
	}
}

func (w *Writer) validatePointerIndex(index int) {
	if index == 1 || index == 0 {
		return
	}
	panic("unexpected pointer index")
}

func (w *Writer) validateTempIndex(index int) {
	if index >= 0 || index <= 7 {
		return
	}
	panic("unexpected pointer index")
}

func (w *Writer) write(bin string) error {
	_, err := w.outfile.Write([]byte(bin))
	if err != nil {
		fmt.Println("failed to write file")
		os.Exit(1)
	}
	return nil
}

func (w *Writer) Close() {
	w.outfile.Close()
}
