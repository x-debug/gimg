package main

import (
	"context"
	"flag"
	"gimg/handlers"
	"gimg/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port     string
	savePath string
)

func init() {
	flag.StringVar(&port, "port", "8080", "port of server listen to")
	flag.StringVar(&savePath, "save_path", "./", "path of save to")
}

func main() {
	router := gin.Default()

	router.MaxMultipartMemory = 100 << 20
	ctx := pkg.CreateCtx(savePath)

	//register handlers
	router.GET("/", handlers.GetHandler(ctx))
	router.POST("/upload", handlers.UploadHandler(ctx))
	router.StaticFile("/demo", "./examples/demo.html")

	srv := &http.Server{
		Addr:    ":" + port,
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
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
