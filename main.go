package main

import (
	"context"
	"encoding/json"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless"
	proxyModel "github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/proxy"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var (
	client *http.Client
)

type RequestParam struct {
	Q  string `json:"q" url:"q"`
	Hl string `json:"hl" url:"hl"`

	Window string `json:"window" url:"window"`
}

var (
	windowMapping = map[string]int{
		"1D":  1,
		"5D":  2,
		"1M":  3,
		"6M":  4,
		"YTD": 5,
		"1Y":  6,
		"5Y":  7,
		"MAX": 8,
	}
)

func main() {
	// new actor
	actor := scrapeless.New(scrapeless.WithProxy(), scrapeless.WithStorage())
	defer actor.Close()
	var param = &RequestParam{}
	if err := actor.Input(param); err != nil {
		log.Fatal(err)
	}
	// get proxy url
	proxy, err := actor.Proxy.Proxy(context.TODO(), proxyModel.ProxyActor{
		Country:         "us",
		SessionDuration: 10,
	})

	if err != nil {
		panic(err)
	}
	parse, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	// init client with proxy
	client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(parse)}}

	// do crawl logic
	data, err := GetFinance(context.TODO(), param.Hl, param.Q)
	if param.Window != "" {
		if _, ok := windowMapping[param.Window]; !ok {
			log.Errorf("success=false,  err=%v", err)
			return
		}
		index := windowMapping[param.Window]
		data.Graph, err = GetFinanceByWindow(context.TODO(), param.Q, param.Hl, index)
		if err != nil {
			log.Errorf("success=false,  err=%v", err)
			return
		}
	}
	if err != nil {
		log.Errorf("success=false,  err=%v", err)
		return
	}
	resultBytes, _ := json.Marshal(data)
	ok, err := actor.Storage.GetKv().SetValue(context.Background(), "data", string(resultBytes), 0)
	if !ok || err != nil {
		log.Errorf("set kv failed,  err=%v", err)
		return
	}
	log.Info("set kv success")
}
