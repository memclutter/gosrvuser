package main

import (
	"log"
	"strings"

	"time"

	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

func MiddlewareHeader(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Set common headers
		ctx.Response.Header.SetContentType("application/json")

		// Call handler
		h(ctx)
	})
}

// Logging middleware
func MiddlewareLogging(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Call handler
		h(ctx)

		// Extract data
		method := strings.ToUpper(string(ctx.Method()))
		request := string(ctx.RequestURI())
		status := ctx.Response.StatusCode()

		log.Printf("%v %v - %v", method, request, status)
	})
}

// Database middleware
func MiddlewareMongodb(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Copy session
		session := mongodb.Copy()

		// Get default database
		db := session.DB("")

		// Save in context
		ctx.SetUserValue("mongodb", session)
		ctx.SetUserValue("mongodb.db", db)

		// Run handler
		h(ctx)
	})
}

// AMQP middleware
func MiddlewareAmqp(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Select amqp channel
		var ch *amqp.Channel
		for {
			var err error

			if mq != nil {
				if ch, err = mq.Channel(); err == nil {
					break
				}
			}

			log.Println("Trying to reselect amqp channel")
			time.Sleep(500 * time.Millisecond)
		}

		// Save in context
		ctx.SetUserValue("amqp.ch", ch)

		// Run handler
		h(ctx)
	})
}

func MiddlewareResponse(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Call handler
		h(ctx)

		// Get status code
		status := ctx.Response.StatusCode()

		// Detect success flag
		success := status >= fasthttp.StatusOK && status < fasthttp.StatusBadRequest

		// Trying get error message
		messageInterface := ctx.UserValue("response.message")
		message := ""
		if !success {
			message = fasthttp.StatusMessage(status)
		}
		if messageInterface != nil {
			message = messageInterface.(string)
		}

		// Get data
		data := ctx.UserValue("response.data")

		// Build response and set in body
		response := Response{Success: success, Status: status}

		if data != nil {
			response.Data = data
		}

		if len(message) != 0 {
			response.Message = message
		}

		// Set body
		writer := json.NewEncoder(ctx.Response.BodyWriter())
		writer.Encode(response)
	})
}
