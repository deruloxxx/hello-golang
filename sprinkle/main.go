package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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

func readLine(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	// 乱数をの元になるシード値から乱数を作成
	rand.Seed(time.Now().UTC().UnixNano())
	// 標準入力(入力装置,OSが提供するデータ入力機能)を取得
	s := bufio.NewScanner(os.Stdin)
	// データの入力元をscanメソッドで読み込み。データがあればtrueを返して実行
	for s.Scan() {
		if err := readLine("read.txt"); err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	}
}
