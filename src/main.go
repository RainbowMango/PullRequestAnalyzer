package main

import (
	"common"
	"crawler"
	"fmt"
	"time"
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

	// 获取指定日期区间的PR数据, [startDate, endDate)
	//startDate := time.Date(2019, time.July, 15, 0, 0, 0, 0, time.UTC)
	//endDate := time.Date(2019, time.July, 16, 0, 0, 0, 0, time.UTC)
	startDate := time.Date(2019, time.July, 15, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2019, time.July, 16, 0, 0, 0, 0, time.Local)

	// 循环获取数据
	nexPage := KubernetesMasterCommitPage
	pageIndex := 0
	shouldStop := false
	for !shouldStop{
		fmt.Printf("Get PR from page : %s, index: %d\n", nexPage, pageIndex)
		prList := crawler.CrawlPrListFromPage(nexPage) // 此时pr拥有 pr.URL, pr.MergeTime

		// 向本页中获取的Pr 列表中填充数据
		for index, pr := range prList {
			// 如果PR合入时间早于指定时间，则退出循环
			if pr.MergeTime.Before(startDate){
				fmt.Printf("Found PR(%s) merged at %s, before merge time:%s\n", pr.URL, pr.MergeTime.Local().String(), startDate.String())
				shouldStop = true
				break
			}

			// 因为是从前往后查找，前面的可能不在统计区间内，所以前面的只需要忽略，不需要退出
			if pr.MergeTime.After(endDate) {
				continue
			}

			prList[index].MergeTime = pr.MergeTime

			// 获取PR label 和kind
			prWithAttribute := crawler.GetPRLables(pr.URL)
			prList[index].Labels = append(prList[index].Labels, prWithAttribute.Labels...)
			prList[index].Kind = prWithAttribute.Kind

			// 将该PR追加到全局列表中
			prs = append(prs, prList[index])
		}

		if !shouldStop {
			nexPage = crawler.GetNextPageLink(nexPage)
			pageIndex++
		}
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
		fmt.Printf("PR: %s, Kind: %s, Merged At: %s\n", pr.URL, pr.Kind, pr.MergeTime.Local().String())
	}

	fmt.Printf("Finally Got %d PRs.\n", len(prs))
	fmt.Println("kindcleanupNum: ", kindcleanupNum)
	fmt.Println("kindapichangeNum: ", kindapichangeNum)
	fmt.Println("kindbugNum: ", kindbugNum)
	fmt.Println("kindfeatureNumber: ", kindfeatureNumber)
	fmt.Println("kindOtherNumber: ", kindOtherNumber)
}
