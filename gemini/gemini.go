package gemini

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func NewGenAIClient(ctx context.Context) (*genai.Client, error) {
	apiKey, envExists := os.LookupEnv("GEMINI_API_KEY")
	if !envExists || strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("unable to find GEMINI_API_KEY")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error in creating Gemini client: %w", err)
	}

	return client, nil

}
