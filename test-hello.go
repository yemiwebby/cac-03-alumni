package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	phoneID := os.Getenv("WA_PHONE_ID")
	token := os.Getenv("WA_TOKEN")
	to := os.Getenv("WA_TO_LIST")

	// Test with hello_world template
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phoneID)
	
	body := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "template",
		"template": map[string]any{
			"name":     "hello_world",
			"language": map[string]string{"code": "en_US"},
		},
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %+v\n", result)
	
	if resp.StatusCode >= 300 {
		fmt.Printf("❌ Failed to send message\n")
	} else {
		fmt.Printf("✅ Message sent successfully!\n")
	}
}
