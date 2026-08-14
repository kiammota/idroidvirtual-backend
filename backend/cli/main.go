package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
)

type Usuario struct {
	ID        int64     `json:"id"`
	EmpresaID int64     `json:"empresa_id"`
	Email     string    `json:"email"`
	Nome      string    `json:"nome"`
	Senha     string    `json:"senha"`
	CriadoEm  time.Time `json:"criado_em"`
}

func main() {
	fmt.Println("CADASTRO DE USUÁRIO")

	var usuario Usuario
	var empresaID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Empresa ID: ").
				Value(&empresaID),

			huh.NewInput().
				Title("Nome: ").
				Value(&usuario.Nome),

			huh.NewInput().
				Title("Email: ").
				Value(&usuario.Email).
				Validate(func(str string) error {
					if !strings.Contains(str, "@") || !strings.Contains(str, ".") {
						return errors.New("o e-mail deve conter '@' e um ponto '.'")
					}

					return nil
				}),

			huh.NewInput().
				Title("Senha: ").
				Value(&usuario.Senha).
				EchoMode(huh.EchoModePassword),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Normalização
	usuario.Nome = strings.ToUpper(strings.TrimSpace(usuario.Nome))
	usuario.Email = strings.TrimSpace(usuario.Email)

	// Empresa ID
	id, err := strconv.ParseInt(strings.TrimSpace(empresaID), 10, 64)
	if err != nil {
		fmt.Println("erro: empresa ID inválido:", err)
		return
	}

	usuario.EmpresaID = id

	// Debug: mostra exatamente o que será enviado
	fmt.Println()
	fmt.Println("--- DADOS ---")
	fmt.Println("Empresa ID:", usuario.EmpresaID)
	fmt.Println("Nome:", usuario.Nome)
	fmt.Println("Email:", usuario.Email)
	fmt.Println("-------------")

	// JSON
	usuarioEnviado, err := json.Marshal(usuario)
	if err != nil {
		fmt.Println("não foi possível criar JSON do usuário:", err)
		return
	}

	fmt.Println()
	fmt.Println("--- JSON ENVIADO ---")
	fmt.Println(string(usuarioEnviado))
	fmt.Println("--------------------")

	// Requisição
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/users",
		bytes.NewBuffer(usuarioEnviado),
	)
	if err != nil {
		fmt.Println("erro ao criar requisição:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Println()
	fmt.Println("Enviando requisição...")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("erro ao enviar requisição para o servidor:", err)
		return
	}

	defer resp.Body.Close()

	// Lê resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("erro ao ler resposta do servidor:", err)
		return
	}

	fmt.Println()
	fmt.Println("--- RESPOSTA DA API ---")
	fmt.Println("Status:", resp.Status)

	if len(body) > 0 {
		fmt.Println("Body:", string(body))
	} else {
		fmt.Println("Body: <vazio>")
	}

	fmt.Println("-----------------------")

	// Interpreta status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println()
		fmt.Println("✓ Usuário cadastrado com sucesso!")
		return
	}

	fmt.Println()
	fmt.Println("✗ Não foi possível cadastrar o usuário.")

	switch resp.StatusCode {
	case http.StatusBadRequest:
		fmt.Println("A API rejeitou os dados enviados.")
	case http.StatusUnauthorized:
		fmt.Println("Você não está autenticado.")
	case http.StatusForbidden:
		fmt.Println("Você não possui permissão para essa operação.")
	case http.StatusConflict:
		fmt.Println("Já existe um usuário com esses dados.")
	case http.StatusInternalServerError:
		fmt.Println("A API encontrou um erro interno.")
	default:
		fmt.Println("A API retornou um erro inesperado.")
	}
}
