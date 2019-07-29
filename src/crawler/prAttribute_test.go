package crawler

import (
	"strings"
	"testing"
)

func TestGetTitle(t *testing.T) {
	page := "https://github.com/kubernetes/kubernetes/pull/80259"
	expected := "Kubeadm Networking Configuration E2E Tests"

	title := GetTitle(page)
	if !strings.Contains(title, expected) {
		t.Errorf("Failed get pr title from page: %s\n", page)
		return
	}

	t.Logf("Get PR: %s title successfull: %s", page, title)
}
