package main

import (
	"encoding/json"
	"fmt"
	"hello-golang/meander"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	// プログラムから使用できるCPU数の最大値を指定
	runtime.GOMAXPROCS(runtime.NumCPU())

	// meander.APIKeyをセット
	meander.APIKey = os.Getenv("GOOGLE_PLACES_API_KEY")
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
