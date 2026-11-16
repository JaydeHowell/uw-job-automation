package jobrequest

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type JobPostingInfo struct {
	Title          string
	Location       string
	PostedDate     time.Time
	TimeType       string
	ReqID          string
	JobDescription string
	URL            string
}

type JobDescriptionDTO struct {
	JobPostingInfo struct {
		Title          string `json:"title"`
		JobDescription string `json:"jobDescription"`
		ExternalURL    string `json:"externalUrl"`
		TimeType       string `json:"timeType"`
		ReqID          string `json:"jobReqId"`
	} `json:"jobPostingInfo"`
}

func (c *WorkdayClient) GetJobDescription(externalPath string) (JobDescriptionDTO, error) {
	jobURL := c.baseURL + externalPath
	request, err := http.NewRequest(http.MethodGet, jobURL, nil)
	if err != nil {
		return JobDescriptionDTO{}, fmt.Errorf("error in GET request: %w", err)
	}

	request.Header.Set("Accept", c.contentType)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return JobDescriptionDTO{}, fmt.Errorf("error in client GET response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Warn("error closing GET response body", "error", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return JobDescriptionDTO{}, fmt.Errorf("error in client GET response: %d: %s", response.StatusCode, response.Status)
	}

	var jobDescriptionDTO JobDescriptionDTO
	if err := json.NewDecoder(response.Body).Decode(&jobDescriptionDTO); err != nil {
		return JobDescriptionDTO{}, fmt.Errorf("error in decoding GET response JSON: %w", err)
	}

	return jobDescriptionDTO, nil
}
