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

	if err := json.Unmarshal(body, &result); err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("resposta inválida da API: %s", err),
		}
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

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf(
				"resposta inválida da API: %s",
				err,
			),
		}
	}

	return result
}
