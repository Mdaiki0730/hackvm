package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mdaiki0730/hackvm/code"
	"github.com/Mdaiki0730/hackvm/parser"
)

func main() {
	// arg validation
	if len(os.Args) != 2 {
		fmt.Println("please enter only target file name")
		os.Exit(1)
	}

	// file setting
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("no such directory")
		os.Exit(1)
	}

	// setting code writer
	tfn := strings.Split(os.Args[1], "/")
	of, err := os.Create(fmt.Sprintf("%s/%s.asm", os.Args[1], tfn[len(tfn)-1]))
	writer := code.NewWriter(of)
	defer writer.Close()
	if err != nil {
		fmt.Println("failed to file create")
		os.Exit(1)
	}

	// main process
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) != ".vm" {
			continue
		} else {
			f, err := os.Open(os.Args[1] + "/" + file.Name())
			defer f.Close()
			if err != nil {
				fmt.Println("cannot open file")
				os.Exit(1)
			}

			writer.SetFileName(filename)
			p := parser.NewParser(f, writer)
			for p.HasMoreCommands() {
				p.Advance()
			}
		}
	}
}
