package main

import (
	"common"
	"crawler"
	"fmt"
)

var prs []common.PullRequestItem

const KubernetesMasterCommitPage = "https://github.com/kubernetes/kubernetes/commits/master"

const (
	KindcleanupLable = "kind/cleanup"
	KindapichangeLable = "kind/api-change"
	KindbugLable = "kind/bug"
	Kindfeature = "kind/feature"
)

func main() {

	// 先获取第一页的PR列表
	prList := crawler.CrawlPrListFromPage(KubernetesMasterCommitPage)
	prs = append(prs, prList...)

	// 遍历PR,获取PR属性
	for index, pr := range prs {
		prWithAtrribute := crawler.GetPRLables(pr.URL)
		prs[index].Labels = append(prs[index].Labels, prWithAtrribute.Labels...)
		prs[index].Kind = prWithAtrribute.Kind
		prs[index].MergeTime = prWithAtrribute.MergeTime
	}

	// 分析结果
	kindcleanupNum := 0
	kindbugNum := 0
	kindapichangeNum := 0
	kindfeatureNumber := 0
	kindOtherNumber := 0

	for _, pr := range prs {
		switch pr.Kind {
		case KindcleanupLable:
			kindcleanupNum++
		case KindapichangeLable:
			kindapichangeNum++
		case KindbugLable:
			kindbugNum++
		case Kindfeature:
			kindfeatureNumber++
		default:
			kindOtherNumber++
		}
	}

	fmt.Printf("Finally Got %d PRs.\n", len(prs))
	fmt.Println("kindcleanupNum: ", kindcleanupNum)
	fmt.Println("kindapichangeNum: ", kindapichangeNum)
	fmt.Println("kindbugNum: ", kindbugNum)
	fmt.Println("kindfeatureNumber: ", kindfeatureNumber)
	fmt.Println("kindOtherNumber: ", kindOtherNumber)

	// 获取下一页链接
	nextPageLink := crawler.GetNextPageLink(KubernetesMasterCommitPage)
	fmt.Printf("Next page: %s\n", nextPageLink) // 暂时不考虑处理下一页
}
