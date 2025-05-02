// Package openai implements the OpenAI language model integration for Gendo.
// It provides functionality to interact with OpenAI's API, supporting
// configurable models, API keys, and base URLs. The package handles
// authentication, request formatting, and response parsing.
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gendo/pkg/log"
)

// LLM implements the llm.LLM interface for OpenAI
type LLM struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents the request body for OpenAI API
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Response represents the response from OpenAI API
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// New creates a new OpenAI LLM
func New(apiKey string, cliModel string) *LLM {
	// Try GENDO_API_KEY first, then fall back to OPENAI_API_KEY
	if apiKey == "" {
		apiKey = os.Getenv("GENDO_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
	}

	// Try GENDO_API_BASE first, then fall back to standard OpenAI environment variables
	baseURL := os.Getenv("GENDO_API_BASE")
	if baseURL == "" {
		baseURL = os.Getenv("OPENAI_API_BASE")
		if baseURL == "" {
			baseURL = os.Getenv("OPENAI_BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:9100/v1"
			}
		}
	}
	baseURL = strings.TrimRight(baseURL, "/")

	// Try CLI model first, then GENDO_MODEL, then default to bitnet
	model := cliModel
	if model == "" {
		model = os.Getenv("GENDO_MODEL")
		if model == "" {
			model = "bitnet"
		}
	}

	return &LLM{
		apiKey:     apiKey,
		baseURL:    baseURL,
		model:      model,
		httpClient: &http.Client{},
	}
}

// Process implements the llm.LLM interface
func (l *LLM) Process(prompt, input string) (string, error) {
	log.Debug("Processing with OpenAI LLM - Model: %s, Prompt: %q, Input: %q", l.model, prompt, input)

	if l.apiKey == "" {
		log.Debug("No API key set, passing through input")
		return input, nil
	}

	reqBody := Request{
		Model: l.model,
		Messages: []Message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: input},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Debug("Failed to marshal request: %v", err)
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", l.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Debug("Failed to create request: %v", err)
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)

	log.Debug("Sending request to %s with body: %s", req.URL.String(), string(jsonData))
	resp, err := l.httpClient.Do(req)
	if err != nil {
		log.Debug("Failed to call OpenAI API: %v", err)
		return "", fmt.Errorf("failed to call OpenAI API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug("Failed to read response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	log.Debug("Response from OpenAI API: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Debug("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp Response
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&openAIResp); err != nil {
		log.Debug("Failed to decode OpenAI response: %v", err)
		return "", fmt.Errorf("failed to decode OpenAI response: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		log.Debug("No response from OpenAI API")
		return "", fmt.Errorf("no response from OpenAI API")
	}

	result := openAIResp.Choices[0].Message.Content
	log.Debug("OpenAI LLM returned: %q", result)
	return result, nil
}
