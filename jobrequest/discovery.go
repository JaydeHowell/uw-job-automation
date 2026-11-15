package jobrequest

import (
	"fmt"
	"strings"
	"time"
)

const (
	BaseApiPrefix = "https://uw.wd5.myworkdayjobs.com/wday/cxs/uw/UWHires"
	postURL       = BaseApiPrefix + "/jobs"
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
