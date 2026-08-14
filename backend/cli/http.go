package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	Data    any    `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const (
	API_URL = "http://localhost:8080/"
)

func deserializeResult(body []byte) (Result, error) {
	var result Result

	fmt.Printf("BODY RECEBIDO: %q\n", string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return Result{}, fmt.Errorf(
			"DESSERIALIZAÇÃO INVÁLIDA: %w",
			err,
		)
	}

	return result, nil
}


func SendGet(route string) Result {
	req, err := http.NewRequest(
		http.MethodGet,
		API_URL+route,
		nil,
	)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	var result Result
	result, err = deserializeResult(body)
	if err !=nil {
		println("erro tentando deserializar result: ", err.Error())
	}

	return result
}

func SendPost(body []byte) Result {
	req, err := http.NewRequest(
		http.MethodPost,
		API_URL,
		bytes.NewReader(body),
	)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}

	var result Result

	result, err = deserializeResult(responseBody)	
	if err != nil {
		println("erro tentando deserializar result:", err.Error())
	}

	return result
}
