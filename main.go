package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}

	port := ":8080"

	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/public/", http.StripPrefix("/public", fs))

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			r.ParseForm()
			mermaidDiagram := r.FormValue("mermaid-input")

			terraFormString, err := makeOpenAIRequest(mermaidDiagram, openAIKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, "%s", terraFormString)

		}
	})

	fmt.Printf("listening on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeOpenAIRequest(mermaidDiagram string, apiKey string) (string, error) {

	requestMessage := fmt.Sprintf("Analyze the provided Mermaid diagram to identify all resources, their types, attributes, and interconnections. Translate these details into a Terraform configuration using HCL syntax. Structure the Terraform configuration accurately to reflect the architecture shown in the diagram, considering resource dependencies and relationships. The response should be formatted as an HTML code block. Wrap the Terraform configuration in <pre> and <code> HTML tags to ensure proper formatting and readability when embedded in a web application. Ensure the syntax is valid and adheres to Terraform's best practices. %s", mermaidDiagram)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": requestMessage},
		},
		"temperature": 0.7,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	choices := response["choices"].([]interface{})
	if len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})
		content := message["content"].(string)
		return content, nil
	}

	return "", fmt.Errorf("Error")

}
