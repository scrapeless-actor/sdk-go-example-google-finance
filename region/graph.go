package region

import (
	"github.com/tidwall/gjson"
	"time"
)

type GraphInfo struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency,omitempty"`
	Date     string  `json:"date"`
	Volume   int     `json:"volume,omitempty"`
}

func Graph(data string) any {
	var (
		currency  = gjson.Parse(data).Get("0.0.2").String()
		graphInfo []GraphInfo
	)

	for _, v := range gjson.Parse(data).Get("0.0.3.0.1").Array() {
		year := v.Get("0.0").Int()
		month := v.Get("0.1").Int()
		day := v.Get("0.2").Int()
		hour := v.Get("0.3").Int()
		minute := v.Get("0.4").Int()
		price := v.Get("1.0").Float()
		volume := v.Get("2").Int()
		date := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), 0, 0, time.UTC)
		graphInfo = append(graphInfo, GraphInfo{
			Price:    price,
			Currency: currency,
			Date:     date.String(),
			Volume:   int(volume),
		})
	}
	return graphInfo
}
