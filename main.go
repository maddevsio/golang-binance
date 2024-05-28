package main

import (
	"github.com/gin-gonic/gin"
)

const BinanceAddr = "wss://stream.binance.com/ws"
const BinanceOrigin = "https://stream.binance.com"
const BinanceReconnectDelay = 3

func main() {

	binance := NewExchange("binance", BinanceAddr, BinanceOrigin)
	go binance.Run()

	s := NewServer()
	s.AddExchange(binance)

	r := gin.Default()
	r.GET("/subscribes", s.Subscribes)
	r.POST("/subscribe", s.Subscribe)
	r.POST("/unsubscribe", s.Unsubscribe)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
