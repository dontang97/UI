package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/router"
	"github.com/dontang97/ui/ui"
)

func main() {
	db := &pg.Client{}
	db.Connect()
	defer db.Disconnect()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srv := router.Route(ui.New(db))
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
