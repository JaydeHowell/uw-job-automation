package jobrequest

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

const (
	BaseApiPrefix = "https://uw.wd5.myworkdayjobs.com/wday/cxs/uw/UWHires"
	postURL       = BaseApiPrefix + "/jobs"
)

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

func toJobPosting(dto workdayDTO) (PostResponse, error) {
	jobTitle := dto.Title
	externalPath := dto.ExternalPath
	location := dto.Location

	reqID := dto.BulletFields[0]
	if !strings.HasPrefix("REQ-", reqID) {
		errorMessage := fmt.Sprintf("reqID is malformed: %s", reqID)
		slog.Error(errorMessage)
		return PostResponse{}, errors.New(errorMessage)
	}

	postedDate := dto.BulletFields[0]
	parsedDate, err := time.Parse(time.DateOnly, postedDate)
	if err != nil {
		return PostResponse{}, fmt.Errorf("malformed date %d: %w", postedDate, err)
	}

	return PostResponse{
		JobTitle:     jobTitle,
		ExternalPath: externalPath,
		Location:     location,
		ReqID:        reqID,
		PostedDate:   parsedDate,
	}, nil
}
