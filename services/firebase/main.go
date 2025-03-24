package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Substitua pela sua API KEY do Firebase
const firebaseAPIKey = "YOUR_API_KEY_HERE"

// Substitua pelo seu email e senha
const userEmail = "test@test.com"
const userPassword = "123456"

func main() {
	idToken, err := getFirebaseIDToken(userEmail, userPassword)
	if err != nil {
		log.Fatal("Erro ao obter Firebase ID Token:", err)
	}
	fmt.Println("Firebase ID Token gerado:\n", idToken)
}

func getFirebaseIDToken(email, password string) (string, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", firebaseAPIKey)

	body := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	idToken, ok := result["idToken"].(string)
	if !ok {
		return "", fmt.Errorf("idToken não encontrado na resposta: %v", string(respBody))
	}
	return idToken, nil
}
