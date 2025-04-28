package region

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

type SummaryInfo struct {
	Title          string              `json:"title"`
	Stock          string              `json:"stock"`
	Exchange       string              `json:"exchange"`
	Price          string              `json:"price"`
	ExtractedPrice string              `json:"extracted_price"`
	Currency       string              `json:"currency"`
	PriceMovement  PriceMovement       `json:"price_movement"`
	Market         *SummaryMarketsInfo `json:"market,omitempty"`
	Extensions     any                 `json:"extensions,omitempty"`
}

type SummaryInfoWithExtensions struct {
	Title          string  `json:"title"`
	Stock          string  `json:"stock"`
	Exchange       string  `json:"exchange"`
	Price          string  `json:"price"`
	ExtractedPrice float64 `json:"extracted_price"`
	Extensions     any     `json:"extensions"`
}

func Summary(dc *goquery.Document) SummaryInfoWithExtensions {
	var (
		summaryInfoWithExtensions SummaryInfoWithExtensions
	)
	summaryInfoWithExtensions.Price = dc.Find("div[jscontroller='NdbN0c']").Children().Eq(0).Children().Eq(0).Children().Eq(0).Text()
	extractedPriceStr := strings.Replace(summaryInfoWithExtensions.Price, ",", "", -1)
	extractedPrice, _ := strconv.ParseFloat(extractedPriceStr, 64)
	summaryInfoWithExtensions.ExtractedPrice = extractedPrice
	dc.Find("div[jscontroller='NdbN0c']").Children().Eq(1).Contents().Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) == "#text" {
			split := strings.Split(selection.Text(), " · ")
			summaryInfoWithExtensions.Extensions = split[:len(split)-1]
		}
	})
	dc.Find("div[jscontroller='DrJTUc'] div[class='xJwwl']").Children().Eq(0).Contents().Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) == "#text" {
			split := strings.Split(selection.Text(), " • ")
			summaryInfoWithExtensions.Stock = split[0]
			summaryInfoWithExtensions.Exchange = split[1]
		}
	})
	title := dc.Find("div[jscontroller='DrJTUc'] div[class='xJwwl']").Children().Eq(1).Text()
	summaryInfoWithExtensions.Title = title
	return summaryInfoWithExtensions
}

func SummaryWithMarket(dc *goquery.Document, data string) SummaryInfo {
	price := dc.Find("div[jscontroller='NdbN0c']").Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Eq(0).Children().Text()
	currency := price[:1]
	extractedPrice := price[1:]
	d := gjson.Parse(data).Get("0.0.0")
	stock := d.Get("1.0").String()
	exchange := d.Get("1.1").String()
	title := d.Get("2").String()
	value := d.Get("5.1").Float()
	percentage := d.Get("5.2").Float()
	movement := "Up"
	if value < 0 {
		movement = "Down"
	}
	return SummaryInfo{
		Title:          title,
		Stock:          stock,
		Exchange:       exchange,
		Price:          price,
		ExtractedPrice: extractedPrice,
		Currency:       currency,
		PriceMovement: PriceMovement{
			Value:      value,
			Percentage: percentage,
			Movement:   movement,
		},
	}
}

type SummaryMarketsInfo struct {
	Trading        string        `json:"trading"`
	Price          string        `json:"price"`
	ExtractedPrice string        `json:"extracted_price"`
	Currency       string        `json:"currency"`
	PriceMovement  PriceMovement `json:"price_movement"`
}

func SummaryMarkets(dc *goquery.Document) *SummaryMarketsInfo {
	if _, ok := dc.Find("div[jscontroller='NdbN0c']").Find("div[class='ivZBbf ygUjEc']").Attr("class"); !ok {
		return nil
	}
	data := dc.Find("div[jscontroller='NdbN0c']").Children().Eq(1).Text()
	split := strings.Split(data, ":")
	split1 := strings.Split(split[1], "(")
	price := split1[0]
	extractedPrice := price[1:]
	currency := price[:1]
	split2 := strings.Split(split1[1], ")")
	movement := "Down"
	value, _ := strconv.ParseFloat(split2[1][1:], 64)
	percentage, _ := strconv.ParseFloat(split2[0][:len(split2[0])-1], 64)
	if strings.Contains(split2[1], "+") {
		movement = "Up"
	}
	return &SummaryMarketsInfo{
		Trading:        split[0],
		Price:          price,
		ExtractedPrice: extractedPrice,
		Currency:       currency,
		PriceMovement: PriceMovement{
			Value:      value,
			Percentage: percentage,
			Movement:   movement,
		},
	}
}

func SummaryExtensions(dc *goquery.Document) []string {
	var data []string
	split := strings.Split(dc.Find("div[jscontroller='NdbN0c']").Find("div[class='ygUjEc']").Text(), " · ")
	for i := 0; i < len(split)-1; i++ {
		data = append(data, split[i])
	}
	return data
}
