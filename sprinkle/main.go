package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{
	otherWord,
	otherWord,
	otherWord,
	otherWord + "app",
	otherWord + "site",
	"get" + otherWord,
	"go" + otherWord,
	"lets " + otherWord,
}

func main() {
	// 乱数をの元になるシード値から乱数を作成
	rand.Seed(time.Now().UTC().UnixNano())
	// 標準入力(入力装置,OSが提供するデータ入力機能)を取得
	s := bufio.NewScanner(os.Stdin)
	// データの入力元をscanメソッドで読み込み。データがあればtrueを返して実行
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
