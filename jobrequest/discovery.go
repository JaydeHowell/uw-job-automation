package jobrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type searchEnvelope struct {
	Total       int             `json:"total"`
	JobPostings []jobPostingDTO `json:"jobPostings"`
}

type jobPostingDTO struct {
	Title        string   `json:"title"`
	ExternalPath string   `json:"externalPath"`
	Location     string   `json:"locationsText"`
	BulletFields []string `json:"bulletFields"`
}

type JobPosting struct {
	JobTitle     string
	ExternalPath string
	Location     string
	ReqID        string
	PostedDate   time.Time
}

func (c *WorkdayClient) SendWorkdayPost(searchText string, AccumulatorMap map[string]JobPosting) error {
	payloadLimit := 20
	payloadOffset := 0
	totalPostings := 0
	postURL := c.baseURL + "/jobs"
	for {
		payload := map[string]any{
			"appliedFacets": map[string]any{},
			"limit":         payloadLimit,
			"offset":        payloadOffset,
			"searchText":    searchText,
		}
		SearchEnvelope, err := c.httpRequest(payload, postURL)
		if err != nil {
			return fmt.Errorf("error in http request: %w", err)
		}

		if len(SearchEnvelope.JobPostings) == 0 {
			break
		}

		for _, posting := range SearchEnvelope.JobPostings {
			jobPosting, err := toJobPosting(posting)
			if err != nil {
				slog.Warn("malformed job posting", "title", posting.Title, "error", err)
				continue
			}
			AccumulatorMap[jobPosting.ReqID] = jobPosting
		}
		if payloadOffset == 0 {
			totalPostings = SearchEnvelope.Total
		}
		if len(SearchEnvelope.JobPostings) < payloadLimit || payloadOffset >= totalPostings {
			break
		}
		payloadOffset += payloadLimit
	}
	return nil
}

func (c *WorkdayClient) httpRequest(payload map[string]any, postURL string) (searchEnvelope, error) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(&payload); err != nil {
		return searchEnvelope{}, fmt.Errorf("failed to assemble JSON: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, postURL, &buffer)
	if err != nil {
		return searchEnvelope{}, fmt.Errorf("error in POST send: %w", err)
	}

	request.Header.Set("Content-Type", c.contentType)
	request.Header.Set("Accept", c.contentType)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return searchEnvelope{}, fmt.Errorf("error in client POST response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Warn("error closing POST response body", "error", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return searchEnvelope{}, fmt.Errorf("error in client POST response: %d: %s", response.StatusCode, response.Status)
	}

	var SearchEnvelope searchEnvelope
	if err := json.NewDecoder(response.Body).Decode(&SearchEnvelope); err != nil {
		return searchEnvelope{}, fmt.Errorf("error in decoding POST response JSON: %w", err)
	}
	return SearchEnvelope, nil
}

func toJobPosting(dto jobPostingDTO) (JobPosting, error) {
	if len(dto.BulletFields) < 1 {
		return JobPosting{}, fmt.Errorf("bullet fields missing %q", dto.BulletFields)
	}

	jobTitle := dto.Title
	externalPath := dto.ExternalPath
	location := dto.Location

	var reqID string
	var postedDate time.Time

	for _, item := range dto.BulletFields {
		if strings.HasPrefix(item, "REQ-") {
			reqID = item
			continue
		}

		if strings.HasPrefix(item, "Posting Date: ") {
			splitDate := strings.Split(item, ":")
			trimmedDate := strings.Trim(splitDate[1], " ")
			parsedDate, err := time.Parse("01/02/2006", trimmedDate)
			if err != nil {
				return JobPosting{}, fmt.Errorf("malformed date %s: %w", trimmedDate, err)
			}
			postedDate = parsedDate
			continue
		}

		continue
	}
	if reqID == "" {
		return JobPosting{}, fmt.Errorf("reqID not found: %q", dto.BulletFields)
	}

	return JobPosting{
		JobTitle:     jobTitle,
		ExternalPath: externalPath,
		Location:     location,
		ReqID:        reqID,
		PostedDate:   postedDate,
	}, nil
}
