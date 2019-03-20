package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
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
	fmt.Println("  Configuration checking - calculator v0.8   ")
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

	// Use fasthttprouter mux
	router := fasthttprouter.New()
	router.GET("/", rootHandler)
	router.GET("/echo/:name", echoHandler)
	router.GET("/factorialIterative/:number", factorialIterativeHandler)
	router.GET("/factorialRecursive/:number", factorialRecursiveHandler)


	// Defining server
	srv := &fasthttp.Server{
		// https://stackoverflow.com/questions/29334407/creating-an-idle-timeout-in-go
		//WriteTimeout: 					time.Second * 60,
		ReadTimeout: 						time.Second * 20,
		//IdleTimeout:  					time.Second * 120,
		SleepWhenConcurrencyLimitsExceeded: time.Second * 5,
		Handler:                            router.Handler,
	}

	// Lanuch server in a thread
	go func() {
		fmt.Println("Starting server...")
		log.Info().Msg("Starting server...")
		if err := srv.ListenAndServe(viper.GetString("port")); err != nil {
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
