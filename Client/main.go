package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Erro ao salvar no banco: %v\n", err)
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}
	//fmt.Printf("Valor atual do câmbio: %v\n", string(body))
	println(string(body))
}
