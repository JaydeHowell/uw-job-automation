package jobrequest

import "time"

type JobPostingInfo struct {
	Title      string
	Location   string
	PostedDate time.Time
	TimeType   string
	ReqID      string
	JobDetails string
	URL        string
}

func ConstructJobOpening(posting JobPosting, description JobDescription) JobPostingInfo {
	return JobPostingInfo{
		Title:      posting.JobTitle,
		Location:   posting.Location,
		PostedDate: posting.PostedDate,
		TimeType:   description.TimeType,
		ReqID:      posting.ReqID,
		JobDetails: description.JobDetails,
		URL:        description.ExternalURL,
	}
}
