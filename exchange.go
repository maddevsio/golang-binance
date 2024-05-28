package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Command struct {
	Name string
	Args []string
}

type Exchange struct {
	sync.RWMutex

	Name string

	addr   string
	origin string
	ws     *websocket.Conn
	reqID  int64
}

func NewExchange(name, addr, origin string) *Exchange {
	return &Exchange{
		Name:   name,
		addr:   addr,
		origin: origin,
		reqID:  1,
	}
}

func (e *Exchange) Run() {
	for {
		ws, err := websocket.Dial(e.addr, "", e.origin)
		if err != nil {
			log.Println(fmt.Errorf("unable to connect to %s, reconnecting... %s", e.Name, err))
			time.Sleep(BinanceReconnectDelay * time.Second)
			continue
		}
		e.ws = ws

		e.readMessages()

		log.Printf("%s run", e.Name)
	}
}

func (e *Exchange) readMessages() {
	var data []byte
	for {
		err := websocket.Message.Receive(e.ws, &data)
		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(data))
	}
}

func (e *Exchange) sendMessage(message []byte) {
	log.Printf("send: %s: %s", e.Name, string(message))
	err := websocket.Message.Send(e.ws, string(message))
	if err != nil {
		log.Println(err)
		return
	}
}

func (e *Exchange) FormatPair(pair string) string {

	pair = strings.ReplaceAll(pair, "/", "")
	pair = strings.ToLower(pair)

	return pair
}

func (e *Exchange) incrReqID() {
	e.Lock()
	defer e.Unlock()
	e.reqID++
}

func (e *Exchange) Subscribe(pair string) {
	req := "{\"method\": \"SUBSCRIBE\", \"params\": [\"%s@bookTicker\"], \"id\": %d}\n"

	e.sendMessage([]byte(fmt.Sprintf(req, e.FormatPair(pair), e.reqID)))
	e.incrReqID()

}

func (e *Exchange) Unsubscribe(pair string) {
	req := "{\"method\": \"UNSUBSCRIBE\", \"params\": [\"%s@bookTicker\"], \"id\": %d}\n"

	e.sendMessage([]byte(fmt.Sprintf(req, e.FormatPair(pair), e.reqID)))
	e.incrReqID()
}
