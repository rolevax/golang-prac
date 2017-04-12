package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/rolevax/golang-prac/chat/messages"
)

type Conns struct {
	conns []net.Conn
	mutex sync.Mutex
}

func (this *Conns) Add(conn net.Conn) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.conns = append(this.conns, conn)
}

func (this *Conns) Remove(conn net.Conn) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for i, c := range this.conns {
		if c == conn {
			// replace by back and pop back
			if len(this.conns) > 1 {
				this.conns[i] = this.conns[len(this.conns)-1]
			}
			this.conns = this.conns[:len(this.conns)-1]
		}
	}
}

func (this *Conns) Broadcast(pb proto.Message) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, c := range this.conns {
		messages.WritePb(c, pb)
	}
}

var conns Conns

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("server listening at 8080")
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println(err)
		} else {
			go loop(conn)
		}
	}
}

func loop(conn net.Conn) {
	defer conn.Close()
	var username string

	reader := bufio.NewReader(conn)
	for {
		var req messages.Req
		err := messages.ReadPb(reader, &req)
		if err != nil {
			return // disconnect client
		}
		switch req.Type {
		case messages.CONNECT:
			conns.Add(conn)
			defer conns.Remove(conn)
			username = req.Who
			log.Println(username, "connected")
			defer log.Println(username, "disconnected")
			resp := messages.Resp{
				Type: messages.CONNECT,
				What: username + " joined room",
			}
			conns.Broadcast(&resp)
		case messages.SAY:
			s := fmt.Sprintf("%s: %s", username, req.What)
			log.Println(s)
			resp := messages.Resp{
				Type: messages.SAY,
				What: s,
			}
			conns.Broadcast(&resp)
		case messages.NICK:
			s := fmt.Sprintf("%s renamed to %s", username, req.What)
			log.Println(s)
			resp := messages.Resp{
				Type: messages.NICK,
				What: s,
			}
			conns.Broadcast(&resp)
			username = req.What
		}
	}
}
