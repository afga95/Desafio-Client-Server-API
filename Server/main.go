package main

import (
	"context"
	"desafio-server-api/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ContacaoBRL struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	db.InitDB()

	http.HandleFunc("/cotacao", BuscaCambioDolarHandler)

	println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuscaCambioDolar(ctx context.Context) (*ContacaoBRL, error) {
	// Faz a requisição para a API
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Lê a resposta da API e converte para um slice de bytes
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Converte o slice de bytes para uma struct
	var cotacao ContacaoBRL
	err = json.Unmarshal(res, &cotacao)
	if err != nil {
		return nil, err
	}

	// Retorna a struct com a cotação do dólar
	return &cotacao, nil
}

func BuscaCambioDolarHandler(w http.ResponseWriter, r *http.Request) {
	// Cria um contexto com timeout de 200ms para buscar a cotação do dólar
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	log.Println("Iniciando busca da cotação do dólar")
	defer log.Println("Finalizando busca da cotação do dólar")

	// Busca a cotação do dólar
	contacaoBRL, err := BuscaCambioDolar(ctx)
	if err != nil {
		log.Printf("Erro ao buscar cotação: %v\n", err)
		http.Error(w, "Erro ao buscar a cotação do dólar", http.StatusInternalServerError)
		return
	}

	// Cria um contexto com timeout de 10ms para a persistência no banco
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	// Salva a cotação no banco de dados
	err = db.AdicionarCotacao(ctxDB, contacaoBRL.USDBRL.Bid)
	if err != nil {
		log.Printf("Erro ao salvar no banco: %v\n", err)
		http.Error(w, "Erro ao salvar a cotação no banco de dados", http.StatusInternalServerError)
		return
	}

	// Chama a função que cria o arquivo de texto com a cotação
	criarArquivoCotacao(contacaoBRL.USDBRL.Bid)

	// Define o cabeçalho da resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Escreve a cotação no corpo da resposta
	json.NewEncoder(w).Encode(contacaoBRL)
}

func criarArquivoCotacao(bid string) {
	// Cria o arquivo de texto
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo: %v\n", err)
		return
	}
	defer file.Close()

	// Escreve a cotação no arquivo
	_, err = file.WriteString(fmt.Sprintf("Dólar: {%v}", bid))
	if err != nil {
		fmt.Printf("Erro ao escrever no arquivo: %v\n", err)
		return
	}
}
