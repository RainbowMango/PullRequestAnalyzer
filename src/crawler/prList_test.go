package crawler

import (
	"strings"
	"testing"
)

func TestCrawlPrListFromPageExperiment(t *testing.T) {
	page := "https://github.com/kubernetes/kubernetes/commits/master"

	PullRequests := CrawlPrListFromPageExperiment(page)
	if len(PullRequests) == 0 {
		t.Errorf("No pull request found in page: %s\n", page)
		return
	}

	for _, pr := range PullRequests {
		if !strings.Contains(pr.URL, "https://github.com/kubernetes/kubernetes/pull/") {
			t.Errorf("pr.URL fomat error: %s\n", pr.URL)
			return
		}
	}
}
