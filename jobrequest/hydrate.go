package jobrequest

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type JobDescription struct {
	Title       string `json:"title"`
	JobDetails  string `json:"jobDescription"`
	ExternalURL string `json:"externalUrl"`
	TimeType    string `json:"timeType"`
	ReqID       string `json:"jobReqId"`
}

type jobDescriptionEnvelope struct {
	JobDescription JobDescription `json:"jobPostingInfo"`
}

func (c *WorkdayClient) GetWorkdayDescription(externalPath string) (JobDescription, error) {
	jobURL := c.baseURL + externalPath
	request, err := http.NewRequest(http.MethodGet, jobURL, nil)
	if err != nil {
		return JobDescription{}, fmt.Errorf("error in GET request: %w", err)
	}

	request.Header.Set("Accept", c.contentType)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return JobDescription{}, fmt.Errorf("error in client GET response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Warn("error closing GET response body", "error", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return JobDescription{}, fmt.Errorf("error in client GET response: %d: %s", response.StatusCode, response.Status)
	}

	var jobDescription jobDescriptionEnvelope
	if err := json.NewDecoder(response.Body).Decode(&jobDescription); err != nil {
		return JobDescription{}, fmt.Errorf("error in decoding GET response JSON: %w", err)
	}

	return jobDescription.JobDescription, nil
}
