package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8000/api/v1"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	// Wait for server to start
	fmt.Println("Waiting for server to be ready...")
	time.Sleep(2 * time.Second)

	email := fmt.Sprintf("test_%d@example.com", time.Now().Unix())
	password := "secret123"

	// 1. Register
	fmt.Println("\n1. Testing Registration...")
	regReq := RegisterRequest{
		Name:     "Test User",
		Email:    email,
		Password: password,
	}
	if err := sendRequest("POST", "/register", regReq); err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	fmt.Println("Registration successful!")

	// 2. Login
	fmt.Println("\n2. Testing Login...")
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}
	if err := sendRequest("POST", "/login", loginReq); err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Println("Login successful!")
}

func sendRequest(method, endpoint string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(respBody))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return nil
}
