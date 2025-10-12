package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/koplec/gospl/internal/reader"
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Gospl REPL")

	for {
		fmt.Print("> ")

		if !scanner.Scan() { //ctrl+Dで、EOFシグナルが送られ、falseになって、終わり。
			break
		}

		input := scanner.Text()

		// Read 入力をS式に変換
		parser := reader.NewParser(input)
		expr, err := parser.Parse()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		result := expr
		fmt.Println(result.String())
	}
}
