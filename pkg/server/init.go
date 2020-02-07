package server

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/server/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var (
	GIN           *gin.Engine
	listenAddress string
)

func Init() {
	GIN = gin.New()
	GIN.Use(gin.Logger())
	GIN.Use(middlewares.CORS())
	GIN.Use(middlewares.Recovery())
	GIN.Use(middlewares.AccessControl())
	GIN.NoMethod(middlewares.Handler404())
	GIN.NoRoute(middlewares.Handler404())
	listenAddress = os.Getenv("DOMAINMAN_SERVER_LISTENADDRESS")
	registerRouter()
}

func Run() (err error) {
	err = database.Connect()
	if err != nil {
		log.Println("Failed to connect database")
		return
	}
	err = GIN.Run(listenAddress)
	if err != nil {
		log.Println("Failed to listen server")
		return
	}
	return
}
