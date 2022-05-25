package meander

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Google Places APIからのレスポンスを使いやすいオブジェクトに変換するための定義
type Place struct {
	*googleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}

var APIKey string

type googleResponse struct {
	Results []*Place `json:"results"`
}

type googleGeometry struct {
	*googleLocation `json:"location"`
}

type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type googlePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}

type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"icon":     p.Icon,
		"photos":   p.Photos,
		"vicinity": p.Vicinity,
		"lat":      p.Lat,
		"lng":      p.Lng,
	}
}

func (q *Query) find(types string) (*googleResponse, error) {
	u := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	// リクエストURLを作成
	vals := make(url.Values)
	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("types", types)
	vals.Set("key", APIKey)

	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		vals.Set("minprice", fmt.Sprintf("%d", int(r.From)-1))
		vals.Set("maxprice", fmt.Sprintf("%d", int(r.To)-1))
	}
	res, err := http.Get(u + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var responce googleResponse
	if err := json.NewDecoder(res.Body).Decode(&responce); err != nil {
		return nil, err
	}
	return &responce, nil
}

// 問い合わせを一斉に行い、その結果を返します
func (q *Query) Run() []interface{} {
	// ナノ秒単位で時間を取得
	rand.Seed(time.Now().UnixNano())
	var w sync.WaitGroup
	var l sync.Mutex
	places := make([]interface{}, len(q.Journey))
	// Query.Findメソッドを並行に呼び出す。早くリクエストを送信するため。
	for i, r := range q.Journey {
		w.Add(1)
		go func(types string, i int) {
			// WaitGroupオブジェクトに対してリクエストの完了を伝える
			defer w.Done()
			// リクエストを実行
			response, err := q.find(types)
			if err != nil {
				fmt.Errorf("施設の検索に失敗しました:", err)
				return
			}
			if len(response.Results) == 0 {
				fmt.Errorf("施設が見つかりませんでした", types)
				return
			}
			for _, result := range response.Results {
				for _, photo := range result.Photos {
					// クライアントにAPIを意識させないため写真のURLをこちらで用意
					photo.URL = "https://maps.googleapis.com/maps/api/place/photo?" +
						"maxwidth=1000&photoreference=" + photo.PhotoRef +
						"&key=" + APIKey
				}
			}
			randI := rand.Intn(len(response.Results))
			l.Lock()
			places[i] = response.Results[randI]
			l.Unlock()
		}(r, i)
	}
	w.Wait() // 全てのリクエストの完了を待ちます
	return places
}
