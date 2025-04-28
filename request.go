package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sdk-go-example-google-finance/region"
	"strings"
	"time"
)

type FinanceInfo struct {
	Markets        any                     `json:"markets"`
	Summary        any                     `json:"summary"`
	Graph          any                     `json:"graph"`
	KnowledgeGraph any                     `json:"knowledge_graph"`
	NewsResults    []any                   `json:"news_results"`
	Financials     []region.FinancialsInfo `json:"financials,omitempty"`
	DiscoverMore   any                     `json:"discover_more"`
}

func GetFinance(ctx context.Context, hl string, q string) (*FinanceInfo, error) {
	var (
		financeInfo FinanceInfo
	)
	// get html data
	res, err := getData(ctx, hl, q)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return nil, err
	}

	mapping := GetDataMapping(res)

	// DiscoverAndPeopleSearch
	discoverAndPeopleSearch := region.DiscoverAndPeopleSearch(doc)
	financeInfo.DiscoverMore = discoverAndPeopleSearch
	// financials
	changeAll := region.GetYYChangeAll(mapping["13"])
	financeInfo.Financials = changeAll

	knowledgeGraphKeyStats := region.KnowledgeGraphKeyStats(doc)
	knowledgeGraph := region.KnowledgeGraph{
		KeyStats: knowledgeGraphKeyStats,
		About:    []any{region.About(doc)},
	}
	financeInfo.KnowledgeGraph = knowledgeGraph

	financeInfo.Graph = region.Graph(mapping["10"])
	for _, info := range region.TopNewsResults(doc) {
		financeInfo.NewsResults = append(financeInfo.NewsResults, info)
	}
	for _, info := range region.OpinionAndNumbersResults(doc) {
		financeInfo.NewsResults = append(financeInfo.NewsResults, info)
	}
	for _, info := range region.NewsResults(doc) {
		financeInfo.NewsResults = append(financeInfo.NewsResults, info)
	}
	if len(changeAll) == 0 {
		financeInfo.Summary = region.Summary(doc)

	} else {
		summaryInfo := region.SummaryWithMarket(doc, mapping["17"])
		if region.SummaryMarkets(doc) != nil {
			summaryInfo.Market = region.SummaryMarkets(doc)
		}
		summaryInfo.Extensions = region.SummaryExtensions(doc)
		financeInfo.Summary = summaryInfo
	}
	marketsInfo := region.Markets(mapping["16"], region.GetMarketTitle(doc))
	financeInfo.Markets = marketsInfo

	return &financeInfo, nil
}

func getData(ctx context.Context, hl string, q string) (string, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://www.google.com/finance/quote/%s?hl=%s", q, hl), nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/jpeg,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Priority", "u=0, i")
	request.Header.Set("Sec-Ch-Ua", `"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`)
	request.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	request.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Site", "none")
	request.Header.Set("Sec-Fetch-User", "?1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	do, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	body, err := io.ReadAll(do.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetDataMapping(data string) map[string]string {
	var (
		dataMapping = make(map[string]string)
	)
	re := regexp.MustCompile(`key:.*?sideChannel?`)
	re1 := regexp.MustCompile(`key: '([^']+)'`)
	re2 := regexp.MustCompile(`data:.*?sideChannel?`)
	matchString := re.FindAllString(data, -1)
	for _, v := range matchString {
		key := re1.FindAllString(v, -1)
		a := strings.Replace(key[0], `'`, "", -1)
		split := strings.Split(a, ":")
		dataRes := re2.FindAllString(v, -1)
		s := strings.Replace(dataRes[0], "data:", "", -1)
		s = strings.Replace(s, ", sideChannel", "", -1)
		dataMapping[split[len(split)-1]] = s
	}
	return dataMapping
}

type GraphInfo struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency,omitempty"`
	Date     string  `json:"date"`
	Volume   int     `json:"volume,omitempty"`
}

func GetFinanceByWindow(ctx context.Context, q string, hl string, index int) (any, error) {
	var (
		graphInfo []GraphInfo
	)
	parse, err := url.Parse("https://www.google.com/finance/_/GoogleFinanceUi/data/batchexecute")
	if err != nil {
		return nil, err
	}
	queryValues := &url.Values{}
	queryForm := &url.Values{}
	queryForm.Set("f.req", getForm(q, index))
	queryValues.Add("source-path", fmt.Sprintf("/finance/quote/%s", q))
	queryValues.Add("bl", "boq_finance-ui_20250306.11_p0")
	queryValues.Add("hl", hl)
	queryValues.Add("rt", "c")
	parse.RawQuery = queryValues.Encode()
	request, err := http.NewRequest(http.MethodPost, parse.String(), strings.NewReader(queryForm.Encode()))
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}
	request.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("accept", "*/*")
	request.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	request.Header.Set("referer", "https://www.google.com/")
	request.Header.Set("content-length", fmt.Sprintf("%d", len(queryForm.Encode())))
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	request.Header.Set("Cookie", "NID=522=BzbKuspW0ZDXCQ9hE0RS5R32UWVV-5zduP7MOeRC9dxZoATZsOZq6ILwgBUzmoARDWA2tG8f1Flxyd5HP2nBEwFLWiayizCaV5IlTJR2T2TGVAGAuzqLunPVOkCKwZzG0hM4KyRqam2WmKmgwOQZx37pGpSocl0cp7x6xY4VDp-ojnFOv89KMVkC9Zg; GN_PREF=W251bGwsIkNBSVNEQWlNLWZXOUJoRG80S2pIQVEiXQ__")
	do, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	wrbs := GetWrbs(all)
	result := gjson.Parse(string(wrbs[0]))
	currency := gjson.Parse(result.Get("0.2").String()).Get("0.0.2").String()

	for _, r := range gjson.Parse(result.Get("0.2").String()).Get("0.0.3.0.1").Array() {
		year := r.Get("0.0").Int()
		month := r.Get("0.1").Int()
		day := r.Get("0.2").Int()
		hour := r.Get("0.3").Int()
		date := time.Date(int(year), time.Month(month), int(day), int(hour), 0, 0, 0, time.UTC)
		price := r.Get("1.0").Float()
		value := r.Get("2").Int()
		graphInfo = append(graphInfo, GraphInfo{
			Price:    price,
			Currency: currency,
			Date:     date.String(),
			Volume:   int(value),
		})
		// Sep 11 2024, 04:00 PM UTC-04:00
	}
	return graphInfo, nil
}

func getForm(q string, index int) string {
	array1 := strings.Split(q, ":")
	array2 := []any{nil, array1}
	array3 := []any{array2}
	array4 := []any{array3, index, nil, nil, nil, nil, nil, 0}
	array4Str, _ := json.Marshal(array4)
	array5 := []any{"AiCwsd", string(array4Str), nil, "generic"}
	array6 := []any{array5}
	array7 := []any{array6}
	marshal, _ := json.Marshal(array7)
	return string(marshal)
}

func GetWrbs(respBytes []byte) [][]byte {
	lines := bytes.Split(respBytes, []byte("\n"))
	var wrbs [][]byte
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte(`[["wrb.fr"`)) {
			wrbs = append(wrbs, line)
		}
	}
	return wrbs
}
