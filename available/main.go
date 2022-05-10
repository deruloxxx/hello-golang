package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func exists(domain string) (bool, error) {
	const whoisServer string = "com.whois-servers.net"
	// whoisServer port 43に対してnet.Dialで接続を開く
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return false, err
	}
	// whoisServerの接続情報に関わらず最後に実行される
	defer conn.Close()
	// ドメイン名を送信
	conn.Write([]byte(domain + "\r\n"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), "no match") {
			return false, nil
		}
	}
	return true, nil
}

var marks = map[bool]string{true: "あり", false: "なし"}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		domain := s.Text()
		fmt.Print(domain, "")
		exist, err := exists(domain)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(marks[!exist])
		// WHOISサーバーの負荷上昇を避けるために処理を1秒間休止
		time.Sleep(1 * time.Second)
	}
}
