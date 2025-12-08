package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	httphandler "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/inbound/http"
)

type Server struct {
	router *http.ServeMux
	server *http.Server
}

func NewServer(host, port string) *Server {
	router := httphandler.NewRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		router: router,
		server: server,
	}
}

func (s *Server) Start() error {
	log.Printf("Server listening on %s\n", s.server.Addr)
	log.Println("Available endpoints:")
	log.Println("  GET  /health")
	log.Println("  GET  /swagger/index.html")
	log.Println("  GET  /api/v1/orders/{code}/total")
	log.Println("  GET  /api/v1/customers/{code}/orders")
	log.Println("  GET  /api/v1/customers/{code}/orders/count")
	log.Println("  POST /api/v1/orders")

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Close()
}
