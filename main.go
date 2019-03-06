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

	port := ":9596"
	log.Println("Starting server")

	// Start to read conf file
	log.Println("calculator v0.5")
	log.Println("=============================================")
	log.Println("         Configuration checking")
	log.Println("=============================================")
	file, err := os.Open("conf.json")


	if (err != nil) {
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

	log.Println("Running server....")

	log.Fatal(http.ListenAndServe(port, router))

}
