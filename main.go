package main

import (
	"flag"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Default CLI argument values
const (
	DefaultAddr = ":8000"
)

// CLI arguments
var (
	addr string
)

// Init func
func init() {
	flag.StringVar(&addr, "addr", DefaultAddr, "Server host and port")
}

// Main func
func main() {
	// Parse CLI arguments
	flag.Parse()

	// Listen and serve
	log.Print(fasthttp.ListenAndServe(addr, ApplyMiddleware(NewRouter().Handler)))
}

// Create new http router
func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	// TODO: routes here

	return router
}

// Apply all middleware
func ApplyMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {

	h = MiddlewareLogging(h)

	return h
}
