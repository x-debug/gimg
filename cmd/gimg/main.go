package main

import (
	"context"
	"flag"
	"gimg/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port string
)

func init()  {
	flag.StringVar(&port, "port", "8080", "port of server listen to")
}

func main() {
	router := gin.Default()

	//register handlers
	router.GET("/", handlers.GetHandler)
	router.POST("upload", handlers.UploadHandler)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: router,
	}

	go func() {
		log.Println(srv.ListenAndServe())
	}()

	//create signal wait until server exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//server forced shutting down if context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
