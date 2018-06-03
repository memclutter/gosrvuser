package main

import (
	"log"
	"strings"

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
