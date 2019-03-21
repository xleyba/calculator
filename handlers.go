package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"html"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

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
func index(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Got root")
	fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
}

// Return echo message

func echoHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	fmt.Fprintf(w, "Hello: %s", params["message"])

}

// Handle iterative path and calls iterative calculation
func factorialIterativeHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	log.Debug().Msgf("Received: %s", params["number"])

	number, er := StrToInt(fmt.Sprintf("%s", params["number"]))

	if er != nil {
		log.Error().Msg("Error calculating number")
		w.Header().Set("Status", "500")
		fmt.Fprint(w, "500 - Internal server error processing data")
		return
	}

	fmt.Fprintf(w, "%s", factorialIter(number))

}

// Handle recursive path and calls recursive calculation
func factorialRecursiveHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	log.Debug().Msgf("Received: %s", params["number"])

	mynumbera, er := strconv.ParseInt(params["number"],
		10, 64)

	if er != nil {
		log.Error().Msgf("Error calculating number %s", er.Error())
		w.Header().Set("Status", "500")
		fmt.Fprint(w, "500 - Internal server error processing data")
		return
	}

	// Convert int64 to big.Int
	mynumberBig := big.NewInt(mynumbera)
	// Get a pointer
	numberPointer := &mynumberBig

	fmt.Fprintf(w, "%s", factorialRecursive(*numberPointer))

}
