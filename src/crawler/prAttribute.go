package crawler

import (
	"common"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

// 获取PR的属性
func GetPRLables(targetPage string) common.PullRequestItem {
	var PullRequest common.PullRequestItem
	PullRequest.URL = targetPage

	c := colly.NewCollector()

	/*
		<div class="labels css-truncate js-issue-labels">
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #0ffa16; color: #000000" data-name="approved" title="Indicates a PullRequestItem has been approved by an approver from all required OWNERS files." href="/kubernetes/kubernetes/labels/approved"><span class="css-truncate-target" style="max-width: 100%">approved</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #0052cc; color: #ffffff" data-name="area/conformance" title="Issues or PRs related to kubernetes conformance tests" href="/kubernetes/kubernetes/labels/area%2Fconformance"><span class="css-truncate-target" style="max-width: 100%">area/conformance</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #0052cc; color: #ffffff" data-name="area/test" title="area/test" href="/kubernetes/kubernetes/labels/area%2Ftest"><span class="css-truncate-target" style="max-width: 100%">area/test</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #bfe5bf; color: #000000" data-name="cncf-cla: yes" title="Indicates the PullRequestItem&#39;s author has signed the CNCF CLA." href="/kubernetes/kubernetes/labels/cncf-cla%3A%20yes"><span class="css-truncate-target" style="max-width: 100%">cncf-cla: yes</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #e11d21; color: #ffffff" data-name="kind/bug" title="Categorizes issue or PullRequestItem as related to a bug." href="/kubernetes/kubernetes/labels/kind%2Fbug"><span class="css-truncate-target" style="max-width: 100%">kind/bug</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #15dd18; color: #000000" data-name="lgtm" title="Indicates that a PullRequestItem is ready to be merged." href="/kubernetes/kubernetes/labels/lgtm"><span class="css-truncate-target" style="max-width: 100%">lgtm</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #eb6420; color: #000000" data-name="priority/important-soon" title="Must be staffed and worked on either currently, or very soon, ideally in time for the next release." href="/kubernetes/kubernetes/labels/priority%2Fimportant-soon"><span class="css-truncate-target" style="max-width: 100%">priority/important-soon</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #c2e0c6; color: #000000" data-name="release-note-none" title="Denotes a PullRequestItem that doesn&#39;t merit a release note." href="/kubernetes/kubernetes/labels/release-note-none"><span class="css-truncate-target" style="max-width: 100%">release-note-none</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #d2b48c; color: #000000" data-name="sig/testing" title="Categorizes an issue or PullRequestItem as relevant to sig-testing." href="/kubernetes/kubernetes/labels/sig%2Ftesting"><span class="css-truncate-target" style="max-width: 100%">sig/testing</span></a>
		  <a class="sidebar-labels-style box-shadow-none width-full d-block IssueLabel v-align-text-top" style="background-color: #009900; color: #000000" data-name="size/XS" title="Denotes a PullRequestItem that changes 0-9 lines, ignoring generated files." href="/kubernetes/kubernetes/labels/size%2FXS"><span class="css-truncate-target" style="max-width: 100%">size/XS</span></a>
		</div>
	*/
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Attr("class"), "sidebar-labels-style") {  // 必须是边栏的label
			return
		}

		if !strings.Contains(e.Attr("href"), "/kubernetes/kubernetes/labels") {
			return
		}

		label := e.Attr("data-name")
		if label != "" {
			PullRequest.Labels = append(PullRequest.Labels, label)

			if strings.Contains(label, "kind/") && PullRequest.Kind == "" { // PR可能会有多个kind，第一个kind做为PR类型标识
				PullRequest.Kind = label
			}
			//fmt.Printf("PullRequestItem: %s, lable : %s\n", PullRequest.URL, e.Attr("data-name"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		fmt.Println("Something is wrong: ", e)
	})

	c.Visit(targetPage)

	return PullRequest
}