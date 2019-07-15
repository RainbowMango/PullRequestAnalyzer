package crawler

import (
	"common"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

// CrawlPrListFromPage 将抓取指定页面的PR列表。
// PR列表将包含PR的链接、PR合入时间信息
/*
<div class="table-list-cell">

  <p class="commit-title h5 mb-1 text-gray-dark ">
      <a aria-label="Merge pull request #80151 from nikhita/rules-cleanup publishing: bump go versions for 1.13 and 1.14" class="message js-navigation-open" data-pjax="true" href="/kubernetes/kubernetes/commit/c6eb9a8ed51f5c63cb351e2a4c13494bf5c303a2">Merge pull request</a>
      <a class="issue-link js-issue-link"
         data-error-text="Failed to load issue title" data-id="467922764"
         data-permission-text="Issue title is private"
         data-url="https://github.com/kubernetes/kubernetes/issues/80151"
         data-hovercard-type="pull_request"
         data-hovercard-url="/kubernetes/kubernetes/pull/80151/hovercard"
         href="https://github.com/kubernetes/kubernetes/pull/80151"> // PR的链接
         #80151
      </a>
      <a aria-label="Merge pull request #80151 from nikhita/rules-cleanup publishing: bump go versions for 1.13 and 1.14"
         class="message js-navigation-open"
         data-pjax="true"
         href="/kubernetes/kubernetes/commit/c6eb9a8ed51f5c63cb351e2a4c13494bf5c303a2">
         from nikhita/rules-cleanup
      </a>

      <span class="hidden-text-expander inline">
        <button type="button" class="ellipsis-expander js-details-target" aria-expanded="false">&hellip;</button>
      </span>
  </p>

  <div class="commit-meta commit-author-section no-wrap d-flex flex-items-center mt-1">
      <div class="AvatarStack flex-self-start ">
      </div>
      <div>
		  <a href="/kubernetes/kubernetes/commits?author=k8s-ci-robot"
             class="commit-author tooltipped tooltipped-s user-mention"
             aria-label="View all commits by k8s-ci-robot">
             k8s-ci-robot
          </a>
          committed
          <relative-time datetime="2019-07-15T04:47:05Z">Jul 15, 2019</relative-time>  // 合入时间
    </div>
  </div>
  <div class="commit-desc"><pre class="text-small">publishing: bump go versions for 1.13 and 1.14</pre>
  </div>
</div>
 */
func CrawlPrListFromPage(targetPage string) []common.PullRequestItem {

	var prForPage []common.PullRequestItem

	c := colly.NewCollector()

	/*
		<a class="issue-link js-issue-link"
		   data-error-text="Failed to load issue title"
		   data-id="467063436"
		   data-permission-text="Issue title is private"
		   data-url="https://github.com/kubernetes/kubernetes/issues/80053"
		   data-hovercard-type="pull_request"
		   data-hovercard-url="/kubernetes/kubernetes/pull/80053/hovercard"
		   href="https://github.com/kubernetes/kubernetes/pull/80053">
		   #80053
		</a>
	*/
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		href := e.Attr("href")
		if strings.Contains(href, "https://github.com/kubernetes/kubernetes/pull") {
			//fmt.Println(href)   // https://github.com/kubernetes/kubernetes/pull/80053
			//fmt.Println(e.Name) // a
			//fmt.Println(e.Text) // #80053
			//fmt.Println(e.Attr("class"))

			pr := common.PullRequestItem{
				URL: href,
			}

			prForPage = append(prForPage, pr)
		}

		//e.Request.Visit(e.Attr("href"))
		//fmt.Println(e.Name)
		//fmt.Println(e.Text)
		//fmt.Println(e.Attr("ria-label"))

		//fmt.Print(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		fmt.Println("Something is wrong: ", e)
	})

	c.Visit(targetPage)

	return prForPage
}

func CrawlPrListFromPageExperiment(targetPage string) []common.PullRequestItem {

	var prForPage []common.PullRequestItem

	c := colly.NewCollector()
	c.OnHTML("div[class=table-list-cell]", func(e *colly.HTMLElement) { // 选择器另外的写法".table-list-cell"
		link, exist := e.DOM.Find("a.issue-link").Attr("href")
		if !exist {
			return
		}

		timeStr, exist := e.DOM.Find("relative-time[datetime]").Attr("datetime")
		timeStamp, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			fmt.Errorf("Trans pr commit time to timestamp failed. timestr: %s\n", timeStr)
			return
		}

		newItem := common.PullRequestItem{
			URL: link,
			MergeTime: timeStamp,
		}

		prForPage = append(prForPage, newItem)

		//fmt.Printf("%s\n", link)
		//fmt.Printf("%s\n", timeStr)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		fmt.Println("Something is wrong: ", e)
	})

	c.Visit(targetPage)

	return prForPage
}

func GetNextPageLink(targetPage string) string {
	var nextPage string

	c := colly.NewCollector()

	/*
		<a rel="nofollow" class="btn btn-outline BtnGroup-item" href="https://github.com/kubernetes/kubernetes/commits/master?after=de091d102f7d1777abbad9977af7f089743e6d68+34">Older</a>
	*/
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if !strings.Contains(e.Attr("class"), "btn") {  // 必须是个按钮
			return
		}

		if e.Text != "Older" { // 按钮名字必须是"Older"
			return
		}

		nextPage = e.Attr("href")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, e error) {
		fmt.Println("Something is wrong: ", e)
	})

	c.Visit(targetPage)

	return nextPage
}