package region

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"strings"
)

type MarketsInfo struct {
	Stock         string        `json:"stock"`
	Link          string        `json:"link"`
	Name          string        `json:"name"`
	Price         float64       `json:"price"`
	PriceMovement PriceMovement `json:"price_movement"`
}

func Markets(data string, marketTitle []string) map[string]any {
	var (
		resp  = make(map[string]any)
		count = 0
		index = 0
	)
	for _, result := range gjson.Parse(data).Get("0").Array() {
		for _, k := range result.Array() {
			for _, k1 := range k.Array() {
				if k1.IsArray() {
					if count%5 == 0 && count != 0 {
						index++
					}
					info := parseMarkets(k1.String())
					if resp[marketTitle[index]] == nil {
						resp[marketTitle[index]] = make([]any, 0)
					}
					resp[marketTitle[index]] = append(resp[marketTitle[index]].([]any), info)
					count++
				}
			}
		}
	}
	return resp
}

func parseMarkets(data string) (marketsInfo MarketsInfo) {
	d := gjson.Parse(data)
	name := d.Get("2").String()
	if name == "" {
		name = d.Get("1.0.2").String()
	}
	info := d.Get("1.0")
	stock := info.Get("21").String()
	link := fmt.Sprintf("https://www.google.com/finance/quote/%s", stock)
	price := info.Get("5.0").Float()
	priceMovementValue := info.Get("5.1").Float()
	priceMovementPercentage := info.Get("5.2").Float()
	priceMovementMovement := "Up"
	if priceMovementPercentage < 0 {
		priceMovementMovement = "Down"
	}
	return MarketsInfo{
		Stock: stock,
		Link:  link,
		Name:  name,
		Price: price,
		PriceMovement: PriceMovement{
			Value:      priceMovementValue,
			Percentage: priceMovementPercentage,
			Movement:   priceMovementMovement,
		},
	}
}

func GetMarketTitle(dc *goquery.Document) (marketTitle []string) {
	dc.Find("div[jscontroller='LMhoGc']").Children().Eq(1).Children().Eq(0).Children().Children().Each(func(i int, selection *goquery.Selection) {
		val, _ := selection.Attr("class")
		if strings.Contains(val, "AHyjFe") {
			marketTitle = append(marketTitle, selection.Text())
		}
	})
	return
}

type TopNewsInfo struct {
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
	Source  string `json:"source"`
	Date    string `json:"date"`
}

// TopNews
func TopNews(dc *goquery.Document) TopNewsInfo {
	link, _ := dc.Find("div[jscontroller='LMhoGc']").Children().Eq(1).Children().Eq(1).Children().Attr("href")
	snippet := dc.Find("div[jscontroller='LMhoGc']").Children().Eq(1).Children().Eq(1).Children().Children().Eq(0).Text()
	source := dc.Find("div[jscontroller='LMhoGc']").Children().Eq(1).Children().Eq(1).Children().Children().Eq(1).Text()
	date := dc.Find("div[jscontroller='LMhoGc']").Children().Eq(1).Children().Eq(1).Children().Children().Eq(2).Text()
	source = getResult(source)
	date = getResult(date)
	return TopNewsInfo{
		Link:    link,
		Snippet: snippet,
		Source:  source,
		Date:    date,
	}
}

func getResult(data string) string {
	split := strings.Split(data, "•")
	if len(split) < 2 {
		return ""
	}
	return strings.TrimSpace(split[1])
}
