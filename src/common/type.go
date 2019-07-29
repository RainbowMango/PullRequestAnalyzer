package common

import "time"

type PullRequestItem struct {
	URL       string
	Title     string
	Labels    []string
	Kind      string
	MergeTime time.Time
}
