package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luminishion/wildberries-l0/orders"
	"github.com/luminishion/wildberries-l0/subscriber"
	"github.com/luminishion/wildberries-l0/web"
)

func main() {
	orders.Connect()
	web.RunHTTP()
	subscriber.RunNats()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutdown ....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	web.StopHTTP(ctx)
	subscriber.StopNats()
	orders.Close()

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Exiting")
}
