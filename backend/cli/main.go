package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	usuarioSend, err := json.Marshal(usuario)

	result := SendPost(usuarioSend)
	fmt.Printf("RESULTADO: %s", result)
}
