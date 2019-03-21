package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Load conf and default values
func start() {

	fmt.Println("Starting server")

	// Start to read conf file
	fmt.Print("\n\n")
	fmt.Println("=============================================")
	fmt.Println("  Configuration checking - calculator v1.0   ")
	fmt.Println("=============================================")

	// loading configuration
	viper.SetConfigName("conf")          // name of config file (without ext)
	viper.AddConfigPath(".")             // default path for conf file
	viper.SetDefault("port", ":9596")    // default port value
	viper.SetDefault("loglevel", "info") // default port value
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		fmt.Printf("Fatal error config file: %v \n", err)
		panic(err)
	}

	fmt.Println("-- Using port:       ", viper.GetString("port"))

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "@timestamp"
	switch viper.GetString("loglevel") {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	fmt.Println("-- Log level:        ", viper.GetString("loglevel"))

	fmt.Println("=============================================")

}

func main() {

	start()

	// CPU profiling by default
	//defer profile.Start().Stop()
	// Memory profiling
	//defer profile.Start(profile.MemProfile).Stop()

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/echo/{message}", echoHandler)
	r.HandleFunc("/factorialIterative/{number}", factorialIterativeHandler)
	r.HandleFunc("/factorialRecursive/{number}", echoHandlerD(viper.GetString("calledServiceURL"),
		SetClientD()))

	srv := &http.Server{
		Addr:           viper.GetString("port"),
		Handler:        r,
		ReadTimeout:    2 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    time.Second * 2,
	}

	srv.ConnState = func(c net.Conn, cs http.ConnState) {
		switch cs {
		case http.StateIdle:
			c.SetReadDeadline(time.Now().Add(time.Second * 2))
			log.Debug().Msgf("StateIddle: %s", c.LocalAddr().String())
		case http.StateNew:
			c.SetReadDeadline(time.Now().Add(time.Second * 2))
			log.Debug().Msgf("StateNew: %s", c.LocalAddr().String())
		case http.StateActive:
			log.Debug().Msgf("Active: %s", c.LocalAddr().String())
			c.SetReadDeadline(time.Now().Add(time.Second * 2))
		}
	}

	// Lanuch server in a thread
	go func() {
		fmt.Println("Starting server...")
		log.Info().Msg("Starting server...")
		if err := srv.ListenAndServe(); err != nil {
			log.Panic().Msgf("%s", err)
		}
	}()

	// Process a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	errShutdown := srv.Shutdown()
	if errShutdown != nil {
		panic(fmt.Sprintf("Error shutting down %s", errShutdown))
	}

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Print("\n\n")
	fmt.Println("shutting down")
	fmt.Println("Goddbye!....")
	os.Exit(0)

}
