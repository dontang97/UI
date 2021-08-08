package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/dontang97/ui/router"
	"github.com/dontang97/ui/secret"
	"github.com/dontang97/ui/ui"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	keyDir := flag.String("jwt-key-folder", "./secret", "the folder of RSA key pair used to generate JWT")

	flag.Parse()

	secret.InitSecretKey(*keyDir)

	_ui := ui.New()
	_ui.Connect()
	defer _ui.Disconnect()

	srv := router.Route(_ui)
	go func() {
		fmt.Println("Start ui server...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	time.Sleep(3 * time.Second)
	os.Exit(0)
}
