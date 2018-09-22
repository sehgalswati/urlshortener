package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/sehgalswati/urlshortener/config"
	"github.com/sehgalswati/urlshortener/handler"
	"github.com/sehgalswati/urlshortener/backend/postgres"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "path of the config file")

	flag.Parse()

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		fmt.Printf("Error getting config %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("after parsing the values are ",config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DB,  config.Server.Host, config.Server.Port, config.Options.Prefix)
	fmt.Printf("Done with getting config from the file\n")
	
	srvc, err := postgres.New(config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DB)
	if err != nil {
		fmt.Printf("Got error here in opening creating new postgress service %v\n", err)
		log.Fatal(err)
	}
	fmt.Printf("Done initiating postgres service\n")
	defer srvc.Close()

	// Create a server
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: handler.New(config.Options.Prefix, srvc),
		}

	// Check for a closing signal
	go func() {
		// Graceful shutdown
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Cannot shut down server : %v", err)
		} else {
			log.Println("Server shut down successful")
		}
	}()

	// Start server
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}
