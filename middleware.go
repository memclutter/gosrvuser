package main

import (
	"log"
	"strings"

	"time"

	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

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

			if ch, err = mq.Channel(); err == nil {
				break
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
