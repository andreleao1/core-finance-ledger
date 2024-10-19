package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	contextPath string
	port        string
}

func NewHttpServer(contextPath string, port string) *HttpServer {
	return &HttpServer{
		contextPath: contextPath,
		port:        port,
	}
}

func (hs *HttpServer) StartHttpServer() error {
	httpServer := gin.Default()

	httpServer.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("HTTP server started on " + hs.port)
	return httpServer.Run(hs.port) // Retorna um erro se falhar
}
