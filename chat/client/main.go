package main

import (
	"bufio"
	"log"
	"net"

	console "github.com/AsynkronIT/goconsole"
	"github.com/rolevax/golang-prac/chat/messages"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}

	req := messages.Req{
		Type: messages.CONNECT,
		Who:  "aaa",
	}
	messages.WritePb(conn, &req)

	go readLoop(conn)

	cons := console.NewConsole(func(text string) {
		reqSay := messages.Req{
			Type: messages.SAY,
			What: text,
		}
		messages.WritePb(conn, &reqSay)
	})
	//write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		reqSay := messages.Req{
			Type: messages.NICK,
			What: newNick,
		}
		messages.WritePb(conn, &reqSay)
	})
	cons.Run()
}

func readLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		var resp messages.Resp
		err := messages.ReadPb(reader, &resp)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(resp.What)
	}
}
