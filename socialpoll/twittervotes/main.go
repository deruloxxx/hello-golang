package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDBにダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("データベース接続が閉じられました")
}

type poll struct {
	Options []string
}

// Twitter検索に使われる選択肢を取り出す
func loadOptions() ([]string, error) {
	var options []string
	// ballotsデータベースに含まれるコレクションpollsを取り出す。nilはフィルタリングを行わないという意味
	// Find(nil)はフィルタリングを行わないということ
	// 流れるようなインターフェースに基づく(メソッド呼び出しの連鎖)
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func main() {
	// 穏やかな起動と終了
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("停止します...")
		stopChan <- struct{}{}
		closeConn()
	}()
	// プログラムを終了させようとしたらシグナルを送信
	// syscall.SIGTERMがUnixシグナル
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("MongoDBへのダイヤルに失敗しました:", err)
	}
	defer closedb()

	// 処理を開始します
	votes := make(chan string) // 投票結果のためのチャネル
	publihserStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publihserStoppedChan
}
