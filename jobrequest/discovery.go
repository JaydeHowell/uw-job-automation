package jobrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaseApiPrefix     = "https://uw.wd5.myworkdayjobs.com/wday/cxs/uw/UWHires"
	postURL           = BaseApiPrefix + "/jobs"
	contentType       = "application/json"
	clientPostTimeout = 10
)

type searchEnvelope struct {
	Total       int          `json:"total"`
	JobPostings []workdayDTO `json:"jobPostings"`
}

type workdayDTO struct {
	Title        string   `json:"title"`
	ExternalPath string   `json:"externalPath"`
	Location     string   `json:"locationsText"`
	BulletFields []string `json:"bulletFields"`
}

type PostResponse struct {
	JobTitle     string
	ExternalPath string
	Location     string
	ReqID        string
	PostedDate   time.Time
}

func SendWorkdayPost(searchText string) ([]PostResponse, error) {
	payload := map[string]any{
		"appliedFacets": map[string]any{},
		"limit":         20,
		"offset":        0,
		"searchText":    searchText,
	}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(&payload); err != nil {
		return []PostResponse{}, fmt.Errorf("failed to assemble JSON: %w", err)
	}

	request, err := http.NewRequest("POST", postURL, &buffer)
	if err != nil {
		return []PostResponse{}, fmt.Errorf("error in POST send: %w", err)
	}

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", contentType)

	client := &http.Client{Timeout: clientPostTimeout * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return []PostResponse{}, fmt.Errorf("error in client response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Errorf("error closing response body: %w", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return []PostResponse{}, fmt.Errorf("error in client response: %d: %s", response.StatusCode, response.Status)
	}

	var SearchEnvelope searchEnvelope
	if err := json.NewDecoder(response.Body).Decode(&SearchEnvelope); err != nil {
		return []PostResponse{}, fmt.Errorf("error in decoding JSON: %w", err)
	}

	var JobPostings []PostResponse
	for _, posting := range SearchEnvelope.JobPostings {
		jobPosting, err := toJobPosting(posting)
		if err != nil {
			return []PostResponse{}, err
		}
		JobPostings = append(JobPostings, jobPosting)
	}
	return JobPostings, nil
}

func toJobPosting(dto workdayDTO) (PostResponse, error) {
	if len(dto.BulletFields) < 1 {
		return PostResponse{}, fmt.Errorf("bullet fields missing %q", dto.BulletFields)
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
				return PostResponse{}, fmt.Errorf("malformed date %s: %w", trimmedDate, err)
			}
			postedDate = parsedDate
			continue
		}

		continue
	}
	if reqID == "" {
		return PostResponse{}, fmt.Errorf("reqID not found: %q", dto.BulletFields)
	}

	return PostResponse{
		JobTitle:     jobTitle,
		ExternalPath: externalPath,
		Location:     location,
		ReqID:        reqID,
		PostedDate:   postedDate,
	}, nil
}
