package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"flag"
	"context"
	"os/signal"
)

// Configuration
type Configuration struct {
	Port string
}

// Convert String to int
func StrToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

// Calculate factorial with iterative method
func factorialIter(x int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= x; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}

	return result
}

// Calculates factorial recursive
func factorialRecursive(x *big.Int) *big.Int {
	n := big.NewInt(1)
	if x.Cmp(big.NewInt(0)) == 0 {
		return n
	}
	return n.Mul(x, factorialRecursive(n.Sub(x, n)))
}

// Return default message for root routing
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Return echo message
func echoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	message := params["message"]

	fmt.Fprintf(w, "%s", message)
}

// Handle iterative path and calls iterative calculation
func factorialIterativeHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	number, er := StrToInt(params["number"])

	if er != nil {
		fmt.Fprintln(w, "Error calculating number")
		log.Fatal("Error calculating number")
		return
	}

	fmt.Fprintf(w, "%s", factorialIter(number))

}

// Handle recursive path and calls recursive calculation
func factorialRecursiveHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	// Convert received string to int64
	mynumbera, er := strconv.ParseInt(params["number"], 10, 64)

	if er != nil {
		fmt.Fprintln(w, "Error calculating number")
		log.Fatal("Error calculating number")
		return
	}

	// Convert int64 to big.Int
	mynumberbig := big.NewInt(mynumbera)
	// Get a pointer
	numberPointer := &mynumberbig

	if er != nil {
		fmt.Fprintln(w, "Error calculating number")
		log.Fatal("Error calculating number")
		return
	}

	fmt.Fprintf(w, "%s", factorialRecursive(*numberPointer))

}



func main() {

	// Used for gracefull server shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Config data
	port := ":9596"
	log.Println("Starting server")

	// Start to read conf file
	log.Println("\n\n")
	log.Println("=============================================")
	log.Println("   Configuration checking - calculator v0.6")
	log.Println("=============================================")
	file, err := os.Open("conf.json")

	if err != nil {
		log.Println("No conf file, using port 9596 by default")
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err := decoder.Decode(&configuration)

		if err != nil {
			fmt.Println("error:", err)
			log.Fatal()
		} else {

			// Check port parameter
			if len(configuration.Port) == 0 {
				log.Println("-- No port inf config file, using: ", port)
			} else {
				port = configuration.Port
				log.Println("-- Using port: ", port)
			}

		}
	}

	log.Println("=============================================")


	router := mux.NewRouter() //.StrictSlash(true)

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/echo/{message}", echoHandler).Methods("GET")
	router.HandleFunc("/factorialIterative/{number}", factorialIterativeHandler).Methods("GET")
	router.HandleFunc("/factorialRecursive/{number}", factorialRecursiveHandler).Methods("GET")

	// set timeout
	muxWithMiddlewares := http.TimeoutHandler(router, time.Second*3, "Timeout!")

	srv := &http.Server{
		Addr:    port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		// Using just the read parameter due to this article
		// https://stackoverflow.com/questions/29334407/creating-an-idle-timeout-in-go
		//WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		////:  time.Second * 120,
		//Handler: router,
		Handler: muxWithMiddlewares,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Running server....")

		//log.Fatal(http.ListenAndServe(port, router))
		//log.Println(s.ListenAndServe())
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Println("\n\n")
	log.Println("shutting down")
	log.Println("Goddbye!....")
	os.Exit(0)


}
