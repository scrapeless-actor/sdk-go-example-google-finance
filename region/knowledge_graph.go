package region

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type KnowledgeGraph struct {
	KeyStats any `json:"key_stats"`
	About    any `json:"about"`
}

type KeyStatsTag struct {
	Text        string `json:"text"`
	Description string `json:"description"`
}
type KeyStatsStats struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

func KnowledgeGraphKeyStats(dc *goquery.Document) map[string]any {
	var (
		keyStatsTag   []KeyStatsTag
		keyStatsStats []KeyStatsStats
	)
	//fmt.Println(dc.Find("div[class='T4LgNb'] div[role='region']").Attr("class"))
	dc.Find("div[jscontroller='DrJTUc'] div[class='eYanAe']").Children().Each(func(i int, s *goquery.Selection) {
		// 获取href属性
		val, _ := s.Attr("class")
		if val == "vvDK2c" { //tags
			s.Find("div[class='UaHgge']").Children().Each(func(i int, selection *goquery.Selection) {
				description := selection.Find("div").Text()
				text := selection.Find("a span span").Text()
				keyStatsTag = append(keyStatsTag, KeyStatsTag{
					Text:        text,
					Description: description,
				})
			})
		}
		if val == "gyFHrc" { //stats
			label := s.Find("span").Children().Eq(0).Text()
			description := s.Find("span").Children().Eq(1).Text()
			value := s.Children().Eq(1).Text()
			keyStatsStats = append(keyStatsStats, KeyStatsStats{
				Label:       label,
				Description: description,
				Value:       value,
			})
		}
	})
	return map[string]any{
		"tags":  keyStatsTag,
		"stats": keyStatsStats,
	}
}

type AboutInfo struct {
	Title       string           `json:"title"`
	Description AboutDescription `json:"description"`
	Info        []AboutInfoInfo  `json:"info,omitempty"`
}

type AboutDescription struct {
	Snippet  string `json:"snippet"`
	Link     string `json:"link"`
	LinkText string `json:"link_text"`
}

type AboutInfoInfo struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Link  string `json:"link,omitempty"`
}

func About(dc *goquery.Document) any {
	var (
		aboutInfo     AboutInfo
		aboutInfoInfo []AboutInfoInfo
	)
	dc.Find("div[jscontroller='BPOkb']").Each(func(i int, selection *goquery.Selection) {
		if selection.Children().Eq(0).Children().Eq(0).Text() == "About" {
			descriptionSnippet := selection.Children().Eq(1).Children().Eq(0).Text()
			descriptionLink, _ := selection.Children().Eq(1).Children().Eq(0).Find("a").Attr("href")
			descriptionText := selection.Children().Eq(1).Children().Eq(0).Find("a").Text()
			descriptionSnippet = strings.Replace(descriptionSnippet, descriptionText, "", -1)
			selection.Children().Eq(1).Find("div[class='gyFHrc']").Each(func(i int, s1 *goquery.Selection) {
				infoLabel := s1.Children().Eq(0).Text()
				infoValue := s1.Children().Eq(1).Text()
				infoLink, _ := s1.Children().Eq(1).Find("a").Attr("href")
				aboutInfoInfo = append(aboutInfoInfo, AboutInfoInfo{
					Label: infoLabel,
					Value: infoValue,
					Link:  infoLink,
				})
			})
			aboutInfo.Title = "About"
			aboutInfo.Description.Snippet = descriptionSnippet
			aboutInfo.Description.Link = descriptionLink
			aboutInfo.Description.LinkText = descriptionText
			aboutInfo.Info = aboutInfoInfo
		}
	})
	return aboutInfo
}
