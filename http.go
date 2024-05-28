package main

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*sync.RWMutex
	Pairs     map[string]bool
	exchanges map[string]*Exchange
}

func NewServer() *Server {
	return &Server{
		RWMutex:   &sync.RWMutex{},
		Pairs:     make(map[string]bool),
		exchanges: make(map[string]*Exchange),
	}
}

func (s *Server) Subscribes(c *gin.Context) {
	s.RLock()
	defer s.RUnlock()

	for pair := range s.Pairs {
		c.Writer.Write([]byte(pair + "\n"))
	}
}

func (s *Server) Subscribe(c *gin.Context) {
	s.Lock()
	defer s.Unlock()

	pairs := strings.Split(c.PostForm("pairs"), ",")
	for _, pair := range pairs {
		for _, exchange := range s.exchanges {
			exchange.Subscribe(pair)
		}
	}
}

func (s *Server) Unsubscribe(c *gin.Context) {
	s.Lock()
	defer s.Unlock()

	pairs := strings.Split(c.PostForm("pairs"), ",")
	for _, pair := range pairs {
		for _, exchange := range s.exchanges {
			exchange.Unsubscribe(pair)
		}
	}
}

func (s *Server) AddExchange(exchange *Exchange) {
	s.Lock()
	defer s.Unlock()

	s.exchanges[exchange.Name] = exchange
}
