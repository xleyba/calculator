package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"math/big"
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
func rootHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBody([]byte("Hello!"))
	// If we arrived here then everything is OK. :)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

// Return echo message
func echoHandler(ctx *fasthttp.RequestCtx) {

	log.Debug().Msgf("Received: %s", ctx.UserValue("name"))
	fmt.Fprintf(ctx, "Hello: %s", ctx.UserValue("name"))
	// If we arrived here then everything is OK. :)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

}

// Handle iterative path and calls iterative calculation
func factorialIterativeHandler(ctx *fasthttp.RequestCtx) {

	log.Debug().Msgf("Received: %s", ctx.UserValue("number"))

	number, er := StrToInt(fmt.Sprintf("%s", ctx.UserValue("number")))

	if er != nil {
		log.Error().Msg("Error calculating number")
		ctx.Error(er.Error(), fasthttp.StatusInternalServerError)
		ctx.Response.Header.Set("Status", strconv.Itoa(fasthttp.StatusInternalServerError))
		ctx.Response.SetBody([]byte("500 - Internal server error processing data"))
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	fmt.Fprintf(ctx, "%s", factorialIter(number))


}

// Handle recursive path and calls recursive calculation
func factorialRecursiveHandler(ctx *fasthttp.RequestCtx) {

	mynumbera, er := strconv.ParseInt(fmt.Sprintf("%s", ctx.UserValue("number")),
		10, 64)

	if er != nil {
		log.Error().Msg("Error calculating number")
		ctx.Error(er.Error(), fasthttp.StatusInternalServerError)
		ctx.Response.Header.Set("Status", strconv.Itoa(fasthttp.StatusInternalServerError))
		ctx.Response.SetBody([]byte("500 - Internal server error processing data"))
		return
	}

	// Convert int64 to big.Int
	mynumberBig := big.NewInt(mynumbera)
	// Get a pointer
	numberPointer := &mynumberBig


	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	fmt.Fprintf(ctx, "%s", factorialRecursive(*numberPointer))

}
