package web

import (
	"context"
	"log"
	"net/http"

	"github.com/luminishion/wildberries-l0/config"
	"github.com/luminishion/wildberries-l0/web/routes"

	"github.com/gin-gonic/gin"
)

var sv *http.Server

func RunHTTP() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	routes.Orders(router)
	routes.Generate(router)

	sv = &http.Server{
		Addr:    config.WebAddr,
		Handler: router,
	}

	go func() {
		if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("web listen: ", err)
		}
	}()
}

func StopHTTP(ctx context.Context) {
	if err := sv.Shutdown(ctx); err != nil {
		log.Fatal("web shutdown: ", err)
	}
}
