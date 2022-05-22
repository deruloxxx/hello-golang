package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDBにダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closeDB() {
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

func main() {}
