package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"log/slog"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
)

const geminiAuthHeader = "x-goog-api-key"
const geminigenerateContent = "/models/gemini-2.5-flash-lite:generateContent"
const generatorSeconds = 5
const defaultTimeout = 60

type GeminiFileData struct {
	MimeType string `json:"mime_type"`
	FileURI  string `json:"file_uri"`
}

type GeminiPart struct {
	Text     string          `json:"text,omitempty"`
	FileData *GeminiFileData `json:"file_data,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
	Tools    []interface{}   `json:"tools,omitempty"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiService struct {
	promptManaget PromptManager
	httpClient    *http.Client
}

func NewGeminiService() *GeminiService {
	return &GeminiService{
		promptManaget: *NewPromptManager(),
		httpClient: &http.Client{
			Timeout: time.Second * defaultTimeout,
		},
	}
}

func GetStreamResponse(ctx context.Context, prompt string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		for i := 0; i <= generatorSeconds; i++ {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				ch <- fmt.Sprintf("token %d", i)
			}
		}
	}()

	return ch

}

func getApiKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing GEMINI_API_KEY env variable")
	}
	return apiKey, nil
}

func getBaseUrl() (string, error) {
	apiKey := os.Getenv("GEMINI_BASE_URL")
	if apiKey == "" {
		return "", fmt.Errorf("missing GEMINI_BASE_URL env variable")
	}
	return apiKey, nil
}

func (s *GeminiService) GenResponse(ctx context.Context, prompt string) (string, error) {
	slog.Info("Starting GenResponse")

	baseURL, err := getBaseUrl()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(baseURL, geminigenerateContent)
	if err != nil {
		return "", fmt.Errorf("failed to join URL: %w", err)
	}

	method := "POST"

	parts := []GeminiPart{
		{Text: prompt},
	}

	payload := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: parts,
			},
		},
		Tools: []any{
			map[string]any{
				"google_search": map[string]any{},
			},
		},
	}

	return s.executeRequest(ctx, method, url, payload)
}

func (s *GeminiService) GetResponseWithMedia(ctx context.Context, fileLink string, prompt string) (string, error) {
	baseURL, err := getBaseUrl()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(baseURL, geminigenerateContent)
	if err != nil {
		return "", fmt.Errorf("failed to join URL: %w", err)
	}

	method := "POST"

	systemPrompt := fmt.Sprintf(`%s`, prompt)

	parts := []GeminiPart{
		{Text: systemPrompt},
	}

	mimeType := helpers.GetMimeTypeFromUrl(fileLink)

	parts = append(parts, GeminiPart{
		FileData: &GeminiFileData{
			MimeType: mimeType,
			FileURI:  fileLink,
		},
	})

	payload := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: parts,
			},
		},
		Tools: []any{
			map[string]any{
				"google_search": map[string]any{},
			},
		},
	}

	return s.executeRequest(ctx, method, url, payload)
}

func (s *GeminiService) GenResponseWithPreset(ctx context.Context, userInput string, promptType PromptType) (string, error) {
	preset, err := s.promptManaget.GetPromptPreset(promptType)
	if err != nil {
		return "", err
	}

	fullPrompt := preset.FormatPrompt(userInput)

	return s.GenResponse(ctx, fullPrompt)
}

func (s *GeminiService) GenResponseWithMediaPreset(ctx context.Context, userInput string, fileLink string, promptType PromptType) (string, error) {
	preset, err := s.promptManaget.GetPromptPreset(promptType)
	if err != nil {
		return "", err
	}

	fullPrompt := preset.FormatPrompt(userInput)

	return s.GetResponseWithMedia(ctx, fileLink, fullPrompt)
}

func (s *GeminiService) executeRequest(ctx context.Context, method string, url string, payload any) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	apiKey, err := getApiKey()
	if err != nil {
		return "", err
	}
	req.Header.Set(geminiAuthHeader, apiKey)

	req = req.WithContext(ctx)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result GeminiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal Gemini response: %w", err)
	}
	slog.Info("Response JSON parsed")

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	slog.Info("Returning final text from Gemini")
	return result.Candidates[0].Content.Parts[0].Text, nil
}
