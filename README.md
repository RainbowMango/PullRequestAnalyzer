#Install
需要安装`github.com/gocolly/colly`及其依赖

# Run
## 1. 设置PR合入的时间
main.go中手动修改
```
	// 获取指定日期区间的PR数据, (startDate, endDate)
	startDate := time.Date(2019, time.July, 15, 0, 0, 1, 0, time.Local).UTC()
	endDate := time.Date(2019, time.July, 18, 23, 59, 59, 0, time.Local).UTC()
```
这里设置时间使用本地时间方便设置，由于网页中获取的时间为UTC，所以本处设置的本地时间也会转成UTC时间。

TODO: 后续可以考虑使用参数传递

## 2. 启动
直接`go run main.go`

# Output
运行过程中除打印过程信息外，最后爬取结束会打印所有爬取到的PR汇总信息。

比如，查询2019年7月15日合入的PR`[2019.07.15-2019.07.16)`相应的配置为:
```	
startDate := time.Date(2019, time.July, 15, 0, 0, 0, 0, time.Local)
endDate := time.Date(2019, time.July, 16, 0, 0, 0, 0, time.Local)
```

输出结果为：
```
PR: https://github.com/kubernetes/kubernetes/pull/79920, Kind: kind/bug, Merged At: 2019-07-15 23:39:21 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/78774, Kind: kind/bug, Merged At: 2019-07-15 23:39:08 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/76239, Kind: kind/cleanup, Merged At: 2019-07-15 19:43:07 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/80103, Kind: kind/bug, Merged At: 2019-07-15 18:23:06 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/80151, Kind: kind/cleanup, Merged At: 2019-07-15 12:47:05 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/78263, Kind: kind/cleanup, Merged At: 2019-07-15 04:19:07 +0800 CST
PR: https://github.com/kubernetes/kubernetes/pull/80141, Kind: kind/bug, Merged At: 2019-07-15 00:55:03 +0800 CST
Finally Got 7 PRs.
kindcleanupNum:  3
kindapichangeNum:  0
kindbugNum:  4
kindfeatureNumber:  0
kindOtherNumber:  0
```