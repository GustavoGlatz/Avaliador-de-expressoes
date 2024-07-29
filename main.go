package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode" //ASCII
)

func reader() string {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Digite a expressao: ")
	scanner.Scan()
	expressao := scanner.Text()

	return expressao
}

func lexer(expressao string) []string {
	var tokens []string
	var auxtoken strings.Builder // tipo de dado; "acumulador"

	addtoken := func() { // adiciona os caracteres de auxtoken em tokens; declaração de função aninhada muda
		if auxtoken.Len() > 0 {
			tokens = append(tokens, auxtoken.String()) //append(destino, origem)
			auxtoken.Reset()
		}
	}

	// "_" ignora o índice dos caracteres obtido pelo "range"
	for _, caracter := range expressao {
		switch { // executa sem restrição
		case unicode.IsDigit(caracter):
			auxtoken.WriteRune(caracter) // adiciona o número lido ao auxtoken
		case unicode.IsSpace(caracter):
			addtoken()
		default:
			addtoken()                                // útil para casos como: "... 5)"
			tokens = append(tokens, string(caracter)) // adiciona os caracteres não número e não espaço a tokens
		}
	}
	addtoken() //adiciona o(s) ultimo(s) caracter(es), como: "... 10"

	return tokens
}

func parser(tokens []string) {

}

func main() {

	//lexer(reader())

	tokens := lexer(reader())

	fmt.Printf("Tokens: %q", tokens)

}
