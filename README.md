# Summary
This is a simple Go program intended to pull down job postings from the UW Workday website and pass them into a Gemini SDK for a relevancy score. The workflow is scrape → filter → evaluate → report. The goal is to automate and deduplicate relevant job searches for the UW website. I chose to use Gemini only for the evaluation portion to limit the number of tokens. Filtering and reporting are cheap and easy to implement with code, reasoning is not. This is currently a work in progress.
- [ ] POST/GET for job information
- [ ] Filter for relevant terms and deduplicate
- [ ] Gemini reasoning
- [ ] Receive scored output
- [ ] Produce report