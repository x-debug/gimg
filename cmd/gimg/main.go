package main

import (
	"context"
	"flag"
	"gimg/cache"
	"gimg/config"
	"gimg/handlers"
	lg "gimg/logger"
	"gimg/pkg"
	"gimg/processor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

func initEngine(conf *config.Config) {
	engine = processor.NewEngine(processor.Imagick, conf.Engine)
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
	initEngine(conf)
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
	cached, err := cache.NewCache(conf.Cache, logger)
	if err != nil {
		logger.Error("Cache initialize error", lg.Error(err))
		return
	}
	ctx := pkg.CreateCtx(conf, cached, logger, engine)

	//Register handlers
	router.StaticFile("/favicon.ico", "./resources/favicon-16x16.png")
	router.StaticFile("/demo", "./examples/demo.html")

	if conf.Auth.Close {
		router.POST("/upload", handlers.UploadHandler(ctx))
		router.GET("/:hash", handlers.GetHandler(ctx))
		router.GET("", handlers.RemoteGetHandler(ctx))
	} else {
		var group_router *gin.RouterGroup

		if conf.Auth.Type == "basic" {
			group_router = router.Group("/", gin.BasicAuth(gin.Accounts{conf.Auth.User: conf.Auth.Password}))
		} else {
			logger.Error("Auth configuare type error")
			return
		}
		group_router.POST("/upload", handlers.UploadHandler(ctx))
		group_router.GET("/:hash", handlers.GetHandler(ctx))
		group_router.GET("", handlers.RemoteGetHandler(ctx))
	}

	logger.Info("Http listen ", lg.Int("Port", conf.Port))
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(conf.Port),
		Handler: router,
	}

	//Start http server
	go func() {
		_ = srv.ListenAndServe()
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
