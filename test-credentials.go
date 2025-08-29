package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	phoneID := os.Getenv("WA_PHONE_ID")
	token := os.Getenv("WA_TOKEN")

	if phoneID == "" || token == "" {
		fmt.Println("Missing WA_PHONE_ID or WA_TOKEN")
		return
	}

	fmt.Printf("Phone ID: %s\n", phoneID)
	fmt.Printf("Token (first 20 chars): %s...\n", token[:20])

	// Test API connection by getting phone number info
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s", phoneID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("API Response Status: %s\n", resp.Status)

	if resp.StatusCode >= 300 {
		// Read error response
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("Error response: %s\n", buf.String())
	} else {
		// Read success response
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("Success response: %+v\n", result)
	}
}
