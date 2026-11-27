package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"UWScraper/jobrequest"

	"google.golang.org/genai"
)

const (
	modelName string = "gemini-4-hydrogen"
)

type JobEvaluation struct {
	ReqID     string `json:"ReqID"`
	Reasoning string `json:"Reasoning"`
	Score     int    `json:"Score"`
}

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

func GetJobEvaluations(ctx context.Context, client *genai.Client, parts []*genai.Part) ([]JobEvaluation, error) {
	config := &genai.GenerateContentConfig{}
	config.SystemInstruction = &genai.Content{
		Parts: []*genai.Part{
			{Text: SystemPrompt},
		},
	}

	config.ResponseMIMEType = "application/json"

	evaluationSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"ReqID": {
				Type:        genai.TypeString,
				Description: "The reference ID of the job",
			},
			"Reasoning": {
				Type:        genai.TypeString,
				Description: "The reasoning behind why the score was given",
			},
			"Score": {
				Type:        genai.TypeInteger,
				Description: "Evaluation score given to the job between 0 and 100",
			},
		},
		Required:         []string{"ReqID", "Reasoning", "Score"},
		PropertyOrdering: []string{"ReqID", "Reasoning", "Score"},
	}
	config.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: evaluationSchema,
	}

	response, err := client.Models.GenerateContent(ctx, modelName, []*genai.Content{{Parts: parts}}, config)
	if err != nil {
		return nil, fmt.Errorf("error in evaluation request: %w", err)
	}

	var evaluations []JobEvaluation
	if err := json.Unmarshal([]byte(response.Text()), &evaluations); err != nil {
		return nil, fmt.Errorf("decode job evaluations: %w", err)
	}
	return evaluations, nil
}

func CreateGeminiRequest(prompt string, batch []jobrequest.JobInfo) ([]*genai.Part, error) {
	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("marshaling gemini request: %w", err)
	}
	return []*genai.Part{{Text: "<job_postings>\n" + string(batchJSON) + "\n</job_postings>\n" + prompt}}, nil
}
