package main

import (
	"context"
	"fmt"
	"github.com/Celbux/dataapi/business/dataapi"
	"github.com/Celbux/dataapi/services/dataapi/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

func main() {

	err := run(log.New(os.Stdout, "", 0))
	if err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}

}

func run(log *log.Logger) error {

	// Configuration uses github.com/ardanlabs/conf library
	// Your program configuration is attempted to be retrieved in the priority:
	// 1) Environment variable
	// 2) CMD flag
	// 3) Else the default value will be used
	defer log.Println("main: Completed")
	var cfg struct {
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:8082"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:0s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}
	namespace := "DATA_API"
	if err := conf.Parse(os.Args[1:], namespace, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(namespace, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(namespace, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n\n%v\n\n", out)

	// Dependency Injection: Create our Services with their dependencies to
	// attach on for later access via receiver functions
	dataAPI := handlers.DataAPIHandlers{
		Service: dataapi.DataAPIService{Log: log},
	}

	// Make a channel to listen for an interrupt or terminate signal from the
	// OS. Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Create the server that will listen and serve
	api := http.Server{
		Addr: cfg.Web.APIHost,
		Handler: handlers.API(
			dataAPI,
			log,
			shutdown,
		),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this
	// error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(
			context.Background(),
			cfg.Web.ShutdownTimeout,
		)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			_ = api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil

}
