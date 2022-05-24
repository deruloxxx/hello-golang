package main

import (
	"encoding/json"
	"hello-golang/meander"
	"net/http"
	"runtime"
)

func main() {
	// プログラムから使用できるCPU数の最大値を指定
	runtime.GOMAXPROCS(runtime.NumCPU())

	// meander.APIKeyをセット
	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		// meander.Journeysをエンコード化してwに書き出し
		respond(w, r, meander.Journeys)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		// データのスライスに含まれるそれぞれの要素に対してmeander.Publicを呼び出す
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}
