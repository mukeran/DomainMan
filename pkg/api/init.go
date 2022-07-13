package api

import (
	"DomainMan/pkg/api/middlewares"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	GIN *gin.Engine
	srv *http.Server
)

func Init() {
	GIN = gin.New()
	GIN.Use(gin.Logger())
	GIN.Use(middlewares.CORS())
	GIN.Use(middlewares.Recovery())
	GIN.Use(middlewares.AccessControl())
	GIN.NoMethod(middlewares.Handler404())
	GIN.NoRoute(middlewares.Handler404())
	registerRouter()
}

func Run(listen string) (err error) {
	srv = &http.Server{
		Addr:    listen,
		Handler: GIN,
	}
	err = srv.ListenAndServe()
	if err != nil {
		return
	}
	return
}

func Stop() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Failed to stop server: %v", err)
	}
	return
}
