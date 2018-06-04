package main

import (
	"flag"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/globalsign/mgo"
	"github.com/valyala/fasthttp"
)

// Default CLI argument values
const (
	DefaultAddr       = ":8000"
	DefaultMongodbUrl = "mongodb://localhost:27017/user"
)

// CLI arguments
var (
	addr       string
	mongodbUrl string
)

// Global variables (database connection pools, etc)
var (
	err     error
	mongodb *mgo.Session
)

// Init func
func init() {
	flag.StringVar(&addr, "addr", DefaultAddr, "Server host and port")
	flag.StringVar(&mongodbUrl, "mongodbUrl", DefaultMongodbUrl, "Mongodb connection url")
}

// Main func
func main() {
	// Parse CLI arguments
	flag.Parse()

	// Connect to mongodb
	mongodb, err = mgo.Dial(mongodbUrl)
	if err != nil {
		log.Fatalf("mgo.Dial: %v", err)
		return
	}
	defer mongodb.Close()

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
	h = MiddlewareMongodb(h)

	return h
}
