package main

import (
	"log/slog"

	"UWScraper/constants"
	"UWScraper/jobrequest"

	"github.com/joho/godotenv"
)

const (
	batchRequests  int    = 5
	profileDefault string = "profile/profile.md"
)

var searchTerms = [...]string{"robotics", "control", "APL", "autonomous", "Rust"}

func main() {
	workdayClient := jobrequest.NewWorkdayClient()
	searchTerms := constants.SearchTerms
	jobPostingMap := make(map[string]jobrequest.JobPosting)
	for _, term := range searchTerms {
		if err := workdayClient.SendWorkdayPost(term, jobPostingMap); err != nil {
			slog.Warn("error creating posting map", "error", err)
		}
	}
	jobs := make([]jobrequest.JobInfo, 0, len(jobPostingMap))
	for _, jobPosting := range jobPostingMap {
		jobDescription, err := workdayClient.GetWorkdayDescription(jobPosting.ExternalPath)
		if err != nil {
			slog.Warn("error getting job description", "error", err)
			continue
		}
		jobs = append(jobs, jobrequest.ConstructJobOpening(jobPosting, jobDescription))
	}

	if err := godotenv.Load(); err != nil {
		slog.Debug("error in loading environment", "error", err)
	}
}
