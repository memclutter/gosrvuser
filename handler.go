package main

import (
	"time"

	"github.com/valyala/fasthttp"
)

// Handle GET /status. This endpoint need for health check.
func HandleStatus(ctx *fasthttp.RequestCtx) {
	data := ResponseDataStatus{}
	data.Time = time.Now().UTC()

	data.Health.Db = ctx.UserValue("mongodb.db") != nil
	data.Health.Amqp = ctx.UserValue("amqp.ch") != nil

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetUserValue("response.data", data)
}
