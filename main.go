package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode" //ASCII
)

type no struct {
	valor    string
	direito  *no
	esquerdo *no
}

func precedencia(operador string) int {
	switch operador {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func parser(tokens []string) *no {
	var pilha []*no
	var filaSaida []*no

	for _, token := range tokens {
		if num, err := strconv.Atoi(token); err == nil {
			filaSaida = append(filaSaida, &no{valor: strconv.Itoa(num)})
		} else if isOperator(token) {
			for len(pilha) > 0 && precedencia(pilha[len(pilha)-1].valor) >= precedencia(token) {
				filaSaida = append(filaSaida, pilha[len(pilha)-1])
				pilha = pilha[:len(pilha)-1] //remove o topo da pilha
			}
			pilha = append(pilha, &no{valor: token})
		} else if token == "(" {
			pilha = append(pilha, &no{valor: token})
		} else if token == ")" {
			for len(pilha) > 0 && pilha[len(pilha)-1].valor != "(" {
				filaSaida = append(filaSaida, pilha[len(pilha)-1])
				pilha = pilha[:len(pilha)-1]
			}
			pilha = pilha[:len(pilha)-1]
		}
	}

	// Tira os elementos restantes da pilha para a fila
	for len(pilha) > 0 {
		filaSaida = append(filaSaida, pilha[len(pilha)-1])
		pilha = pilha[:len(pilha)-1]
	}

	// Construção da arvore pela fila
	for len(filaSaida) > 1 {
		for i := 0; i < len(filaSaida); i++ {
			if isOperator(filaSaida[i].valor) && filaSaida[i].esquerdo == nil && filaSaida[i].direito == nil {
				// Tira o operador e operandos
				operador := filaSaida[i]
				esquerda := filaSaida[i-2]
				direita := filaSaida[i-1]

				operador.esquerdo = esquerda
				operador.direito = direita

				filaAux := []*no{}

				// Adiciona os elementos anteriores ao operador
				filaAux = append(filaAux, filaSaida[:i-2]...)

				// Adiciona o operador (que agora é o nó da subárvore)
				filaAux = append(filaAux, operador)

				// Adiciona os elementos restantes após os dois operandos e o operador
				filaAux = append(filaAux, filaSaida[i+1:]...)

				filaSaida = filaAux

				break
			}
		}
	}

	return filaSaida[0]
}

// Função que converte a árvore de expressão para uma string
func toString(node *no) string {
	if node == nil {
		return ""
	}

	// Se o nó é folha, retorna o valor
	if node.esquerdo == nil && node.direito == nil {
		return node.valor
	}

	// Concatena a expressão em ordem in-fixa
	leftStr := toString(node.esquerdo)
	rightStr := toString(node.direito)

	return fmt.Sprintf("%s %s %s", leftStr, node.valor, rightStr)
}

func reader(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Não foi possível abrir o arquivo: %v", err)
	}
	defer file.Close()

	var expressions []string

	//Le o arquivo linha por linha
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			expressions = append(expressions, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Erro ao ler o arquivo: %v", err)
	}

	return expressions, nil
}

func lexer(input string) []string {
	var tokens []string
	var currentToken string

	for i := 0; i < len(input); i++ {
		caracter := rune(input[i])

		if unicode.IsSpace(caracter) {
			continue
		}

		switch caracter {
		case '+', '*', '/', '(', ')':
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(caracter))

		case '-':
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}

			// Verifica se o '-' é um operador unário ou binário
			if len(tokens) == 0 || tokens[len(tokens)-1] == "(" || tokens[len(tokens)-1] == "+" ||
				tokens[len(tokens)-1] == "-" || tokens[len(tokens)-1] == "*" || tokens[len(tokens)-1] == "/" {
				currentToken = "-" //binário
			} else {
				tokens = append(tokens, "-") //unário
			}

		default:
			if unicode.IsDigit(caracter) || caracter == '.' {
				currentToken += string(caracter)
			} else {
				if currentToken != "" {
					tokens = append(tokens, currentToken)
					currentToken = ""
				}
				tokens = append(tokens, string(caracter))
			}
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

func evalStep(node *no) *no {
	if node == nil {
		return nil
	}

	if node.esquerdo != nil && node.esquerdo.esquerdo == nil && node.esquerdo.direito == nil &&
		node.direito != nil && node.direito.esquerdo == nil && node.direito.direito == nil {

		leftValue, _ := strconv.Atoi(node.esquerdo.valor)
		rightValue, _ := strconv.Atoi(node.direito.valor)
		var result int

		switch node.valor {
		case "+":
			result = leftValue + rightValue
		case "-":
			result = leftValue - rightValue
		case "*":
			result = leftValue * rightValue
		case "/":
			result = leftValue / rightValue
		}

		return &no{valor: strconv.Itoa(result)}

	} else if node.esquerdo.esquerdo != nil && node.esquerdo.direito != nil {
		node.esquerdo = evalStep(node.esquerdo)
		return node
	} else if node.direito.esquerdo != nil && node.direito.direito != nil {
		node.direito = evalStep(node.direito)
		return node
	}

	return nil
}

func resultado(arvore *no) {
	for arvore.esquerdo != nil && arvore.direito != nil {
		arvore = evalStep(arvore)
		fmt.Print("\n", toString(arvore))
	}

}

func main() {

	start := time.Now()

	filename := "Casos-de-teste.txt"

	expressions, err := reader(filename)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %v", err)
		return
	}

	for _, expr := range expressions {
		fmt.Println("\nExpressão:", expr)
		tree := parser(lexer(expr))
		fmt.Print("Resultado: ")
		resultado(tree)
		fmt.Print("\n")
	}

	fmt.Print("Tempo de execucao: ", time.Since(start))

}
