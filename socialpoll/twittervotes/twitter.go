package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type tweet struct {
	Text string
}

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	// 接続を示すconnが閉じられているかを確認
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	// connの値を更新(最新のデータを取得する)
	conn = netc
	return netc, nil
}

// twitterとの接続をclean upする
var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

// リクエストの認証に使用するOAuthオブジェクトをセットアップ
var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

func setupTwitterAuth() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	var ts struct {
		Consumerkey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken    string `env:"SP_TWITTER_AccessToken,required"`
		AccessSecret   string `env:"SP_TWITTER_AccessSecret,required"`
	}
	// 環境変数のrequiredを取得できなかったらerrorを出す
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.Consumerkey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	return httpClient.Do(req)
}

// Twitter上で投票が行われたことを通知
func readFromTwitter(votes chan<- string) {
	// loadOptions関数で全ての投票での選択肢を取得
	options, err := loadOptions()
	if err != nil {
		log.Println("選択肢の読み込みに失敗しました:", err)
		return
	}
	// Twitter側のエンドポイントを指定
	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("URLの解析に失敗しました:", err)
		return
	}
	query := make(url.Values)
	// 選択肢のリクエストをカンマ区切りで指定
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("検索リクエストの作成に失敗しました:", err)
		return
	}
	resp, err := makeRequest(req, query)
	if err != nil {
		log.Println("検索のリクエストに失敗しました:", err)
		return
	}
	reader = resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var tweet tweet
		if err := decoder.Decode(&tweet); err != nil {
			break
		}
		for _, option := range options {
			if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(option)) {
				log.Println("投票:", option)
				votes <- option
			}
		}
	}
}

// chan struct{}は受信専用 votesは投票内容が送信
func startTwitterStream(stopchan <-chan struct{},
	votes chan<- string) <-chan struct{} {
	stoppedchan := make(chan struct{}, 1)
	go func() {
		defer func() {
			// ここで受信したら処理を終了してリターンさせる
			stoppedchan <- struct{}{}
		}()
		for {
			select {
			case <-stopchan:
				log.Println("Twitterへの問い合わせを終了します...")
				return
			default:
				log.Println("Twitterに問い合わせます...")
				readFromTwitter(votes)
				log.Println("(待機中)")
				time.Sleep(10 * time.Second) // 待機してから再接続します
			}
		}
	}()
	// goチャンネルの終了を伝える
	return stoppedchan
}
