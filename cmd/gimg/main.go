package main

import (
	"context"
	"flag"
	"gimg/config"
	"gimg/handlers"
	lg "gimg/logger"
	"gimg/pkg"
	"gimg/processor"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	port      int
	savePath  string
	configure string
	engine    processor.Engine
)

func init() {
	flag.IntVar(&port, "port", 0, "port of server listen to")
	flag.StringVar(&savePath, "save_path", "", "path of save to")
	flag.StringVar(&configure, "conf", "./gimg.yml", "path of configure")
}

func initProcessor(conf *config.Config) {
	engine = processor.NewEngine(processor.Imagick)
	engine.Initialize()
}

func overrideConf(conf *config.Config) *config.Config {
	if port != 0 {
		conf.Port = port
	}

	if savePath != "" {
		conf.Engine.SavePath = savePath
	}

	return conf
}

func main() {
	flag.Parse()
	conf, err := config.Load(configure)
	if err != nil {
		log.Fatal("Read configure file error :", err)
	}
	conf = overrideConf(conf)

	//Init logger
	logger := lg.New(conf.Logger)
	initProcessor(conf)
	defer engine.Terminate()

	//Set run mode
	if conf.Debug {
		logger.Info("Run ", lg.String("Mode", gin.DebugMode))
		gin.SetMode(gin.DebugMode)
	} else {
		logger.Info("Run ", lg.String("Mode", gin.ReleaseMode))
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20
	ctx := pkg.CreateCtx(conf, logger, engine)

	//Register handlers
	router.POST("/upload", handlers.UploadHandler(ctx))
	router.StaticFile("/demo", "./examples/demo.html")
	router.GET("/:hash", handlers.GetHandler(ctx))

	logger.Info("Http listen ", lg.Int("Port", conf.Port))
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(conf.Port),
		Handler: router,
	}

	//Start http server
	go func() {
		err = srv.ListenAndServe()
		logger.Error("Serve http ", lg.Error(err))
	}()

	//Create signal wait until server exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server")

	//Server forced shutting down if context timeout
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.Panic("Server forced to shutdown", lg.Error(err))
	}
	
	logger.Info("Server exited")
}
