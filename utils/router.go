package utils

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(buildTime string) *gin.Engine {
	if buildTime == "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// use logger and recovery by default
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// check platform
	cdnPlatform := os.Getenv("CDN_PLATFORM")
	switch cdnPlatform {
	case "cloudflare":
		router.TrustedPlatform = gin.PlatformCloudflare
	case "google":
		router.TrustedPlatform = gin.PlatformGoogleAppEngine
	default:
		// check trustdd proxies
		trustedProxies := os.Getenv("TRUSTED_PROXIES")
		if trustedProxies != "" {
			router.SetTrustedProxies(strings.Split(trustedProxies, ","))
		} else {
			router.SetTrustedProxies(nil)
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buildTime": buildTime,
			"runMode":   gin.Mode(),
		})
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	router.GET("/readyz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	return router
}

func ListenAndServe(router *gin.Engine) {
	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		listener, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			log.Fatalf("Listen: %s\n", err)
		}
		log.Printf("Server Started on %s\n", srv.Addr)
		err = srv.Serve(listener)
		if err != nil {
			log.Fatalf("Serve: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server...")

	shutdownWait := GetIntValue(os.Getenv("SHUTDOWN_WAIT"), 5)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownWait)*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error Server Shutdown: %s\n", err)
	}
	log.Println("Server Shutdown Normally.")
}
