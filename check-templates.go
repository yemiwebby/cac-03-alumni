package main

import (
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

	// Get message templates
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/message_templates", phoneID)
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	
	fmt.Printf("Status: %s\n", resp.Status)
	
	if resp.StatusCode >= 300 {
		fmt.Printf("❌ Failed to get templates\n")
		fmt.Printf("Response: %+v\n", result)
		return
	}

	// Parse templates
	if data, ok := result["data"].([]interface{}); ok {
		fmt.Printf("\n📋 Message Templates:\n")
		fmt.Printf("==================\n")
		
		for _, template := range data {
			if t, ok := template.(map[string]interface{}); ok {
				name := t["name"]
				status := t["status"]
				category := t["category"]
				language := "unknown"
				
				if lang, ok := t["language"].(string); ok {
					language = lang
				}
				
				fmt.Printf("📝 Name: %v\n", name)
				fmt.Printf("   Status: %v\n", status)
				fmt.Printf("   Category: %v\n", category)
				fmt.Printf("   Language: %v\n", language)
				fmt.Printf("   ---\n")
			}
		}
	} else {
		fmt.Printf("No templates found or unexpected response format\n")
		fmt.Printf("Response: %+v\n", result)
	}
}
