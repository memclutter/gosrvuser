package main

import (
	"flag"
	"log"

	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/globalsign/mgo"
	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

// Default CLI argument values
const (
	DefaultAddr       = ":8000"
	DefaultMongodbUrl = "mongodb://localhost:27017/user"
	DefaultAmqpUrl    = "amqp://guest:guest@localhost:5672/"
)

// CLI arguments
var (
	addr       string
	mongodbUrl string
	amqpUrl    string
)

// Global variables (database connection pools, etc)
var (
	err          error
	mongodb      *mgo.Session
	mq           *amqp.Connection
	mqCloseError chan *amqp.Error
)

// Init func
func init() {
	flag.StringVar(&addr, "addr", DefaultAddr, "Server host and port")
	flag.StringVar(&mongodbUrl, "mongodbUrl", DefaultMongodbUrl, "Mongodb connection url")
	flag.StringVar(&amqpUrl, "amqpUrl", DefaultAmqpUrl, "AMQP connection url")

	mqCloseError = make(chan *amqp.Error)
}

// Main func
func main() {
	// Parse CLI arguments
	flag.Parse()

	// Connect to mongodb
	for {
		mongodb, err = mgo.Dial(mongodbUrl)
		if err == nil {
			break
		}

		log.Println(err)
		log.Printf("Trying to reconnect to Mongodb at %s\n", mongodbUrl)
		time.Sleep(500 * time.Millisecond)
	}
	defer mongodb.Close()

	// Connect to amqp broker
	go func() {
		mqErr := new(amqp.Error)

		for {
			mqErr = <-mqCloseError
			if mqErr != nil {
				log.Printf("Connecting to %s\n", amqpUrl)

				for {
					mq, err = amqp.Dial(amqpUrl)

					if err == nil {
						break
					}

					log.Println(err)
					log.Printf("Trying to reconnect to AMQP at %s\n", amqpUrl)
					time.Sleep(500 * time.Millisecond)
				}

				mqCloseError = make(chan *amqp.Error)
				mq.NotifyClose(mqCloseError)
			}
		}
	}()
	mqCloseError <- amqp.ErrClosed
	defer mq.Close()

	// Listen and serve
	log.Print(fasthttp.ListenAndServe(addr, ApplyMiddleware(NewRouter().Handler)))
}

// Create new http router
func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.NotFound = HandleNotFound
	router.MethodNotAllowed = HandleMethodNotAllowed

	router.GET("/status", HandleStatus)
	router.POST("/sign-up", HandleSignUp)

	return router
}

// Apply all middleware
func ApplyMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {

	h = MiddlewareLogging(h)
	h = MiddlewareMongodb(h)
	h = MiddlewareAmqp(h)
	h = MiddlewareHeader(h)
	h = MiddlewareResponse(h)

	return h
}
