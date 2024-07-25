package main

import (
	"bufio"
	"fmt"
	"os"
)

func reader() string {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Digite a expressao: ")
	scanner.Scan() // use `for scanner.Scan()` to keep reading
	expressao := scanner.Text()

	return expressao
}

func lexer(expressao string) {

}

func main() {

	lexer(reader())

}
