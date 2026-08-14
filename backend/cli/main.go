package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	var usuario Usuario = Usuario{}

	var empresaID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Empresa ID: ").Value(&empresaID),
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
				Value(&usuario.Senha).EchoMode(huh.EchoModePassword),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	usuario.Nome = strings.ToUpper(strings.TrimSpace(usuario.Nome))
	usuario.Email = strings.TrimSpace(usuario.Email)

	id, err := strconv.ParseInt(strings.TrimSpace(empresaID), 10, 64)
	if err != nil {
		println("erro cadastrando usuário: ", err.Error())
		return
	}
	usuario.EmpresaID = id

	usuarioEnviado, er := json.Marshal(usuario)
	if er != nil {
		println("não foi possível criar json do usuário: ", err.Error())
		return
	}
	println("usuário cadastrado!")

	req, err := http.NewRequest("POST", "http://localhost:8080/users", bytes.NewBuffer(usuarioEnviado))
	if err != nil {
		println("erro a fazer requisição:", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("erro ao enviar requisição para o servidor: %v", err)
	}
	defer resp.Body.Close()
	println("Body do servidor: ", resp.Body)
	println("Requisição enviada! Status do servidor:", resp.Status)

}
