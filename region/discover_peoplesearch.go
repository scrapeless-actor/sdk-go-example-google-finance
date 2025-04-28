package region

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"unicode"
)

type DiscoverAndPeopleSearchInfo struct {
	Title string                        `json:"title"`
	Items []DiscoverAndPeopleSearchItem `json:"items"`
}

type DiscoverAndPeopleSearchItem struct {
	Stock          string  `json:"stock"`
	Link           string  `json:"link"`
	Name           string  `json:"name"`
	Price          string  `json:"price"`
	ExtractedPrice float64 `json:"extracted_price"`
	Currency       string  `json:"currency,omitempty"`
	PriceMovement  PriceMovement
}
type PriceMovement struct {
	Value      float64 `json:"value,omitempty"`
	Percentage float64 `json:"percentage"`
	Movement   string  `json:"movement"`
}

func DiscoverAndPeopleSearch(dc *goquery.Document) any {
	var (
		discoverAndPeopleSearchInfo []DiscoverAndPeopleSearchInfo
		interestedIn                DiscoverAndPeopleSearchInfo
		search                      DiscoverAndPeopleSearchInfo
	)
	// You may be interested in
	dc.Find("section[role='complementary']").Children().Eq(1).Children().Eq(0).Contents().Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) == "#text" {
			title := strings.TrimSpace(selection.Text())
			interestedIn.Title = title
		}
	})
	dc.Find("section[role='complementary']").Children().Eq(1).Children().Eq(1).Children().Eq(0).Children().Eq(0).Find("div[role='listitem']").Each(func(i int, selection *goquery.Selection) {
		a := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0)
		stock, _ := a.Attr("href")
		split := strings.Split(stock, `/`)
		stock = split[len(split)-1]
		link := fmt.Sprintf("https://www.google.com/finance/quote/%s", stock)
		name := a.Contents().Eq(1).Text()
		percentage := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Contents().Eq(1).Text()
		movement, _ := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Find("path").Attr("d")
		if strings.Contains(movement, "M20") {
			movement = "Down"
		} else {
			movement = "Up"
		}
		price := a.Children().Eq(2).Children().Eq(0).Text()
		currency := ""
		if !isFirstDigit(price) {
			currency = price[:1]
			price = price[1:]
		}
		replacePrice := strings.Replace(price, ",", "", -1)
		priceF, _ := strconv.ParseFloat(replacePrice, 64)
		percentageF, _ := strconv.ParseFloat(strings.Replace(percentage, "%", "", -1), 64)
		interestedIn.Items = append(interestedIn.Items, DiscoverAndPeopleSearchItem{
			Stock:          stock,
			Link:           link,
			Name:           name,
			Price:          price,
			ExtractedPrice: priceF,
			Currency:       currency,
			PriceMovement: PriceMovement{
				Movement:   movement,
				Percentage: percentageF,
			},
		})
	})

	// People also search for
	search.Title = dc.Find("section[role='complementary']").Children().Eq(2).Children().Eq(0).Text()
	dc.Find("section[role='complementary']").Children().Eq(2).Children().Eq(1).Children().Eq(0).Children().Eq(0).Find("div[role='listitem']").Each(func(i int, selection *goquery.Selection) {
		a := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0)
		stock, _ := a.Attr("href")
		split := strings.Split(stock, `/`)
		stock = split[len(split)-1]
		link := fmt.Sprintf("https://www.google.com/finance/quote/%s", stock)
		name := a.Contents().Eq(1).Text()
		percentage := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Contents().Eq(1).Text()
		movement, _ := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Find("path").Attr("d")
		if strings.Contains(movement, "M20") {
			movement = "Down"
		} else {
			movement = "Up"
		}
		price := a.Children().Eq(2).Children().Eq(0).Text()
		currency := ""
		if !isFirstDigit(price) {
			currency = price[:1]
			price = price[1:]
		}
		priceF, _ := strconv.ParseFloat(price, 64)
		percentageF, _ := strconv.ParseFloat(strings.Replace(percentage, "%", "", -1), 64)
		search.Items = append(interestedIn.Items, DiscoverAndPeopleSearchItem{
			Stock:          stock,
			Link:           link,
			Name:           name,
			Price:          price,
			ExtractedPrice: priceF,
			Currency:       currency,
			PriceMovement: PriceMovement{
				Movement:   movement,
				Percentage: percentageF,
			},
		})
	})

	discoverAndPeopleSearchInfo = append(discoverAndPeopleSearchInfo, interestedIn, search)
	return discoverAndPeopleSearchInfo
}

func isFirstDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsDigit(rune(s[0]))
}

func DiscoverMarkets(dc *goquery.Document) any {
	var (
		discoverAndPeopleSearchInfo []DiscoverAndPeopleSearchInfo
		interestedIn                DiscoverAndPeopleSearchInfo
	)
	// You may be interested in
	dc.Find("section[role='complementary']").Children().Eq(1).Children().Eq(0).Contents().Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) == "#text" {
			title := strings.TrimSpace(selection.Text())
			interestedIn.Title = title
		}
	})

	dc.Find("section[role='complementary']").Children().Eq(1).Children().Eq(1).Children().Eq(0).Children().Eq(0).Find("div[role='listitem']").Each(func(i int, selection *goquery.Selection) {
		a := selection.Children().Eq(0).Children().Eq(0).Children().Eq(0)
		stock, _ := a.Attr("href")
		split := strings.Split(stock, `/`)
		stock = split[len(split)-1]
		link := fmt.Sprintf("https://www.google.com/finance/quote/%s", stock)
		name := a.Contents().Eq(1).Text()
		percentage := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Contents().Eq(1).Text()
		movement, _ := a.Children().Eq(2).Find("span[jsname='Fe7oBc']").Children().Eq(0).Children().Eq(0).Find("path").Attr("d")
		if strings.Contains(movement, "M20") {
			movement = "Down"
		} else {
			movement = "Up"
		}
		price := a.Children().Eq(2).Children().Eq(0).Text()
		currency := ""
		if !isFirstDigit(price) {
			currency = price[:1]
			price = price[1:]
		}
		replacePrice := strings.Replace(price, ",", "", -1)
		priceF, _ := strconv.ParseFloat(replacePrice, 64)
		percentageF, _ := strconv.ParseFloat(strings.Replace(percentage, "%", "", -1), 64)
		interestedIn.Items = append(interestedIn.Items, DiscoverAndPeopleSearchItem{
			Stock:          stock,
			Link:           link,
			Name:           name,
			Price:          price,
			ExtractedPrice: priceF,
			Currency:       currency,
			PriceMovement: PriceMovement{
				Movement:   movement,
				Percentage: percentageF,
			},
		})
	})
	discoverAndPeopleSearchInfo = append(discoverAndPeopleSearchInfo, interestedIn)
	return discoverAndPeopleSearchInfo
}
