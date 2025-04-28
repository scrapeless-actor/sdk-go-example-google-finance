package region

import (
	"github.com/PuerkitoBio/goquery"
)

type NewsResultsInfo struct {
	Snippet   string `json:"snippet"`
	Link      string `json:"link"`
	Source    string `json:"source"`
	Date      string `json:"date"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func NewsResults(dc *goquery.Document) []NewsResultsInfo {
	var (
		resp = make([]NewsResultsInfo, 0)
	)
	dc.Find("div[jscontroller='DrJTUc']").
		Find("div[jscontroller='ZpnVYd'] div[class='nkXTJ']").
		Each(func(i int, selection *goquery.Selection) {
			source := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Text()
			date := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(1).Text()
			snippet := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(1).Text()
			link, _ := selection.Children().Eq(1).Children().Eq(0).Attr("href")
			thumbnail, _ := selection.Children().Eq(1).Children().Eq(0).Children().Eq(0).Attr("src")
			resp = append(resp, NewsResultsInfo{
				Snippet:   snippet,
				Link:      link,
				Source:    source,
				Date:      date,
				Thumbnail: thumbnail,
			})
		})
	return resp
}

type TopNewsResultsInfo struct {
	Title string            `json:"title"`
	Items []NewsResultsInfo `json:"items"`
}

func NewsResultsMarkets(dc *goquery.Document) []NewsResultsInfo {
	var (
		resp = make([]NewsResultsInfo, 0)
	)
	dc.
		Find("div[jscontroller='ZpnVYd'] div[class='nkXTJ']").
		Each(func(i int, selection *goquery.Selection) {
			source := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Text()
			date := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(1).Text()
			snippet := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(1).Text()
			link, _ := selection.Children().Eq(1).Children().Eq(0).Attr("href")
			thumbnail, _ := selection.Children().Eq(1).Children().Eq(0).Children().Eq(0).Attr("src")
			resp = append(resp, NewsResultsInfo{
				Snippet:   snippet,
				Link:      link,
				Source:    source,
				Date:      date,
				Thumbnail: thumbnail,
			})
		})
	return resp
}

func TopNewsResults(dc *goquery.Document) []TopNewsResultsInfo {
	var (
		topNewsResultsInfos []TopNewsResultsInfo
	)
	if dc.Find("div[jscontroller='DrJTUc']").Find("div[class='qQfHId']").Length() == 0 {
		return topNewsResultsInfos
	}
	selection := dc.Find("div[jscontroller='DrJTUc']").Find("div[class='qQfHId']").First()
	title := selection.Children().Eq(0).Text()
	var topNewsResultsInfo TopNewsResultsInfo
	topNewsResultsInfo.Title = title
	selection.Children().Eq(1).
		Find("div[jscontroller='ZpnVYd']").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Children().Eq(0).Attr("href")
		thumbnail, _ := selection.Children().Eq(0).Children().Eq(0).Find("img").Attr("src")
		source := selection.Children().Eq(0).Children().Eq(1).Text()
		snippet := selection.Children().Eq(0).Children().Eq(2).Text()
		date := selection.Children().Eq(0).Children().Eq(3).Text()
		newsResultsInfo := NewsResultsInfo{
			Snippet:   snippet,
			Link:      link,
			Source:    source,
			Date:      date,
			Thumbnail: thumbnail,
		}
		topNewsResultsInfo.Items = append(topNewsResultsInfo.Items, newsResultsInfo)
	})
	topNewsResultsInfos = append(topNewsResultsInfos, topNewsResultsInfo)
	return topNewsResultsInfos
}

func OpinionAndNumbersResults(dc *goquery.Document) []TopNewsResultsInfo {
	var (
		topNewsResultsInfos []TopNewsResultsInfo
	)
	if dc.Find("div[jscontroller='DrJTUc']").Find("div[class='qQfHId']").Length() == 0 {
		return topNewsResultsInfos
	}
	dc.Find("div[jscontroller='DrJTUc']").Find("div[class='qQfHId']").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}
		title := selection.Children().Eq(0).Text()
		var topNewsResultsInfo TopNewsResultsInfo
		topNewsResultsInfo.Title = title
		selection.Children().Eq(1).
			Find("div[jscontroller='ZpnVYd']").Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Children().Eq(0).Attr("href")
			source := selection.Children().Eq(0).Children().Eq(0).Text()
			snippet := selection.Children().Eq(0).Children().Eq(1).Text()
			date := selection.Children().Eq(0).Children().Eq(2).Text()
			newsResultsInfo := NewsResultsInfo{
				Snippet: snippet,
				Link:    link,
				Source:  source,
				Date:    date,
			}
			topNewsResultsInfo.Items = append(topNewsResultsInfo.Items, newsResultsInfo)
		})
		topNewsResultsInfos = append(topNewsResultsInfos, topNewsResultsInfo)
	})

	return topNewsResultsInfos
}
